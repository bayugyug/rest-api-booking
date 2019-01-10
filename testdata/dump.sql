

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
	`pass`        varchar(150) DEFAULT NULL,
	`status`      varchar(20), 
	`logged`      int(1) default 0, 
	 
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
	`logged`      int(1) default 0, 
	
	`latitude`    decimal(20,20),
	`longitude`   decimal(20,20),

	`status`        varchar(20), 
	`vehiclestatus` varchar(20), 
	`created_dt`  datetime,
	`modified_dt` datetime,
  PRIMARY KEY (`id`),
  UNIQUE  KEY (`mobile` )
) ENGINE=InnoDB;

 

CREATE TABLE IF NOT EXISTS `bookings` (
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


INSERT INTO customers (mobile,pass,status,firstname,lastname,created_dt,modified_dt)
VALUES ('6581579058',md5('dabis'),'active','bayugyug','hehehe',now(),now());

drop table bookings;
drop table drivers;
drop table customers;



