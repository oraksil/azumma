package mq

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/streadway/amqp"
)

type Message struct {
	ReplyTo string // application use - address to reply to (ex: RPC)
	// CorrelationId string    // application use - correlation identifier
	// MessageId     string    // application use - message identifier
	Timestamp time.Time // application use - message timestamp
	Type      string    // application use - message type name
	// UserId        string    // application use - creating user - should be authenticated user
	// AppId         string    // application use - creating application id

	Payload []byte `mapstructure:"Body"`
}

type MessageService interface {
	Identifier() string
	AllParticipants() []string

	Send(to string, msgType string, payload interface{})
	Request(to string, msgType string, payload interface{}) Message
	Broadcast(msgType string, payload interface{})
}

type DefaultMessageServiceImpl struct {
	MqService *MqService
}

func (m *DefaultMessageServiceImpl) Identifier() string {
	return m.MqService.peerName
}

func (m *DefaultMessageServiceImpl) AllParticipants() []string {
	return nil
}

func (m *DefaultMessageServiceImpl) Send(to string, msgType string, payload interface{}) {
	m.publishMessage(m.MqService.exchangeP2P, to, msgType, payload, false)
}

func (m *DefaultMessageServiceImpl) Request(to string, msgType string, payload interface{}) Message {
	instantChannel := make(chan Message)
	m.MqService.recvMsgChannelsRpc[to] = instantChannel

	m.publishMessage(m.MqService.exchangeP2P, to, msgType, payload, true)

	recvMsg := <-instantChannel
	delete(m.MqService.recvMsgChannelsRpc, to)
	close(instantChannel)

	return recvMsg
}

func (m *DefaultMessageServiceImpl) Broadcast(msgType string, payload interface{}) {
	m.publishMessage(m.MqService.exchangeBroadcast, "", msgType, payload, false)
}

func (m *DefaultMessageServiceImpl) publishMessage(
	exchange, routingKey string,
	msgType string,
	payload interface{},
	needReply bool) {

	body, _ := json.Marshal(payload)

	msg := amqp.Publishing{
		ContentType: "application/json",
		Type:        msgType,
		Body:        body,
	}

	if needReply {
		msg.ReplyTo = m.MqService.peerName
	}

	err := m.MqService.channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		msg,
	)

	if err != nil {
		panic(err)
	}
}

type Context struct {
	message *Message
}

func (c *Context) GetMessage() *Message {
	return c.message
}

type HandlerFunc func(ctx *Context)

type MessageHandler interface {
	Routes() []Route
}

type MqService struct {
	channel           *amqp.Channel
	peerName          string
	exchangeP2P       string
	exchangeBroadcast string
	queueP2P          *amqp.Queue
	queueBroadcast    *amqp.Queue

	recvMsgChannel     chan Message
	recvMsgChannelsRpc map[string]chan Message

	handlers     []MessageHandler
	handlerFuncs map[string]HandlerFunc
}

func (mq *MqService) Run() {
	// generate unique peer name per host
	mq.peerName = "generated" // routing key

	// declare queue, it will create queue only if it doesn't exist
	mq.queueP2P = mq.enusureQueue(
		fmt.Sprintf("oraksil.mq.q-p2p-%s", mq.peerName))

	mq.queueBroadcast = mq.enusureQueue(
		fmt.Sprintf("oraksil.mq.q-broadcast-%s", mq.peerName))

	// bind the queue with p2p exchange and broadcast one respectively
	mq.bindQueue(mq.queueP2P, mq.exchangeP2P, mq.peerName)
	mq.bindQueue(mq.queueBroadcast, mq.exchangeBroadcast, "")

	// start consuming
	go mq.consumerP2PQueue()
	go mq.consumerBroadcastQueue()

	// message handler
	go mq.messageHandler()
}

func (mq *MqService) consumerP2PQueue() {
	for {
		msgs, _ := mq.channel.Consume(
			mq.queueP2P.Name, // queue
			"",               // consumer
			true,             // auto-ack
			false,            // exclusive
			false,            // no-local
			false,            // no-wait
			nil,              // args
		)

		for m := range msgs {
			var recv Message
			mapstructure.Decode(m, &recv)

			if recv.ReplyTo != "" {
				if _, ok := mq.recvMsgChannelsRpc[recv.ReplyTo]; !ok {
					panic("p2p channel must be created at sending.")
				}
				mq.recvMsgChannelsRpc[m.ReplyTo] <- recv
			} else {
				mq.recvMsgChannel <- recv
			}
		}
	}
}

func (mq *MqService) consumerBroadcastQueue() {
	for {
		msgs, _ := mq.channel.Consume(
			mq.queueBroadcast.Name, // queue
			"",                     // consumer
			true,                   // auto-ack
			false,                  // exclusive
			false,                  // no-local
			false,                  // no-wait
			nil,                    // args
		)

		for m := range msgs {
			var recvMsg Message
			mapstructure.Decode(m, &recvMsg)
			mq.recvMsgChannel <- recvMsg
		}
	}

}

func (mq *MqService) messageHandler() {
	func() {
		for msg := range mq.recvMsgChannel {
			if handlerFunc, ok := mq.handlerFuncs[msg.Type]; ok {
				handlerFunc(&Context{message: &msg})
			}
		}
	}()
}

func (mq *MqService) enusureQueue(name string) *amqp.Queue {
	q, err := mq.channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
	return &q
}

func (mq *MqService) bindQueue(q *amqp.Queue, exchange string, routingKey string) {
	err := mq.channel.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		panic(err)
	}
}

func (mq *MqService) AddHandler(handler MessageHandler) {
	mq.handlers = append(mq.handlers, handler)

	for _, r := range handler.Routes() {
		if _, ok := mq.handlerFuncs[r.MsgType]; ok {
			panic(fmt.Sprintf("handler for %s already exists.", r.MsgType))
		}

		mq.handlerFuncs[r.MsgType] = r.Handler
	}
}

func NewMqService(url, exchangeP2P, exchangeBroadcast string) *MqService {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return &MqService{
		channel:            ch,
		exchangeP2P:        exchangeP2P,
		exchangeBroadcast:  exchangeBroadcast,
		recvMsgChannel:     make(chan Message, 1024),
		recvMsgChannelsRpc: make(map[string]chan Message),
		handlerFuncs:       make(map[string]HandlerFunc),
	}
}
