### Create & run Docker image. Создание Docker образа.    

Создание Docker контейнера для gRPC-сервера (build container of server):      

```shell script
docker build -t gw-mtls-service .
```

Развернуть задание с серверным gRPC-приложением:         

```shell script
kubectl apply -f gw-grpc-mtls-service.yaml
```  

## DSN for DBMS  
Code of data exchange with Redis. Обмен данными с Redis. 


