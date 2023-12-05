### Тестирование gRPC-сервиса. Test of gRPC-service    
Для тестирования методов gRPC-сервиса выполнить. To test gRPC-service:  

```shell script
go run grpc-service.go
go test -v grpc-service.go  
```

Для тестирования без подключения к gRPC-сервису. To test without conn to gRPC-service:  

```shell script
go test -v service-logic-unit_test.go 
``` 







