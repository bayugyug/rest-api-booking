## rest-api-booking

* [x] Sample golang rest-api that simulates a simple car-booking


### Pre-Requisite
	
	- Please run this in your command line to ensure packages are in-place.
	  (normally these will be handled when compiling the api binary)
	
```sh

		go get -u -v github.com/go-chi/chi
		go get -u -v github.com/go-chi/chi/middleware
		go get -u -v github.com/go-chi/cors
		go get -u -v github.com/go-chi/jwtauth
		go get -u -v github.com/go-chi/render
		go get -u -v github.com/dgrijalva/jwt-go
		go get -u -v github.com/go-sql-driver/mysql

```

### Compile

```sh

     git clone https://github.com/bayugyug/rest-api-booking.git && cd rest-api-booking

     git pull && make clean && make

```

### Required Data Preparation

    - Create sample mysql db
	
	- Needs to create the api database and grant the necessary permissions
	
	- Refer the testdata/dump.sql
	
```sh

	create database restapi;
	create user restapi;
	grant all privileges on restapi.* to restapi@localhost identified by 'xxxx';
	grant all privileges on restapi.* to restapi@127.0.0.1 identified by 'xxxx';
	flush privileges;

```

### List of End-Points-Url


```sh

#Customer Create
curl -v -X POST 'http://127.0.0.1:8989/v1/api/customer'  -d '{
					"mobile":"6581577001",
					"pass":"8888",
					"latitude":1.304832,
					"longitude":103.852844,
					"firstname":"customer",
					"lastname": "dabis"
					}'
#Customer OTP
curl -v -X PUT 'http://127.0.0.1:8989/v1/api/otp'     -d '{"mobile":"6581577001","otp":"07814","type":"customer"}'

#Customer Login			
curl -v -X POST 'http://127.0.0.1:8989/v1/api/login'    -d '{"mobile":"6581577001","pass":"8888","type":"customer"}'

#Customer Info	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}"  -X GET 'http://127.0.0.1:8989/v1/api/customer/6581577001' 

#Customer Update Password
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/password/customer'   -d '{"mobile":"6581577001","pass":"1234"}'

#Customer Update GPS Coordinates	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/location'   -d '{"mobile":"6581577001","type":"customer","latitude":1.35821,"longitude":103.85615}'

#Customer Update
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/customer'  -d '{
					"mobile":"6581579000",
					"latitude":1.304832,
					"longitude":103.852855,
					"firstname":"customer",
					"lastname": "dabis"
					}'
#Customer Delete	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X DELETE 'http://127.0.0.1:8989/v1/api/customer/6581579000' 

#Customer Update Status		
curl -v -H "Authorization: BEARER " -X PUT 'http://127.0.0.1:8989/v1/api/status/customer/6581579001/active'

#Driver Create
curl -v -X POST 'http://127.0.0.1:8989/v1/api/driver'  -d '{
					"mobile":"6581755001",
					"pass":"8888",
					"latitude":1.304832,
					"longitude":103.852844,
					"firstname":"driver",
					"lastname": "dabis"
					}'
#Driver OTP
curl -v -X PUT 'http://127.0.0.1:8989/v1/api/otp' -d '{"mobile":"6581755001","otp":"03790","type":"driver"}'

#Driver Login			
curl -v -X POST 'http://127.0.0.1:8989/v1/api/login'    -d '{"mobile":"6581755001","pass":"8888","type":"driver"}'

#Driver Info	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}"  -X GET 'http://127.0.0.1:8989/v1/api/driver/6581755001' 

#Driver Update GPS Coordinates	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/location'   -d '{"mobile":"6581755001","type":"driver","latitude":1.35991,"longitude":102.85615}'


#Driver Update
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/driver'  -d '{
					"mobile":"6581755001",
					"latitude":1.304832,
					"longitude":103.852855,
					"firstname":"driver",
					"lastname": "dabis"
					}'
#Driver Delete	
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X DELETE 'http://127.0.0.1:8989/v1/api/driver/6581755001'

#Driver Update Status		
curl -v -H "Authorization: BEARER " -X PUT 'http://127.0.0.1:8989/v1/api/status/driver/6581755001/active' 

#Driver Update Password
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/password/driver'   -d '{"mobile":"6581755001","pass":"1234"}'

#Driver List Within Nearest 50 KM Radius /drivers/{LATITUDE}/{LONGITUDE}
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}"  -X GET 'http://127.0.0.1:8989/v1/api/drivers/1.336209/103.737326'     

#Driver Update Vehicle Status
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/vehiclestatus'   -d '{"mobile":"6581755001","status":"canceled","latitude":1.35991,"longitude":102.85615}'


#Booking Create
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X POST 'http://127.0.0.1:8989/v1/api/booking'  -d '{
					"mobile_customer":"6581579000",
					"src":"kembangan",
					"src_latitude":1.371572,
					"src_longitude":103.956551,
					"mobile_driver":"6581755001",
					"dst":"bugis",
					"dst_latitude":1.371572,
					"dst_longitude":103.956551
					}'

#Booking Info
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X GET 'http://127.0.0.1:8989/v1/api/booking/4'   

#Booking Update Pickup Time  & Set Status to Trip-Start
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/booking/pickup-time/4'   

#Booking Update Dropoff Time & Set Status to Trip-End
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/booking/dropoff-time/4'   

#Booking Update Status by Customer (canceled)
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/booking/status/customer/4'   

#Booking Update Status by Driver  (canceled/gas-up/panic)
curl -v -H "Authorization: BEARER {TOKEN_FROM_LOGIN}" -X PUT 'http://127.0.0.1:8989/v1/api/booking/status/driver/4/{canceled|gas-up|panic}' 




```


### Mini-How-To on running the api binary

    [x] The api can accept a json format configuration
	
	[x] Fields:
	
		- port    = port to run the http server (default: 8989)
		
		- driver  = database details for mysql  (user/pass/dbname/host/port)
		
		- showlog = flag for dev't log on std-out
		
		
```sh

	 ./rest-api-booking --config '{"port":"8989",
		"driver":{
			"user":"restapi",
			"pass":"xxxxx",
			"port":3306,
			"name":"restapi",
			"host":"127.0.0.1"},
		"showlog":true}'


```

### Notes

        [x] The formula in determining the nearest location GPS coordinates

                - Earth radius is 3959 in miles

                - Earth radius is 6371 in km

        [x] Api is using the nearest 50KM as default



### Reference
	
	- GPS Coordinates Generator
	
		[x] [www.latlong.net](https://www.latlong.net/)

		[x] [www.mapdevelopers.com](https://www.mapdevelopers.com/draw-circle-tool.php)


### License

[MIT](https://bayugyug.mit-license.org/)

