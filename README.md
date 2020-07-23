## Setup

## Test

```
$ go test ./...
```

## Database

### Run `mariadb`
```
$ docker run -d \
    --name oraksil-azuma \
    -v $PWD/temp/data:/var/lib/mysql \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=<root-password> \
    mariadb:10.5.4
```

### Connect to shell and initialize db
```
$ docker exec -it oraksil-azuma bash
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
$ curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
$ sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
$ sudo apt-get update
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
