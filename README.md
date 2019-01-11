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
curl -v -X GET  'http://127.0.0.1:8989/v1/api/driver/6581579999'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/drivers/addresshere'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/customer/6581578888'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/booking/2'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/location/{driver|customer}/address'
curl -v -X GET  'http://127.0.0.1:8989/v1/api/address'  -d '{"address":"200 Victoria Street Bugis Junction Singapore"}'
curl -v -X POST 'http://127.0.0.1:8989/v1/api/login'    -d '{"mobile":"6581578888","pass":"dabis","type":"customer"}'


curl -v -X GET  -H "Authorization: BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDcyNDg1MDQsIm1vYmlsZSI6IjY1ODE1NzkwNTgifQ.vMpIOmMZXsaWtu4sQj28SoB-SyS6qxZCjD0ikoOyuTU" 'http://127.0.0.1:8989/v1/api/location/address'
```


### Reference




### License

[MIT](https://bayugyug.mit-license.org/)

