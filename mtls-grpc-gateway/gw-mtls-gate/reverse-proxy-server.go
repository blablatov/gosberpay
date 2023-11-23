package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	gw "gw-mtls-proto"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	//gw "github.com/blablatov/mtls-grpc-gateway/gw-mtls-gate"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
	crtFile = filepath.Join("..", "gw-mcerts", "client.crt")
	keyFile = filepath.Join("..", "gw-mcerts", "client.key")
	caFile  = filepath.Join("..", "gw-mcerts", "ca.crt")
)

const (
	grpcServerEndpoint = "localhost:50051"
	//address  = "net-tls-service:50051"
	hostname = "localhost"
)

func main() {
	log.SetPrefix("Gate event: ")
	log.SetFlags(log.Lshortfile)

	// Set up the credentials for the connection.
	// Значение токена OAuth2. Используем строку, прописанную в коде.
	tokau := oauth.NewOauthAccess(fetchToken())

	// Load the client certificates from disk
	// Создаем пары ключей X.509 непосредственно из ключа и сертификата сервера
	certificate, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
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

	// Connect to server. Data of auth. Соединения с сервером
	// Указываем аутентификационные данные для транспортного протокола с помощью DialOption.
	opts := []grpc.DialOption{
		// Указываем один и тот же токен OAuth в параметрах всех вызовов в рамках одного соединения.
		// Если нужно указывать токен для каждого вызова отдельно, используем CallOption.
		grpc.WithPerRPCCredentials(tokau),
		// Регистрация унарного перехватчика на gRPC-клиенте
		// Будет направлять все запросы к функции orderUnaryClientInterceptor
		grpc.WithUnaryInterceptor(orderUnaryClientInterceptor),
		// Указываем транспортные аутентификационные данные в виде параметров соединения
		// Поле ServerName должно быть равно значению Common Name, указанному в сертификате
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
			ServerName:   hostname, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})),
	}

	//ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second) //err context deadline
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server endpoint, gRPC server should be running and accessible
	// Сервер gRPC должен быть запущен и доступен
	mux := runtime.NewServeMux()
	err = gw.RegisterRestRequestsHandlerFromEndpoint(ctx, mux, grpcServerEndpoint, opts)
	if err != nil {
		log.Fatalf("Fail to register gRPC service endpoint: %v", err)
		return
	}

	// TLS connect. Подключение по протоколу TLS
	if err := http.ListenAndServeTLS(":8444", crtFile, keyFile, mux); err != nil {
		log.Fatalf("Could not setup HTTPS endpoint: %v", err)
	}
	/*if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not setup HTTP endpoint: %v", err)
	}*/
}

// The value of OAuth2 token. String of token is in the code
// Значение токена OAuth2. Используется строка прописанная в коде
func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "blablatok-tokblabla-blablatok",
	}
}

func orderUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// Pre-processor phase. Этап предобработки, есть доступ к RPC-запросу перед его отправкой на сервер
	log.Println("Method : " + method)
	// Invoking the remote method. Вызов удаленного RPC-метода с помощью UnaryInvoker.
	err := invoker(ctx, method, req, reply, cc, opts...)
	// Post-processor phase. Этап постобработки, можно обработать ответ или возникшую ошибку.
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("\n req = %v\n reply = %v\n", req, reply)
	return err
}
