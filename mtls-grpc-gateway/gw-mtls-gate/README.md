### Тестирование функциональности gRPC прокси-сервера. Testing code with conn to gRPC-gateway          
  
Для тестирования удаленных методов gRPC-сервиса.  
Перед запуском gRPC-шлюза `reverse-proxy-server` необходимо запустить grpc-сервер `grpc-service` и тестовый веб-сервер `gosberpay`.   
Before starting gRPC-gateway `reverse-proxy-server`, starting gRPC-server `grpc-service` and the test Web-server `gosberpay`:       

```shell script
go run gosberpay.go
go run grpc-service.go

go run reverse-proxy-server.go

go test -v get_test.go
go test -v post_test.go
```

Для тестирования gRPC-шлюза, без подключения к gRPC-серверу. To test gRPC-gateway, without conn to rpc server:     
       
```shell script
go test -v reverse-proxy-server-unit_test.go
```  





