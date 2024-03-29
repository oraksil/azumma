# Azumma

Azumma is a REST API server that provides player creation, games catalog and WebRTC signaling proxy. In a cluster mode, it's responsible to provision Orakki instance on demand. It is built with Golang.

# Setup (Prerequites)

## Database

### Run `mariadb`
```
$ docker run -d \
    --name oraksil-db \
    -v $PWD/temp/data:/var/lib/mysql \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=<root-password> \
    mariadb:10.5.4
```

### Connect to shell and initialize db
```
$ docker exec -it oraksil-db bash
# mysql -u root -p
...
...

MariaDB > create database oraksil;
Query OK, 1 row affected (0.001 sec)

MariaDB > grant all privileges on oraksil.* TO 'oraksil'@'%' identified by '<oraksil-password>';
Query OK, 0 rows affected (0.014 sec)
```

### Migrations

#### Install `golang-migrate`
https://github.com/golang-migrate/migrate

```
$ curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
$ sudo apt-get install -y migrate
```

#### Initial migration
```
$ migrate -path ./migrations -database "mysql://oraksil:<oraksil-password>@(localhost:3306)/oraksil" up
```

#### Migrate to specific version
```
$ migrate -path ./migrations -database "mysql://oraksil:<oraksil-password>@(localhost:3306)/oraksil" up 001
```


## Message Queue (RabbitMQ)

### Run `rabbitmq`
```
$ docker run -d \
    --name oraksil-mq \
    --hostname oraksil-mq \
    -v $PWD/temp/mq:/var/lib/rabbitmq \
    -e RABBITMQ_DEFAULT_USER=oraksil \
    -e RABBITMQ_DEFAULT_PASS=oraksil \
    -p 5672:5672 \
    -p 15672:15672 \
    rabbitmq:3.8.5-management
```

### Declare `exchanges`
```
$ docker exec -it oraksil-mq bash
root@oraksil-mq:/# rabbitmqadmin -u oraksil -p oraksil declare exchange name=mqrpc.oraksil.p2p type=direct
exchange declared

root@oraksil-mq:/# rabbitmqadmin -u oraksil -p oraksil declare exchange name=mqrpc.oraksil.broadcast type=fanout
exchange declared

root@oraksil-mq:/# rabbitmqadmin -u oraksil -p oraksil list exchanges
+---------------------------+---------+
|         name              |  type   |
+---------------------------+---------+
|                           | direct  |
| amqrpc.direct             | direct  |
| amqrpc.fanout             | fanout  |
| amqrpc.headers            | headers |
| amqrpc.match              | headers |
| amqrpc.rabbitmqrpc.trace  | topic   |
| amqrpc.topic              | topic   |
| mqrpc.oraksil.broadcast   | fanout  |
| mqrpc.oraksil.p2p         | direct  |
+---------------------------+---------+
```

# Run

```bash
$ go run cmd/app.go
```