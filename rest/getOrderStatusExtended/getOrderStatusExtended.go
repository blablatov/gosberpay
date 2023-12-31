// Расширенный запрос состояния заказа (getOrderStatusExtended.do)
// Выполнить запрос go run getOrderStatusExtended.go
// URL-адреса для доступа к запросам REST описаны здесь:
// https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start

package getOrderStatusExtended

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	rp "github.com/blablatov/gosberpay/rest/register"
)

type StatusParam struct {
	OrderId string `json:"orderId"` // Номер заказа в платежной системе
	rp.ParamsPay
}

func (sp StatusParam) OrderStatusExtended(sch chan string, crtFile, keyFile string) {

	// URL для проверки ответа боевого сервиса Сбера
	//apiUrl := "https://3dsec.sberbank.ru/payment/rest/getOrderStatusExtended.do"

	// URL локального тестового сервера. Для облака указать внешний IP ВМ
	apiUrl := "https://localhost:8443/getOrderStatusExtended"

	// Формирование json параметров запроса. JSON params of request
	payload, _ := json.Marshal(struct {
		OrderId string `json:"orderId"`
	}{
		OrderId: sp.OrderId,
	})

	// Подгрузка сертификата и ключа. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}
	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Форматирование запроса. Formatting of the request
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	resp, err := client.Do(req) // Выполнение запроса. Send of request
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	// Отложеное выполнение закрытия запроса, до выполнения метода и получения ответа
	// Defer to finished the method and got response
	defer resp.Body.Close()

	fmt.Printf("Status = %v ", resp.Status) // Статус ответа сервера. Status of response

	// Чтение данных сервера, обработка ошибок. Reads data from server, check errors
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println("\nResponse server: \n", string([]byte(body)))

	rb := string([]byte(body))

	data, err := json.MarshalIndent(rb, "=", " ")
	if err != nil {
		log.Fatalf("Сбой маршалинга JSON: %s", err)
	}
	fmt.Printf("Data status response = %s\n", data)

	sch <- rb
}
