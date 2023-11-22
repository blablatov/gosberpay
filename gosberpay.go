// Тестовый вебсервер для отладки REST запросов, интернет-эквайринг сервиса Сбербанка.
// Использование: go run gosberpay.go
// Выполнить тестовый запрос go test -v register_test.go или go test -v getOrderStatusExtended.go
// URL-адреса для доступа к запросам REST описаны здесь:
// https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

var (
	crtFile = filepath.Join(".", "certs", "server.crt")
	keyFile = filepath.Join(".", "certs", "server.key")
)

func main() {
	log.SetPrefix("Event main: ")
	log.SetFlags(log.Lshortfile)

	// TLS or simple connect. Подключение TLS или базовое
	http.HandleFunc("/", handle)
	http.ListenAndServeTLS("localhost:8443", crtFile, keyFile, nil)
	//http.ListenAndServe("localhost:8088", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nRequest of client:\n")

	// Параметры заголовков http-запроса. Parameters of headers
	fmt.Fprintf(w, "Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)
	fmt.Printf("Method = %s\nURL = %s\nProto = %s\n", r.Method, r.URL, r.Proto)

	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		fmt.Printf("Header[%q] = %q\n", k, v)
	}

	fmt.Fprintf(w, "Host = %q\n", r.Host)
	fmt.Printf("Host = %q\n", r.Host)

	fmt.Fprintf(w, "RemoteAddr = %q\n", r.RemoteAddr)
	fmt.Printf("RemoteAddr = %q\n", r.RemoteAddr)

	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		fmt.Printf("Form[%q] = %q\n", k, v)
	}
	// Код обработки запросов. Code of processing of requests
	switch r.URL.Path {
	case "/register": // Запрос регистрации заказа register.do.
		for k, v := range r.Form {
			if k != "" || v != nil {
				fmt.Printf("%s, %s\n", k, v)
			}
		}
		orderId := "70906e55-7114-41d6-8332-4609dc6590f4" // Возвращаемый ID заказа. ID of order
		fmt.Fprintf(w, " orderId: %v", orderId)

		formUrl := "https://3dsec.sberbank.ru/payment/merchants/test/payment_ru.html?mdOrder=" + orderId // URL платёжной формы
		fmt.Fprintf(w, "\n formUrl: %v", formUrl)

	case "/getOrderStatusExtended": //Запрос состояния заказа (getOrderStatusExtended.do)
		for k, v := range r.Form {
			if k != "" || v != nil {
				fmt.Printf("%s, %s\n", k, v)
			}
		}

		// Возвращаемые тестовые json параметры запроса состояния заказа
		res, _ := json.Marshal(struct {
			ErrorCode             string `json: "errorCode"`
			ErrorMessage          string `json: "errorMessage"`
			OrderNumber           string `json: "orderNumber"`
			OrderStatus           string `json: "orderStatus"`
			ActionCode            string `json: "actionCode"`
			ActionCodeDescription string `json: "actionCodeDescription"`
		}{
			ErrorCode:             "0",
			ErrorMessage:          "Успешно",
			OrderNumber:           "0784sse49d0s134567890",
			OrderStatus:           "6",
			ActionCode:            "-2007",
			ActionCodeDescription: "Время сессии истекло",
		})
		fmt.Fprintf(w, "Order status: %s", res)

	default:
	}
}
