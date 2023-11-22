## Service and Client of gRPC via mTLS. gRPC-gateway. Rest. Redis  

### Building and Running Service

In order to build, Go to ``Go`` module root directory location (mtls-grpc-gateway/gw-mtls-service) and execute the following
 shell command,
```
go build -v 
```  
```
./gw-mtls-service
```

### Building and Running Client   

In order to build, Go to ``Go`` module root directory location (mtls-grpc-gateway/gw-mtls-client) and execute the following
 shell command,
```
go build -v
```  
```
./gw-mtls-client
```

### Testing

* Add a new product to the ProductInfo service:  
```
$ curl -X POST https://localhost:8443/v1/product -d '{"name": "Apple", "description": "iphone7", "price": 699}'
```
* Response:   
```
"38e13578-d91e-11e9-819f-6c96cfe0687d"
```

* Or via SoapUI and etc.:  
```
Mon Mar 06 21:24:58 YEKT 2023: DEBUG: http-outgoing >> 
POST /v1/product HTTP/1.1
Accept-Encoding: gzip,deflate
Content-Type: application/json
Content-Length: 57
Host: localhost:8443
Connection: Keep-Alive
User-Agent: Apache-HttpClient/4.5.5 (Java/16.0.1)

{"name": "Apple", "description": "iphone7", "price": 699}

Mon Mar 06 21:24:58 YEKT 2023: DEBUG: http-incoming << 
HTTP/1.1 200 OK
Content-Type: application/json
Grpc-Metadata-Content-Type: application/grpc
Date: Mon, 06 Mar 2023 16:24:58 GMT
Content-Length: 38

"ce01618f-ec7b-4c7f-85ea-0e7841363e59"
```

* Get the existing product using ProductID:  
```
$ curl http://localhost:8443/v1/product/38e13578-d91e-11e9-819f-6c96cfe0687d
```  
* Response: 
```
{"id":"38e13578-d91e-11e9-819f-6c96cfe0687d","name":"Apple","description":"iphone7","price":699}
```  

* Or via request string of web-browser:   
```
https://localhost:8443/v1/product/38e13578-d91e-11e9-819f-6c96cfe0687d  
```  
* Response: 
```
{"id":"38e13578-d91e-11e9-819f-6c96cfe0687d","name":"Apple","description":"iphone7","price":699}
```   

## Additional Information

### Generate Server and Client side code 
``` 
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. \
product_info.proto
```

### Update after changing the service definition
``` 
go get -u github.com/blablatov/mtls-grpc-gateway/gw-mtls-proto
```

### Generate reverse proxy service code
```
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true:. \
product_info.proto
```

### Update after changing the reverse proxy service definition
``` 
go get -u github.com/blablatov/mtls-grpc-gateway/gw-mtls-gate
```

### Generate the swagger file correspond to reverse proxy service
```
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--swagger_out=logtostderr=true:. \
product_info.proto
```

