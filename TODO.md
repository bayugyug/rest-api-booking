
##TODO
env DEPNOLOCK=1 dep init -v && go test ./...

##MSYQL

create database restapi;
create user restapi;
grant all privileges on restapi.* to restapi@localhost identified by 'r3stapi';
grant all privileges on restapi.* to restapi@127.0.0.1 identified by 'r3stapi';
flush privileges;





CREATE TABLE IF NOT EXISTS `customers` (
	`id`          int(11) NOT NULL AUTO_INCREMENT,
	`mobile`      varchar(20) NOT NULL,
	`firstname`   varchar(150) DEFAULT NULL,
	`lastname`    varchar(150) DEFAULT NULL,
	`status`      varchar(20), 
	`pass`        varchar(150) DEFAULT NULL,
	
	`latitude`    decimal(20,20),
	`longitude`   decimal(20,20),
	
	`created_dt`  datetime,
	`modified_dt` datetime,
	PRIMARY KEY (`id`),
	UNIQUE  KEY (`mobile`)
) ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS `drivers` (
	`id`          int(11) NOT NULL AUTO_INCREMENT,
	`mobile`      varchar(20) NOT NULL,
	`firstname`   varchar(150) DEFAULT NULL,
	`lastname`    varchar(150) DEFAULT NULL,
	`pass`        varchar(150) DEFAULT NULL,
	
	`latitude`    decimal(20,20),
	`longitude`   decimal(20,20),

	`status`        varchar(20), 
	`vehiclestatus` varchar(20), 
	`created_dt`  datetime,
	`modified_dt` datetime,
  PRIMARY KEY (`id`),
  UNIQUE  KEY (`mobile` )
) ENGINE=InnoDB;


CREATE TABLE IF NOT EXISTS `booking` (
	`id`               int(11) NOT NULL AUTO_INCREMENT,
	`driver_id`        int NOT NULL,
	`customer_id`      int NOT NULL,
	
	`src`              varchar(255),
	`dst`              varchar(255),
	
	`src_latitude`     decimal(20,20),
	`src_longitude`    decimal(20,20),
	
	`dst_latitude`     decimal(20,20),
	`dst_longitude`    decimal(20,20),
	
    `status`           varchar(20), 
	
	`remarks`          varchar(255),
	`remarks_by`       varchar(20),
	
	`pickup_time`      datetime,
	`dropoff_time`     datetime,
	
	`created_dt`  datetime,
	`modified_dt` datetime,
	
    PRIMARY KEY (`id`)
) ENGINE=InnoDB;


##EMPTY Tables
drop table customers;
drop table drivers;
drop table booking;


##GOOGLEMAP_API
https://play.golang.org/p/VsJ42viGuQX

##CURL CMD
curl -v -X GET  'http://127.0.0.1:8989/v1/api/driver/23432432'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/drivers/addresshere'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/customer/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/booking/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/location/driver/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/location/customer/234324'
curl -v -X POST 'http://127.0.0.1:8989/v1/api/login' -d '{"mobile":"6581579058","pass":"dabis","type":"customer"}'



https://github.com/gustavocd/dao-pattern-in-go
