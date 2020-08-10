package usecases

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/oraksil/azumma/internal/domain/models"
)

func (r *MockGameRepository) SaveConnectionInfo(connectionInfo *models.ConnectionInfo) (*models.ConnectionInfo, error) {
	args := r.Called(connectionInfo)
	return connectionInfo, args.Error(1)
}

func (r *MockGameRepository) FindRunningGameById(id int64) (*models.RunningGame, error) {
	args := r.Called(id)
	return args.Get(0).(*models.RunningGame), args.Error(1)
}

func (r *MockGameRepository) GetConnectionInfo(orakkiId string, playerId int64) (*models.ConnectionInfo, error) {
	args := r.Called(orakkiId, playerId)
	return args.Get(0).(*models.ConnectionInfo), args.Error(1)
}
func TestSignalingUseCaseNewOffer(t *testing.T) {
	mockRepo := new(MockGameRepository)

	mockRepo.On("SaveConnectionInfo", mock.Anything).Return(mock.Anything, nil)

	mockGame := &models.RunningGame{Id: 1, Orakki: &models.Orakki{Id: "orakki1"}}
	mockRepo.On("FindRunningGameById", int64(1)).Return(mockGame, nil)

	// when
	useCase := SignalingUseCase{
		GameRepository: mockRepo,
	}

	const sdpString string = "{\"type\":\"offer\",\"sdp\":\"v=0\\r\\no=- 2040110055720832887 2 IN IP4 127.0.0.1\\r\\ns=-\\r\\nt=0 0\\r\\na=group:BUNDLE 0 1 2\\r\\na=msid-semantic: WMS\\r\\nm=audio 9 UDP/TLS/RTP/SAVPF 111 103 104 9 0 8 106 105 13 110 112 113 126\\r\\nc=IN IP4 0.0.0.0\\r\\na=rtcp:9 IN IP4 0.0.0.0\\r\\na=ice-ufrag:Hlz2\\r\\na=ice-pwd:T3fbJdLjlZ8zi4vUpYl4gQSJ\\r\\na=ice-options:trickle\\r\\na=fingerprint:sha-256 50:46:C7:BA:49:98:B7:0A:CD:13:B7:74:C5:2F:D0:AC:85:AF:43:CF:04:53:0F:0D:0D:71:DF:27:3E:C9:C1:55\\r\\na=setup:actpass\\r\\na=mid:0\\r\\na=extmap:1 urn:ietf:params:rtp-hdrext:ssrc-audio-level\\r\\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\\r\\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\\r\\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\\r\\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\\r\\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\\r\\na=recvonly\\r\\na=rtcp-mux\\r\\na=rtpmap:111 opus/48000/2\\r\\na=rtcp-fb:111 transport-cc\\r\\na=fmtp:111 minptime=10;useinbandfec=1\\r\\na=rtpmap:103 ISAC/16000\\r\\na=rtpmap:104 ISAC/32000\\r\\na=rtpmap:9 G722/8000\\r\\na=rtpmap:0 PCMU/8000\\r\\na=rtpmap:8 PCMA/8000\\r\\na=rtpmap:106 CN/32000\\r\\na=rtpmap:105 CN/16000\\r\\na=rtpmap:13 CN/8000\\r\\na=rtpmap:110 telephone-event/48000\\r\\na=rtpmap:112 telephone-event/32000\\r\\na=rtpmap:113 telephone-event/16000\\r\\na=rtpmap:126 telephone-event/8000\\r\\nm=video 9 UDP/TLS/RTP/SAVPF 96 97 98 99 100 101 102 122 127 121 125 107 108 109 124 120 123 119 114 115 116\\r\\nc=IN IP4 0.0.0.0\\r\\na=rtcp:9 IN IP4 0.0.0.0\\r\\na=ice-ufrag:Hlz2\\r\\na=ice-pwd:T3fbJdLjlZ8zi4vUpYl4gQSJ\\r\\na=ice-options:trickle\\r\\na=fingerprint:sha-256 50:46:C7:BA:49:98:B7:0A:CD:13:B7:74:C5:2F:D0:AC:85:AF:43:CF:04:53:0F:0D:0D:71:DF:27:3E:C9:C1:55\\r\\na=setup:actpass\\r\\na=mid:1\\r\\na=extmap:14 urn:ietf:params:rtp-hdrext:toffset\\r\\na=extmap:2 http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time\\r\\na=extmap:13 urn:3gpp:video-orientation\\r\\na=extmap:3 http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01\\r\\na=extmap:12 http://www.webrtc.org/experiments/rtp-hdrext/playout-delay\\r\\na=extmap:11 http://www.webrtc.org/experiments/rtp-hdrext/video-content-type\\r\\na=extmap:7 http://www.webrtc.org/experiments/rtp-hdrext/video-timing\\r\\na=extmap:8 http://tools.ietf.org/html/draft-ietf-avtext-framemarking-07\\r\\na=extmap:9 http://www.webrtc.org/experiments/rtp-hdrext/color-space\\r\\na=extmap:4 urn:ietf:params:rtp-hdrext:sdes:mid\\r\\na=extmap:5 urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id\\r\\na=extmap:6 urn:ietf:params:rtp-hdrext:sdes:repaired-rtp-stream-id\\r\\na=recvonly\\r\\na=rtcp-mux\\r\\na=rtcp-rsize\\r\\na=rtpmap:96 VP8/90000\\r\\na=rtcp-fb:96 goog-remb\\r\\na=rtcp-fb:96 transport-cc\\r\\na=rtcp-fb:96 ccm fir\\r\\na=rtcp-fb:96 nack\\r\\na=rtcp-fb:96 nack pli\\r\\na=rtpmap:97 rtx/90000\\r\\na=fmtp:97 apt=96\\r\\na=rtpmap:98 VP9/90000\\r\\na=rtcp-fb:98 goog-remb\\r\\na=rtcp-fb:98 transport-cc\\r\\na=rtcp-fb:98 ccm fir\\r\\na=rtcp-fb:98 nack\\r\\na=rtcp-fb:98 nack pli\\r\\na=fmtp:98 profile-id=0\\r\\na=rtpmap:99 rtx/90000\\r\\na=fmtp:99 apt=98\\r\\na=rtpmap:100 VP9/90000\\r\\na=rtcp-fb:100 goog-remb\\r\\na=rtcp-fb:100 transport-cc\\r\\na=rtcp-fb:100 ccm fir\\r\\na=rtcp-fb:100 nack\\r\\na=rtcp-fb:100 nack pli\\r\\na=fmtp:100 profile-id=2\\r\\na=rtpmap:101 rtx/90000\\r\\na=fmtp:101 apt=100\\r\\na=rtpmap:102 H264/90000\\r\\na=rtcp-fb:102 goog-remb\\r\\na=rtcp-fb:102 transport-cc\\r\\na=rtcp-fb:102 ccm fir\\r\\na=rtcp-fb:102 nack\\r\\na=rtcp-fb:102 nack pli\\r\\na=fmtp:102 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42001f\\r\\na=rtpmap:122 rtx/90000\\r\\na=fmtp:122 apt=102\\r\\na=rtpmap:127 H264/90000\\r\\na=rtcp-fb:127 goog-remb\\r\\na=rtcp-fb:127 transport-cc\\r\\na=rtcp-fb:127 ccm fir\\r\\na=rtcp-fb:127 nack\\r\\na=rtcp-fb:127 nack pli\\r\\na=fmtp:127 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42001f\\r\\na=rtpmap:121 rtx/90000\\r\\na=fmtp:121 apt=127\\r\\na=rtpmap:125 H264/90000\\r\\na=rtcp-fb:125 goog-remb\\r\\na=rtcp-fb:125 transport-cc\\r\\na=rtcp-fb:125 ccm fir\\r\\na=rtcp-fb:125 nack\\r\\na=rtcp-fb:125 nack pli\\r\\na=fmtp:125 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=42e01f\\r\\na=rtpmap:107 rtx/90000\\r\\na=fmtp:107 apt=125\\r\\na=rtpmap:108 H264/90000\\r\\na=rtcp-fb:108 goog-remb\\r\\na=rtcp-fb:108 transport-cc\\r\\na=rtcp-fb:108 ccm fir\\r\\na=rtcp-fb:108 nack\\r\\na=rtcp-fb:108 nack pli\\r\\na=fmtp:108 level-asymmetry-allowed=1;packetization-mode=0;profile-level-id=42e01f\\r\\na=rtpmap:109 rtx/90000\\r\\na=fmtp:109 apt=108\\r\\na=rtpmap:124 H264/90000\\r\\na=rtcp-fb:124 goog-remb\\r\\na=rtcp-fb:124 transport-cc\\r\\na=rtcp-fb:124 ccm fir\\r\\na=rtcp-fb:124 nack\\r\\na=rtcp-fb:124 nack pli\\r\\na=fmtp:124 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=4d001f\\r\\na=rtpmap:120 rtx/90000\\r\\na=fmtp:120 apt=124\\r\\na=rtpmap:123 H264/90000\\r\\na=rtcp-fb:123 goog-remb\\r\\na=rtcp-fb:123 transport-cc\\r\\na=rtcp-fb:123 ccm fir\\r\\na=rtcp-fb:123 nack\\r\\na=rtcp-fb:123 nack pli\\r\\na=fmtp:123 level-asymmetry-allowed=1;packetization-mode=1;profile-level-id=64001f\\r\\na=rtpmap:119 rtx/90000\\r\\na=fmtp:119 apt=123\\r\\na=rtpmap:114 red/90000\\r\\na=rtpmap:115 rtx/90000\\r\\na=fmtp:115 apt=114\\r\\na=rtpmap:116 ulpfec/90000\\r\\nm=application 9 UDP/DTLS/SCTP webrtc-datachannel\\r\\nc=IN IP4 0.0.0.0\\r\\na=ice-ufrag:Hlz2\\r\\na=ice-pwd:T3fbJdLjlZ8zi4vUpYl4gQSJ\\r\\na=ice-options:trickle\\r\\na=fingerprint:sha-256 50:46:C7:BA:49:98:B7:0A:CD:13:B7:74:C5:2F:D0:AC:85:AF:43:CF:04:53:0F:0D:0D:71:DF:27:3E:C9:C1:55\\r\\na=setup:actpass\\r\\na=mid:2\\r\\na=sctp-port:5000\\r\\na=max-message-size:262144\\r\\n\"}"
	connectionInfo, err := useCase.NewOffer(1, 1, sdpString)

	// then
	fmt.Print(err)
	assert.NotNil(t, connectionInfo)
	assert.Nil(t, err)
}
