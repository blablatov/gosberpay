### Generate Server and Client side code   
Go to ``Go`` module directory location `mtls-grpc-gateway/gw-mtls-proto` and execute the following shell commands: 

	protoc servicepay.proto --go_out=plugins=grpc:. !!!
```
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--go_out=plugins=grpc:. \
product_info.proto
```

Go to ``Go`` module directory location `mtls-grpc-gateway/gw-mtls-gate` and execute the following shell commands:   

	protoc servicepay.proto --grpc-gateway_out=logtostderr=true:. !!!
```
protoc -I/usr/local/include -I. \
-I$GOPATH/src \
-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
--grpc-gateway_out=logtostderr=true:. \
product_info.proto
```