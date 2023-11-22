package main

import (
	"context"
	"fmt"
	rg "register"
	"sync"

	//rg "getredis"
	"log"
	//rs "setredis"
	//"sync"
	pb "gw-mtls-proto"
	//rs "github.com/blablatov/grpc-dsn-dbms/grpc-redis"
	//pb "github.com/blablatov/mtls-grpc-gateway/gw-mtls-proto"
	"github.com/gofrs/uuid"
	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements server.
// Сервер используется для реализации services
type server struct {
	productMap map[string]*pb.Product
	restMap    map[string]*pb.Register
	paramRest  *rg.ParamsPay
}

type emb struct {
	rg.ParamsPay
}

// Method add of product. Метод сервера AddProduct, добавить товар
func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*wrapper.StringValue, error) {
	// Bad request, generate and sends of error to client.
	// Некорректный запрос. Сгенерировать и отправить клиенту ошибку.
	if in.Name == "-1" {
		log.Printf("Order ID is invalid! -> Received Order Name %s", in.Id)
		// Creates state with code of error. Создаем состояние с кодом ошибки InvalidArgument.
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		// Describes type of error. Описываем тип ошибки BadRequest_FieldViolation
		ds, err := errorStatus.WithDetails(
			&epb.BadRequest_FieldViolation{
				Field:       "Name",
				Description: fmt.Sprintf("Order Name received is not valid %s : %s", in.Id, in.Description),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}
	out, err := uuid.NewV4()
	if err != nil {
		return nil, status.Errorf(codes.Internal, " %v\nError while generating Product ID", err)
	}
	in.Id = out.String()
	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[in.Id] = in

	// chg := make(chan string, 1)
	// chs := make(chan string, 1)
	// chb := make(chan bool, 1)

	// // Input data to redis. Внесение данных в redis
	// var wg sync.WaitGroup
	// wg.Add(1) // Counter of goroutines. Значение счетчика.
	// go rs.SetRedis(in.Id, in.Name, wg, chs)
	// go func() {
	// 	wg.Wait()
	// 	close(chb)
	// }()

	// // Get data from redis. Получение данных
	// wg.Add(1) // Counter of goroutines. Значение счетчика.
	// gval := make(chan string)
	// go func() {
	// 	gval <- rs.GetRedis(in.Id, wg, chg)
	// }()
	// log.Println(gval)
	// rval := <-gval
	// go func() {
	// 	wg.Wait()
	// 	close(chb)
	// }()

	return &wrapper.StringValue{Value: in.Id}, status.New(codes.OK, "").Err()
}

// Method get of product. Метод сервера GetProduct получить товар
func (s *server) GetProduct(ctx context.Context, in *wrapper.StringValue) (*pb.Product, error) {
	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return nil, status.Errorf(codes.NotFound, "%v Product does not exist.", in.Value)
}

// Метод регистрации заказа sberpay
func (s *server) AddRegister(ctx context.Context, in *pb.Register) (*wrapper.StringValue, error) {
	// Bad request, generate and sends of error to client.
	// Некорректный запрос. Сгенерировать и отправить клиенту ошибку.
	if in.Amount == "-1" {
		log.Printf("Amount is invalid! -> Received OrderNumber %s", in.OrderNumber)
		// Creates state with code of error. Создаем состояние с кодом ошибки InvalidArgument.
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		// Describes type of error. Описываем тип ошибки BadRequest_FieldViolation
		ds, err := errorStatus.WithDetails(
			&epb.BadRequest_FieldViolation{
				Field:       "Amount",
				Description: fmt.Sprintf("OrderNumber received is not valid %s : %s", in.Amount, in.Description),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}

	if s.restMap == nil {
		s.restMap = make(map[string]*pb.Register)
	}

	s.restMap[in.UserName] = in
	s.restMap[in.Password] = in
	s.restMap[in.Amount] = in
	s.restMap[in.ReturnUrl] = in

	sm := make([]string, 0, len(s.restMap))
	for k, _ := range s.restMap {
		if k != "" {
			sm = append(sm, k)
		}
	}

	for k, v := range sm {
		if v != "" {
			log.Printf("Param[%v] = %v\n", k, v)
		}
	}

	rd := rg.ParamsPay{
		UserName:  in.UserName,
		Password:  in.Password,
		Amount:    in.Amount,
		ReturnUrl: in.ReturnUrl,
	}

	var mu sync.Mutex
	rch := make(chan string, 2)

	go func() {
		mu.Lock()
		rg.ParamsPay.Register(rd, rch, crtFile, keyFile)
		mu.Unlock()
	}()

	//log.Println(<-rch)
	rs := fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
	return &wrapper.StringValue{Value: rs}, nil
	//return &wrapper.StringValue{Value: in.OrderNumber}, status.New(codes.OK, "").Err()
}

// Method get of product. Метод сервера GetRegister получить параметр
func (s *server) GetRegister(ctx context.Context, in *wrapper.StringValue) (*pb.Register, error) {
	value, exists := s.restMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return &pb.Register{OrderNumber: in.Value}, nil
	//return nil, status.Errorf(codes.NotFound, "%v Param does not exist.", in.Value)
}
