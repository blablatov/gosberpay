// Методы (логика) gRPC-сервиса

package main

import (
	"context"
	"fmt"
	"log"
	_ "net/http/pprof"

	pb "github.com/blablatov/gosberpay/mtls-grpc-gateway/gw-mtls-proto"
	ss "github.com/blablatov/gosberpay/rest/getOrderStatusExtended"
	rg "github.com/blablatov/gosberpay/rest/register"
	wrapper "github.com/golang/protobuf/ptypes/wrappers"
	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Implements server.
// Сервер используется для реализации services
type server struct {
	restMap   map[string]*pb.Register
	statusMap map[string]*pb.Status
	paramRest *rg.ParamsPay
}

type emb struct {
	rg.ParamsPay
}

// Метод запроса регистрации заказа (register.do)
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

	// Слайс для хранения всех значений мапы. Slice for map
	sm := make([]string, 0, len(s.restMap))

	for k, _ := range s.restMap {
		if k != "" {
			sm = append(sm, k)
		}
	}

	// Тестовый вывод значений. For test
	for k, v := range sm {
		if v != "" {
			log.Printf("Param[%v] = %v\n", k, v)
		}
	}

	rch := make(chan string, 2)
	ch := make(chan int, 1)
	var rs string

	// Вызов метода регистрации заказа, register.do
	for i := 0; i < 10; i++ {
		// Выбор мультиплексирования. Select of multiplexing
		select {
		case ch <- i:
		case x := <-ch:
			go func() {
				rd := rg.ParamsPay{
					UserName:  in.UserName,
					Password:  in.Password,
					Amount:    in.Amount,
					ReturnUrl: in.ReturnUrl,
				}
				rg.ParamsPay.Register(rd, rch, crtFile, keyFile)

				rs = fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
				//return &wrapper.StringValue{Value: rs}, nil

				fmt.Println("goroutine1 =", x)
				ch <- 1
			}()
			<-ch

		case x := <-ch:
			go func() {
				rd := rg.ParamsPay{
					UserName:  in.UserName,
					Password:  in.Password,
					Amount:    in.Amount,
					ReturnUrl: in.ReturnUrl,
				}
				rg.ParamsPay.Register(rd, rch, crtFile, keyFile)

				rs = fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
				//return &wrapper.StringValue{Value: rs}, nil

				fmt.Println("goroutine2 =", x)
				ch <- 1
			}()
			<-ch

		case x := <-ch:
			go func() {
				rd := rg.ParamsPay{
					UserName:  in.UserName,
					Password:  in.Password,
					Amount:    in.Amount,
					ReturnUrl: in.ReturnUrl,
				}
				rg.ParamsPay.Register(rd, rch, crtFile, keyFile)

				rs = fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
				//return &wrapper.StringValue{Value: rs}, nil

				fmt.Println("goroutine3 =", x)
				ch <- 1
			}()
			<-ch

		case x := <-ch:
			go func() {
				rd := rg.ParamsPay{
					UserName:  in.UserName,
					Password:  in.Password,
					Amount:    in.Amount,
					ReturnUrl: in.ReturnUrl,
				}
				rg.ParamsPay.Register(rd, rch, crtFile, keyFile)

				rs = fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
				//return &wrapper.StringValue{Value: rs}, nil

				fmt.Println("goroutine4 =", x)
				ch <- 1
			}()
			<-ch

		case x := <-ch:
			go func() {
				rd := rg.ParamsPay{
					UserName:  in.UserName,
					Password:  in.Password,
					Amount:    in.Amount,
					ReturnUrl: in.ReturnUrl,
				}
				rg.ParamsPay.Register(rd, rch, crtFile, keyFile)

				rs = fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
				//return &wrapper.StringValue{Value: rs}, nil

				fmt.Println("goroutine5 =", x)
				ch <- 1
			}()
			<-ch

		default:
		}
	}

	//rs := fmt.Sprintf("orderId=%s formUrl=%s", <-rch, <-rch)
	return &wrapper.StringValue{Value: rs}, nil
}

// Method get of product. Метод сервера GetRegister получить параметр
func (s *server) GetRegister(ctx context.Context, in *wrapper.StringValue) (*pb.Register, error) {
	value, exists := s.restMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}
	return &pb.Register{OrderNumber: in.Value}, nil
}

// Метод запроса состояния заказа (getOrderStatusExtended.do)
func (s *server) GetOrderStatusExtended(ctx context.Context, in *pb.Status) (*wrapper.StringValue, error) {
	// Bad request, generate and sends of error to client.
	// Некорректный запрос. Сгенерировать и отправить клиенту ошибку.
	if in.OrderId == "" {
		log.Printf("OrderId is invalid! -> Received OrderId %s", in.OrderNumber)
		// Creates state with code of error. Создаем состояние с кодом ошибки InvalidArgument.
		errorStatus := status.New(codes.InvalidArgument, "Invalid information received")
		// Describes type of error. Описываем тип ошибки BadRequest_FieldViolation
		ds, err := errorStatus.WithDetails(
			&epb.BadRequest_FieldViolation{
				Field:       "OrderId",
				Description: fmt.Sprintf("OrderId received is not valid %s : %s", in.OrderId, in.UserName),
			},
		)
		if err != nil {
			return nil, errorStatus.Err()
		}
		return nil, ds.Err()
	}

	if s.statusMap == nil {
		s.statusMap = make(map[string]*pb.Status)
	}

	s.statusMap[in.OrderId] = in

	sm := make([]string, 0, len(s.statusMap))
	for k, _ := range s.statusMap {
		if k != "" {
			sm = append(sm, k)
		}
	}

	for k, v := range sm {
		if v != "" {
			log.Printf("Param[%v] = %v\n", k, v)
		}
	}

	sch := make(chan string, 10)

	// Вызов запроса состояния заказа (getOrderStatusExtended.do)
	go func() {
		rd := ss.StatusParam{
			OrderId: in.OrderId,
		}
		ss.StatusParam.OrderStatusExtended(rd, sch, crtFile, keyFile)
	}()

	rs := fmt.Sprintf("OrderStatus:%s", <-sch)
	return &wrapper.StringValue{Value: rs}, nil
}
