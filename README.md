## Тестовый шлюз платежного сервиса Сбербанка. SberPay

### Описание. Description  
Тестовый вебсервер `gosberpay.go` для отладки REST запросов, интернет-эквайринг сервиса Сбербанка.  
Без доступа к боевому серверу Сбербанка, тестирование регистрации заказа осуществляется локально посредством тестового `rest`-запроса с нужными данными.  
Запрос `register` возвращает номер заказа `orderId` и сообщение `Запрос регистрации методом register.do - ОК!` или описание ошибки.    
Запрос `getOrderStatusExtended` возвращает данные статуса заказа или описание ошибки:  
	Order status: {"ErrorCode":"0","ErrorMessage":"Успешно","OrderNumber":"0784sse49d0s134567890","OrderStatus":"6","ActionCode":"-2007","ActionCodeDescription":"Время сессии истекло"}  

URL REST-методов и требования к запросам описаны здесь:
	https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start  
	
TODO разработка всех модулей платежного шлюза.    
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
	go test -v getOrderStatusExtended_test.go  

#### Облако. Cloud.  
	sudo docker run --name gosberpay -p 8443:8443 -d cr.yandex/${REGISTRY_ID}/debian:gosberpay  
	go test -v register_test.go  
	go test -v getOrderStatusExtended_test.go  	
	

### Использование. How use  
	go run gosberpay.go
	go run register.go  
	go run getOrderStatusExtended.go   
	
### Ответ боевого сервера Сбербанка:
#### на запрос регистрации заказа (`register`)с недостоверным ID в запросе (нет регистрации):   
	Status = 200 OK 2023/09/23 15:49:28 
	Response of server:
 	{"errorCode":"5","errorMessage":"Access denied"}
	
#### на запрос состояния заказа (getOrderStatusExtended) нет регистрации:  
	Status = 200 OK 2023/10/10 11:17:08 
	Response of server: 
	{"errorCode":"5","errorMessage":"[userName] or [password] or [token] is empty"}  
	


  




 
