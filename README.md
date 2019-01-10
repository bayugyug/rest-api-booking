## rest-api-booking

* [x] Sample golang rest-api that simulates a simple car-booking


### Prerequisite

```sh

    go get -u -v github.com/go-chi/chi
    go get -u -v github.com/go-chi/chi/middleware
    go get -u -v github.com/go-chi/cors
    go get -u -v github.com/go-chi/jwtauth
    go get -u -v github.com/go-chi/render

```

### Compile

```sh

     git clone https://github.com/bayugyug/rest-api-booking.git && cd rest-api-booking

     git pull && make clean && make

```

### Required data preparation

    - Create sample mysql db (refer the testdata/dump.sql)


### List of end-points-url

```sh
curl -v -X GET  'http://127.0.0.1:8989/v1/api/driver/23432432'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/drivers/addresshere'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/customer/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/booking/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/location/driver/234324'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/location/customer/234324'
curl -v -X POST 'http://127.0.0.1:8989/v1/api/login' -d '{"mobile":"6581579058","pass":"dabis","type":"customer"}'
```


### Reference




### License

[MIT](https://bayugyug.mit-license.org/)

