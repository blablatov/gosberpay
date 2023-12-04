### Тестирование функциональность gRPC прокси-сервера. 
### Testing code with conn to gRPC-server          
  
Для проверки удаленных методов сервиса.  
Перед выполнением `reverse-proxy-server` запустить grpc-сервер `grpc-service` и тестовый вебсервер `gosberpay`.   
(Conventional test that starts a gRPC client test the service with RPC.Before his execute run grpc-server:       

```shell script
go run gosberpay.go
go run grpc-service.go

go run reverse-proxy-server.go

go test -v get_test.go
go test -v post_test.go
```

Для тестирования gRPC-шлюза, без подключения к gRPC-серверу:     
       
```shell script
go test -v reverse-proxy-server-unit_test.go
```  





