// Традиционный тест post-запросов, с запущенным обратным прокси-сервером:
// go run reverse-proxy-server.go
// и тестовым вебсервером для отладки REST запросов.
// go run gosberpay.go
// + Benchmark тестирование

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

func TestRegister(t *testing.T) {

	var wb []byte
	var buf bytes.Buffer

	// URL локального тестового сервера. Для облака указать внешний IP ВМ
	apiUrl := "https://localhost:8444/v1/register"

	// Формирование json параметров запроса. JSON params of request
	payload, _ := json.Marshal(struct {
		userName  string `json:"userName"`
		password  string `json:"password"`
		amount    string `json:"amount"`
		returnUrl string `json:"returnUrl"`
	}{
		userName:  "goman",
		password:  "qwerty",
		amount:    "9999",
		returnUrl: "https://test.ru",
	})

	// Подгрузка сертификата и ключа. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs")
	}

	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Форматирование запроса. Formatting of the request
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-Length", "8")
	req.Header.Set("Grpc-Metadata-Content-Type", "application/grpc")
	req.Header.Set("Content-Type", "application/json")

	// Формирование метаданных структуры запроса. Struct of request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter:       w,
				Certificates:       []tls.Certificate{cert},
				ServerName:         hostname, // NOTE: this is required!
				RootCAs:            certPool,
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
	log.Println("\nResponse of gateway to register: ", string([]byte(body)))

	buf.WriteString(string([]byte(body)))
	wb = buf.Bytes()

	// Запись тестовых значений в лог-файл.
	err = ioutil.WriteFile("register_log.txt", wb, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func TestOrderStatus(t *testing.T) {

	var wb []byte
	var buf bytes.Buffer

	// URL локального тестового сервера. Для облака указать внешний IP ВМ
	apiUrl := "https://localhost:8444/v1/status"

	// Формирование json параметров запроса. JSON params of request
	payload, _ := json.Marshal(struct {
		OrderId string `json:"orderId"`
	}{
		OrderId: "70906e55-7114-41d6-8332-4609dc6590f4",
	})

	// Подгрузка сертификата и ключа. Loads the certs
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("Сертификат и ключ не получены: %v\n", err)
	}

	// Create a certificate pool from the certificate authority
	// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("failed to append ca certs")
	}

	// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
	var w io.Writer
	w = os.Stdout

	// Форматирование запроса. Formatting of the request
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
	// Формирование заголовков запроса. Headers of request
	req.Header.Set("Content-Length", "8")
	req.Header.Set("Grpc-Metadata-Content-Type", "application/grpc")
	req.Header.Set("Content-Type", "application/json")

	// Формирование метаданных структуры запроса. Struct of request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				KeyLogWriter:       w,
				Certificates:       []tls.Certificate{cert},
				ServerName:         hostname, // NOTE: this is required!
				RootCAs:            certPool,
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
	log.Println("\nResponse of gateway to order-status: ", string([]byte(body)))

	buf.WriteString(string([]byte(body)))
	wb = buf.Bytes()

	// Запись тестовых значений в лог-файл.
	err = ioutil.WriteFile("status_log.txt", wb, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkAdd(b *testing.B) {

	var wb []byte
	var buf bytes.Buffer

	b.ReportAllocs()
	for i := 0; i < 350; i++ {

		// URL локального тестового сервера. Для облака указать внешний IP ВМ
		apiUrl := "https://localhost:8444/v1/register"

		// Формирование json параметров запроса. JSON params of request
		payload, _ := json.Marshal(struct {
			userName  string `json:"userName"`
			password  string `json:"password"`
			amount    string `json:"amount"`
			returnUrl string `json:"returnUrl"`
			/*{"password": "qwerty",
			"userName": "goman",
			"amount": "9999",
			"returnUrl": "https://test.ru"}//For tests of Soap UI*/
		}{
			userName:  "goman",
			password:  "qwerty",
			amount:    "9999",
			returnUrl: "https://test.ru",
		})

		// Подгрузка сертификата и ключа. Loads the certs
		cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
		if err != nil {
			log.Fatalf("Сертификат и ключ не получены: %v\n", err)
		}

		// Create a certificate pool from the certificate authority
		// Генерируем пул сертификатов в нашем локальном удостоверяющем центре
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		// Добавляем клиентские сертификаты из локального удостоверяющего центра в сгенерированный пул
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append ca certs")
		}

		// Logs CLIENT_SERVER_HANDSHAKE_TRAFFIC_SECRETS
		var w io.Writer
		w = os.Stdout

		// Форматирование запроса. Formatting of the request
		req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewBuffer(payload))
		// Формирование заголовков запроса. Headers of request
		req.Header.Set("Content-Type", "application/json")

		// Формирование метаданных структуры запроса. Struct of request
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					KeyLogWriter:       w,
					Certificates:       []tls.Certificate{cert},
					ServerName:         hostname, // NOTE: this is required!
					RootCAs:            certPool,
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
		log.Println("\nResponse of grpc gateway: ", string([]byte(body)))

		buf.WriteString(string([]byte(body)))
		wb = buf.Bytes()
	}
	// Запись тестовых значений в лог-файл.
	err := ioutil.WriteFile("bench_log.txt", wb, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
