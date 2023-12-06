// Testing remote functions without using network
// Модульное тестирование бизнес-логики удаленных методов с использованием буфера, без передачи по сети.

package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/blablatov/gosberpay/mtls-grpc-gateway/gw-mtls-proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/test/bufconn"
)

const (
	address = "localhost:50051"
	bufSize = 1024 * 1024
)

var listener *bufconn.Listener

func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRestRequestsServer(s, &server{})
	// Регистрация службы сервиса gRPC. Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// Initialization of BufConn
// Package bufconn provides a net. Conn implemented by a buffer
// Реализует имитацию запуска сервера на реальном порту с использованием буфера
func initGRPCServerBuffConn() {
	listener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterRestRequestsServer(s, &server{})
	// Регистрация службы сервиса gRPC. Register reflection service on gRPC server.
	reflection.Register(s)
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// Conventional test that starts a gRPC server and client test the service with RPC
// Обычный тест, запускает сервер gRPC, клиент проверяет службу с помощью RPC
func TestServer_AddRegister(t *testing.T) {
	// Starting a conventional gRPC server runs on HTTP2
	// Запускаем стандартный gRPC-сервер поверх HTTP/2
	initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Подключаемся к серверному приложению
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRestRequestsClient(conn)

	// Параметры запроса. Params of request
	username := "gotov"
	password := "toor"
	amount := "99999"
	returnUrl := "https://test.ru/"

	clientDeadline := time.Now().Add(time.Duration(2 * time.Second))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	// Calls remote method of AddRegister
	// Вызов удаленного метода AddRegister
	r, err := c.AddRegister(ctx, &pb.Register{UserName: username, Password: password, Amount: amount, ReturnUrl: returnUrl})
	if err != nil { // Checks response. Проверяем ответ
		log.Fatalf("Could not add register: %v", err)
	}
	log.Printf("Res %s", r.Value)
}

// Тест с использованием буфера. Test written using Buffconn
func TestServer_AddProductBufConn(t *testing.T) {
	ctx := context.Background()
	initGRPCServerBuffConn()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRestRequestsClient(conn)

	// Параметры запроса. Params of request
	username := "gotov"
	password := "toor"
	amount := "99999"
	returnUrl := "https://test.ru/"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddRegister(ctx, &pb.Register{UserName: username, Password: password, Amount: amount, ReturnUrl: returnUrl})
	if err != nil {
		log.Fatalf("Could not add register: %v", err)
	}
	log.Printf(r.Value)
}

// Тестирование производительности в цикле за указанное колличество итераций
func BenchmarkServer_AddRegisterBufConn(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 25; i++ {
		ctx := context.Background()
		initGRPCServerBuffConn()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewRestRequestsClient(conn)

		// Параметры запроса. Params of request
		username := "gotov"
		password := "toor"
		amount := "99999"
		returnUrl := "https://test.ru/"

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.AddRegister(ctx, &pb.Register{UserName: username, Password: password, Amount: amount, ReturnUrl: returnUrl})
		if err != nil {
			log.Fatalf("Could not add register: %v", err)
		}
		log.Printf(r.Value)
	}
}
