// Тестовый вебсервер на примере метода register.do, интернет-эквайринг сервиса Сбербанка.
// Использование: $ go run gosberpay.go
// Выполнить тестовый запрос запустив модуль go test -v register_test.go
// URL REST-методов и требования к запросам описаны здесь:
// https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start

package main

import (
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

	// TLS or simple connect. Подключение по протоколу TLS или базовое
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
		fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
		fmt.Printf("Form[%q] = %q\n", k, v)
	}
	// Блок методов обработки запросов. Code of processing methods of requests
	switch r.URL.Path {
	case "/register": // Метод регистрации заказа register.do. Test method
		for k, v := range r.Form {
			if k != "" || v != nil {
				fmt.Printf("%s, %s\n", k, v)
			}
		}

		orderId := "70906e55-7114-41d6-8332-4609dc6590f4" // ID заказа. ID of order
		r.Form.Add("orderId", orderId)

		// Ответ сервера с параметром ID заказа. Response with ID the order
		for k, v := range r.Form {
			fmt.Fprintf(w, "Form[%q] = %q\n", k, v)
		}
	default:
	}
}
