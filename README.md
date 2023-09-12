## Тестовый шлюз платежного сервиса Сбербанка. SberPay

![Payment page](https://github.com/blablatov/gosberpay/raw/master/pagepay.png) 

### Описание. Description  
Тестовый вебсервер `gosberpay.go` на примере метода `register.do`, интернет-эквайринг сервиса Сбербанка.  
За неимением доступа к боевому серверу Сбербанка, тестирование регистрации заказа осуществляется с помощью тестового rest-запроса, с данными указаными для интерфейса REST.  
Запрос возвращает сообщение - `Запрос регистрации методом register.do - ОК!` или описание ошибки.

URL REST-методов и требования к запросам описаны здесь:
	https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start  
	
TODO (private repo) разработка всех модулей шлюза для проведения платежей.  
также для Alfa API  
	https://developers.alfabank.ru/products/alfa-api/documentation/development/specification/introduction  
	
![Gateway](https://github.com/blablatov/gosberpay/raw/master/gateway.png)


### Сборка локально и в Yandex Cloud. Build local and to Yandex Cloud  
#### Локально. Local:  
	docker build -t gosberpay -f Dockerfile  
	
#### Облако. Cloud.  
	sudo docker build . -t cr.yandex/${REGISTRY_ID}/debian:gosberpay -f Dockerfile


### Тестирование локально и в Yandex Cloud. Testing local and to Yandex Cloud       
#### Локально. Local:    
	go test -v register_test.go    

#### Облако. Cloud.  
	sudo docker run --name gosberpay -p 8443:8443 -d cr.yandex/${REGISTRY_ID}/debian:gosberpay  
	go test -v register_test.go  
	

### Использование. How use  
	go run gosberpay.go
	go test -v register_test.go  
	
### Ответ боевого сервера Сбербанка, с недостоверным ID в запросе (нет регистрации)  
	CLIENT_RANDOM 8725e2b3aa5beb71134b1e84f760047e19f5c09082b18f7c6904a19a2f7abf68 d078394d43b2d5c8e413cafccb41a87ae03dff8b9705754a0635e4ed76fc7aa346bb286df9b92406bf524900da2e2ede 
	Status = 405 Not Allowed 2023/09/12 15:25:07 
	Response of server:
	<html>
	<head><title>405 Not Allowed</title></head>
	<body>
	<center><h1>405 Not Allowed</h1></center>
	<hr><center>nginx</center>
	</body>
	</html> 
   

  




 
