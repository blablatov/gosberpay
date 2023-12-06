### Generate gRPC-service and gRPC-gateway code   
Генерация кода сервиса и обратного прокси-сервера из proto-файла: 

```shell script
protoc servicepay.proto --go_out=plugins=grpc:.
```
```shell script
protoc servicepay.proto --grpc-gateway_out=logtostderr=true:.
```