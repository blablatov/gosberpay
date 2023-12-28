// Testing remote functions without using network
// Модульное тестирование бизнес-логики удаленных методов без передачи по сети.
// С запуском стандартного gRPC-сервера поверх HTTP/2 на реальном порту.
// Имитация запуска сервера с использованием буфера.

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

// Implements server on real port
// Реализует имитацию запуска сервера на реальном порту
func initGRPCServerHTTP2() {
	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRestRequestsServer(s, &server{})
	// Регистрация службы сервиса gRPC. Register reflection service on gRPC server
	reflection.Register(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}

// Conventional test, starts gRPC server and client testp the service via RPC
// Обычный тест, запускается сервер gRPC, клиент тестирует через RPC метод AddRegister
func TestAddRegister(t *testing.T) {
	// Starting a conventional gRPC server runs on HTTP2
	// Запускается стандартный gRPC-сервер поверх HTTP/2
	initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Подключение к серверному приложению
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
	if err != nil { // Checks response. Проверка ответа
		log.Fatalf("Could not add register: %v", err)
	}
	log.Printf("Res of AddRegister: %s", r.Value)
}

// Тест метода GetOrderStatusExtended. Test of GetOrderStatusExtended
func TestGetOrderStatusExtended(t *testing.T) {
	// Starting a conventional gRPC server runs on HTTP2
	// Запуск стандартного gRPC-сервер поверх HTTP/2
	//initGRPCServerHTTP2()
	conn, err := grpc.Dial(address, grpc.WithInsecure()) // Подключение к серверному приложению
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRestRequestsClient(conn)

	// Параметры запроса. Params of request
	orderId := "70906e55-7114-41d6-8332-4609dc6590f4"

	clientDeadline := time.Now().Add(time.Duration(2 * time.Second))
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	// Calls remote method of GetOrderStatusExtended
	// Вызов удаленного метода GetOrderStatusExtended
	r, err := c.GetOrderStatusExtended(ctx, &pb.Status{OrderId: orderId})
	if err != nil { // Checks response. Проверка ответа
		log.Fatalf("Could not GetOrderStatusExtended: %v", err)
	}
	log.Printf("Res of GetOrderStatusExtended: %s", r.Value)
}

// ////////////////////////////////////////////
// Initialization of BufConn
// Package bufconn provides a net. Conn implemented by a buffer
// Реализует имитацию запуска сервера на реальном порту с использованием буфера
func initGRPCServerBufConn() {
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

func getBufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

// Тест метода AddRegister с использованием буфера. Test written using Bufconn
func TestAddRegisterBufConn(t *testing.T) {
	ctx := context.Background()
	initGRPCServerBufConn()
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

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second) //err context deadline
	defer cancel()
	r, err := c.AddRegister(ctx, &pb.Register{UserName: username, Password: password, Amount: amount, ReturnUrl: returnUrl})
	if err != nil {
		log.Fatalf("Could not add register: %v", err)
	}
	log.Printf(r.Value)
}

// Тест метода GetOrderStatusExtended с использованием буфера. Test written using Bufconn
func TestGetOrderStatusExtendedBufConn(t *testing.T) {
	ctx := context.Background()
	initGRPCServerBufConn()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewRestRequestsClient(conn)

	// Параметры запроса. Params of request
	orderId := "70906e55-7114-41d6-8332-4609dc6590f4"

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second) //err context deadline
	defer cancel()
	r, err := c.GetOrderStatusExtended(ctx, &pb.Status{OrderId: orderId})
	if err != nil {
		log.Fatalf("Could not GetOrderStatusExtended: %v", err)
	}
	log.Printf("Res of GetOrderStatusExtended: %s", r.Value)
}

// Benchmark tests
// Тестирование производительности в цикле за указанное колличество итераций.
// Метода AddRegister
func Benchmark_AddRegisterBufConn(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < 25; i++ {

		ctx := context.Background()
		initGRPCServerBufConn()
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

// Тестирование производительности в цикле за указанное колличество итераций.
// Метода GetOrderStatusExtended
func Benchmark_GetOrderStatusExtendedBufConn(b *testing.B) {

	b.ReportAllocs()
	for i := 0; i < 25; i++ {

		ctx := context.Background()
		initGRPCServerBufConn()
		conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(getBufDialer(listener)), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewRestRequestsClient(conn)

		// Параметры запроса. Params of request
		orderId := "70906e55-7114-41d6-8332-4609dc6590f4"

		//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second) //err context deadline
		defer cancel()
		r, err := c.GetOrderStatusExtended(ctx, &pb.Status{OrderId: orderId})
		if err != nil {
			log.Fatalf("Could not GetOrderStatusExtended: %v", err)
		}
		log.Printf("Res of GetOrderStatusExtended: %s", r.Value)
	}
}
