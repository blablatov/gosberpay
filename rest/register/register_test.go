// Запрос регистрации заказа (register.do)
// Выполнить тестовый запрос go test -v register_test.go
// URL-адреса для доступа к запросам REST описаны здесь:
// https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start

package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestStrings(t *testing.T) {
	var tests = []struct {
		userName  string
		password  string
		amount    string
		returnUrl string
	}{
		{"username-api", "password", "99999", "https://test.ru"},
		{"-api", "qwerty", "amount000001001abcdef", ",#&U*(()))_+_11234"},
		{"\n\t", "' '", "-1-2---==++", "//http///tcp///udp/"},
	}

	var prevuserName string
	for _, test := range tests {
		if test.userName != prevuserName {
			fmt.Printf("\n%s\n", test.userName)
			prevuserName = test.userName
		}
	}

	var prevpassword string
	for _, test := range tests {
		if test.password != prevpassword {
			fmt.Printf("\n%s\n", test.password)
			prevpassword = test.password
		}
	}

	var prevamount string
	for _, test := range tests {
		if test.amount != prevamount {
			fmt.Printf("\n%s\n", test.amount)
			prevamount = test.amount
		}
	}

	var prevreturnUrl string
	for _, test := range tests {
		if test.returnUrl != prevreturnUrl {
			fmt.Printf("\n%s\n", test.returnUrl)
			prevreturnUrl = test.returnUrl
		}
	}
}

func TestFormRequst(t *testing.T) {
	// URL для проверки ответа боевого сервиса Сбера
	//apiUrl := "https://3dsec.sberbank.ru/payment/rest/register.do"

	// URL тестового сорвера локально. Для облака указать внешний IP ВМ.
	apiUrl := "https://localhost:8443/register"
	//Формирование параметров запроса. Params of request
	params := url.Values{"userName": {"username-api"}}
	params.Set("password", "password")
	params.Set("amount", "99999")
	params.Set("returnUrl", "https://test.ru")

	// Формирование строки параметров запроса. String data of request
	dreq := strings.NewReader(params.Encode())

	// Подгрузка сертификата и ключа. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Формирование метаданных структуры запроса. Struct of request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter:       w,
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			},
		},
	}

	// Форматирование запроса. Formatting of the request
	req, _ := http.NewRequest(http.MethodPost, apiUrl, dreq)
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req) // Выполнение запроса. Send of request
	if err != nil {
		log.Fatalf("Ошибка параметров запроса: %v", err)
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	// Defer to finished the method and got response
	defer resp.Body.Close()

	fmt.Printf("Status = %v ", resp.Status) // Статус ответа сервера. Status of response

	// Чтение данных сервера, обработка ошибок. Reads data from server, check errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения данных сервера: %v", err)
	}
	log.Println("\nResponse of server:\n", string([]byte(body)))

	// Проверка идентификатора заказа. Checks order ID
	res := (string)([]byte(body))
	if strings.Contains(res, "orderId") || strings.Contains(res, "formUrl") {
		fmt.Println("Запрос регистрации методом register.do - ОК!")
	}
}
