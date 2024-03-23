// Традиционный тест get-запросов, с запущенным обратным прокси-сервером.
// go run reverse-proxy-server.go
// и тестовым вебсервером для отладки REST запросов.
// go run gosberpay.go
// + Benchmark тестирование

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestStrings(t *testing.T) {
	var tests = []struct {
		payUrl string
	}{
		{"https://test.ru"},
		{"https://3dsec.sberbank.ru/payment/merchants/test/payment_ru.html"},
		{"https://google.com"},
	}

	var prevpayUrl string
	for _, test := range tests {
		if test.payUrl != prevpayUrl {
			fmt.Printf("\n%s\n", test.payUrl)
			prevpayUrl = test.payUrl
		}
	}
}

func TestGet(t *testing.T) {

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

	var w io.Writer
	w = os.Stdout

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

	count := make(map[string]int)
	pm, err := ioutil.ReadFile("register_log.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "err open: %v\n", err)
		return
	}

	rs := string(pm)
	oid := regexp.MustCompile(`orderId=....................................`)
	boid := oid.FindAllString(rs, -1)
	soid := fmt.Sprint(boid)
	orderId := strings.Trim(soid, `[orderId: "]`)
	fmt.Println(" orderId", orderId)

	for _, ln := range strings.Split(string(orderId), `""`) {
		count[ln]++
	}

	for param, n := range count {
		if n > 0 {
			getUrl := "https://3dsec.sberbank.ru/payment/merchants/test/payment_ru.html?mdOrder" + param
			resp, err := client.Get(getUrl)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if resp.StatusCode > 500 {
				log.Fatalf("Response status code: %d and\nbody: %s\n", resp.StatusCode, b)
			}
			if err != nil {
				log.Fatal((err))
			}

			fmt.Printf("\nPage to pay of SberPay service:\n %s\n", b)

		}
	}
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 350; i++ {
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

		var w io.Writer
		w = os.Stdout

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

		count := make(map[string]int)
		pm, err := ioutil.ReadFile("register_log.txt")
		if err != nil {
			fmt.Fprintf(os.Stderr, "err open: %v\n", err)
			return
		}
		for _, ln := range strings.Split(string(pm), `""`) {
			count[ln]++
		}

		for param, n := range count {
			if n > 0 {
				getUrl := "https://localhost:8444/v1/register/" + param
				resp, err := client.Get(getUrl)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				b, err := ioutil.ReadAll(resp.Body)
				if resp.StatusCode > 555 { //resp.StatusCode > 299
					log.Fatalf("Response status code: %d and\nbody: %s\n", resp.StatusCode, b)
				}
				if err != nil {
					log.Fatal((err))
				}
				fmt.Printf("\nResponse get gateway = %s\n", b)

			}
		}
	}
}
