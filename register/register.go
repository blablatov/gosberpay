// Регистрация заказа с помощью метода register.do
// Выполнить тестовый запрос запустив модуль go run register.go
// URL REST-методов и требования к запросам описаны здесь:
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
	"path/filepath"
	"strings"
)

var (
	crtFile = filepath.Join(".", "certs", "client.crt")
	keyFile = filepath.Join(".", "certs", "client.key")
)

func main() {

	apiUrl := "https://localhost:8443/register"
	//Формирование параметров rest запроса. Params of request
	params := url.Values{"userName": {"username-api"}}
	params.Set("password", "password")
	params.Set("amount", "99999")
	params.Set("returnUrl", "https://test.ru")

	// Формирование строки параметров запроса. String data of request
	dreq := strings.NewReader(params.Encode())

	// Подгрузка сертификатов. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Формирование параметров структуры запроса. Struct of request
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
