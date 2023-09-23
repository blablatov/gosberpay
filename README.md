## Тестовый шлюз платежного сервиса Сбербанка. SberPay

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
	Status = 200 OK 2023/09/23 15:49:28 
	Response of server:
 	{"errorCode":"5","errorMessage":"Access denied"}
   

  




 
