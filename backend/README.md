## Test

```
$ go test ./...
```

## Database

### Run `mariadb`
```
$ docker run -d --name oraksil-azuma -v $PWD/temp/data:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=<root-password> mariadb:10.5.4
```

### Connect to shell and initialize db
```
$ docker run -it oraksil-azuma bash
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