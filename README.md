# go-eagle
Website, SSL certificates, Port and Ping monitoring

## Contentes

1. [Important](#important)
1. [Usage](#usage)
    1. [Database](#database)
        1. [MySQL](#mysql)
        2. [MongoDB](#mongo)
    1. [Binaries](#binaries)
        1. [Install](#install)
             1. [Go-eagle](#go-eagle)
                  1. [Standalone](#standalone)
                  2. [Docker](#docker)
             1. [Client](#client)
1. [Support](#support-us)

## Important
For now "go-eagle" use both MySQL & MongoDB as data storage but MySQL will be deprecated soon !

## Usage

### Database

#### MySQL
1. Download database file
```
wget https://raw.githubusercontent.com/cloudsark/go-eagle/master/db.sql
```
2. Create eagle database & tables schema (MariaDB 10.x)
```bash
$ mysql -u root -p
MySQL [(none)]> CREATE DATABASE eagle;
MySQL [(none)]> CREATE USER 'eagle' IDENTIFIED BY 'type_your_password_here';
MySQL [(none)]> GRANT ALL ON eagle.* TO 'eagle'@'localhost';
MySQL [(none)]> use eagle;
MySQL [(none)]> source db.sql
```

#### MongoDB
Note: for now MongoDB store disks stat only
1. Create database from cli
```
use eagle
```
2. Create disks collection and add sample data
```
db.eagle.insert({"_id":{"$oid":"5fe514f0ca8b569e25d853d7"},"hostname":"example.com","name":"sda2","path":"/","fstype":"xfs","total":"40940.00","free":"40556.00","used":"384.00","percent":{"$numberDouble":"0.94"},"flag":{"$numberInt":"1"}})
```
3. Create username and password
```
use eagle
db.createUser(
  {
    user: "your_username",
    pwd: "your_password",
    roles: [ { role: "readWrite", db: "eagle" } ]
  }
)
```
Note: you can use mongodb cloud as your data store
https://www.mongodb.com/pricing

### Binaries

#### Install

##### Go-eagle

###### Standalone

1. Clone eagle repo
```bash
$ git clone https://github.com/cloudsark/go-eagle.git && cd go-eagle  
```
1. Build go-eagle
```
go build
```
1. Create go-eagle service
```
cp systemd/go-eagle.service /etc/systemd/system/go-eagle.service
vi go-eagle.service
ADD env variables to Environment=

save file
systemctl daemon-reload
systemctl enable go-eagle
systemctl start go-eagle
systemctl status go-eagle
```

1. Validate
```
tail -f /var/log/eagle-log.log
```

###### Doekcer

1. Clone eagle repo
```bash
$ git clone https://github.com/cloudsark/eagle.git && cd eagle  
```
1. Build Docker Image
```bash
$ docker build -t eagle-go .  
```
1. Run Container
```bash
$ docker run --detach --name=eagle-go \ 
                      -e SLACK_TOKEN='' \
                      -e SLACK_CHANNEL='' \
                      -e DB_USER='' \
                      -e DB_PASSWORD='' \ 
                      -e DB_NAME='' \
                      -e DB_HOST='' \
                      -e DB_PORT='' \
                      -e CLIENT_USERNAME='' \
                      -e CLIENT_PASSWORD='' \
                      -e MONGO_DB='' \
                      -e MONGO_USER='' \
                      -e MONGO_PASSWORD='' \
                      -e MONGO_URL='' \
                      -e MONGO_PARN='' \
                      eagle-go
```
1. Check logs
```bash
$ docker logs eagle-go .  

##### Client
you need to install eagle-client on any server you want to monitor its (cpu & disks)

1. Download eagle-client
```
wget https://github.com/cloudsark/go-eagle/releases/download/v0.2/eagle-client
```
1. Create client directory
```
mkdir /eagle-client
mv eagle-client /eagle-client
chmod 777 /eagle-client/eagle-client
```
1. Create eagle-client service
```
cd /etc/systemd/system/
wget https://raw.githubusercontent.com/cloudsark/go-eagle/master/systemd/go-client.service

vi go-client.service
set client username & password in Environment=
modify both (WorkingDirectory & ExecStart)
save file
systemctl daemon-reload
systemctl enable eagle-client
systemctl start eagle-client
systemctl status eagle-client
```

1. Validate client
```
netstat -ntlp | grep :10052
```
1. Allow access to port 10052 from firewall or security group

## Support-us
If this project help you, you can support us on paypal :) 

[![paypal](https://linuxdirection.s3-eu-west-1.amazonaws.com/support-paypal-2.png)](https://paypal.me/cloudsark?locale.x=en_US)