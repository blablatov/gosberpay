// Регистрация заказа с помощью метода register.do
// Выполнить тестовый запрос запустив модуль go run register.go
// URL REST-методов и требования к запросам описаны здесь:
// https://securepayments.sberbank.ru/wiki/doku.php/integration:api:rest:start

package register

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	//"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type ParamsPay struct {
	paymentToken         string `json:"paymentToken"`         // Параметры userName и pаssword передавать не нужно.
	Amount               string `json:"amount"`               // Сумма платежа в минимальных единицах валюты
	currency             string `json:"currency"`             // Код валюты платежа ISO 4217
	language             string `json:"language"`             // Язык в кодировке ISO 639-1
	orderNumber          string `json:"orderNumber"`          // номер заказа в системе магазина
	ReturnUrl            string `json:"returnUrl"`            // Адрес перенаправления в случае успешной оплаты
	jsonParams           string `json:"jsonParams"`           // Дополнительные параметры запроса. Формат: {«Имя1»: «Значение1», «Имя2»: «Значение2»}
	pageView             string `json:"pageView"`             // Какие страницы платёжного интерфейса должны загружаться для клиента:DESKTOP,MOBILE
	expirationDate       string `json:"expirationDate"`       // Дата и время окончания жизни заказа
	merchantLogin        string `json:"merchantLogin"`        // При заказе от имени дочернего продавца, его логин в этом параметре
	features             string `json:"features"`             // AUTO_PAYMENT-платёж проводится без проверки подлинности владельца карты
	orderId              string `json:"orderId"`              // Номер заказа в платежной системе
	UserName             string `json:"userName"`             // Логин служебной учётной записи продавца
	Password             string `json:"password"`             // Пароль служебной учётной записи продавца
	pan                  string `json:"pan"`                  // Номер платёжной карты
	cvc                  string `json:"cvc"`                  // Параметр обязателен, если для продавца не выбрано разрешение-может проводить оплату без подтверждения CVC
	expiry               string `json:"expiry"`               // Год и месяц окончания срока действия карты
	cardHolderName       string `json:"cardHolderName"`       // Имя держателя карты
	description          string `json:"description"`          // Описание заказа в свободной форме. Для включения в финансовую отчётность продавца передаются только первые 24 символа этого поля
	additionalParameters string `json:"additionalParameters"` // Дополнительные параметры заказа, которые сохраняются для просмотра из личного кабинета продавца.
	preAuth              bool   `json:"preAuth"`              // Параметр, определяющий необходимость предварительной авторизации (блокирования средств на счёте клиента до их списания)
	ip                   string `json:"ip"`                   // IP-адрес покупателя
	success              bool   `json:"success"`              // Указывает на успешность запроса: true, false
	code                 int    `json:"code"`                 // Код ошибки
	description_err      string `json:"description_err"`      // Подробное техническое объяснение ошибки для отображения пользователю
	message              string `json:"message"`              // Понятное описание ошибки для отображения пользователю
	mu                   sync.Mutex
}

// var (
// 	crtFile = filepath.Join(".", "certs", "client.crt")
// 	keyFile = filepath.Join(".", "certs", "client.key")
// )

func (pp ParamsPay) Register(rch chan string, crtFile, keyFile string) {
	log.SetPrefix("Rest client event: ")
	log.SetFlags(log.Lshortfile)

	// URL для проверки ответа боевого сервиса Сбера
	//apiUrl := "https://3dsec.sberbank.ru/payment/rest/register.do"

	// Тестовый сервер
	apiUrl := "https://localhost:8443/register"
	//Формирование параметров rest запроса. Params of request
	params := url.Values{"userName": {pp.UserName}}
	params.Set("password", pp.Password)
	params.Set("amount", pp.Amount)
	params.Set("returnUrl", pp.ReturnUrl)

	// Формирование строки параметров запроса. String data of request
	dreq := strings.NewReader(params.Encode())

	// Подгрузка сертификатов. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	fmt.Println("Good cert")
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
	log.Println("\nRest sberpay response:\n", string([]byte(body)))

	// Проверка идентификатора заказа. Checks order ID
	res := (string)([]byte(body))
	if strings.Contains(res, "orderId") || strings.Contains(res, "formUrl") {
		fmt.Println("Запрос регистрации методом register.do - ОК!")
	}

	data, err := json.MarshalIndent(res, "=", " ")
	if err != nil {
		log.Fatalf("Сбой маршалинга JSON: %s", err)
	}
	fmt.Printf("Data response = %s\n", data)

	rs := (string)([]byte(data))
	oid := regexp.MustCompile(`orderId:......................................`)
	boid := oid.FindAllString(rs, -1)
	soid := fmt.Sprint(boid)
	orderId := strings.Trim(soid, `[orderId: "]`)
	fmt.Println(" orderId =", orderId)
	rch <- orderId
}
