# eagle
Website, SSL certificates, Port and Ping monitoring

## Contentes

1. [Usage](#usage)
    1. [Binaries](#binaries)
        1. [Install](#install)

## Usage

### Binaries

#### Install

1. Clone eagle repo
```bash
$ git clone https://github.com/cloudsark/eagle.git && cd eagle  
```
2. Create eagle database & tables schema (MariaDB 10.x)
```bash
$ mysql -u root -p
MariaDB [(none)]> CREATE DATABASE eagle;
MariaDB [(none)]> CREATE USER 'eagle' IDENTIFIED BY 'type_your_password_here';
MariaDB [(none)]> GRANT ALL ON eagle.* TO 'eagle'@'localhost';
MariaDB [(none)]> use eagle;
MariaDB [(none)]> source db.sql
```
3. Build Docker Image
```bash
$ docker build -t eagle-go .  
```
4. Run Container
```bash
$ docker run --detach --name=eagle-go \ 
                      -e SLACK_TOKEN='' \
                      -e SLACK_CHANNEL='' \
                      -e DB_USER='' \
                      -e DB_PASSWORD='' \ 
                      -e DB_NAME='' \
                      -e DB_HOST='' \
                      -e DB_PORT='' \
                      eagle-go
```
4. Check logs
```bash
$ docker logs eagle-go .  
