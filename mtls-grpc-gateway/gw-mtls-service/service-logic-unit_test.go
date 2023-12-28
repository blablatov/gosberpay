// Модульное (unit) тестирование сервиса gPRC.
// Определены тестовые сигнатуры реальных методов сервиса.
// Объявлены тестовые типы соответствующие интерфейсу, как его экземпляры.
// go test -v service-logic-unit_test.go

package main

import (
	"fmt"
	"reflect"
	"testing"
)

type rester interface {
	AddRegister(interface{}, testData) (bool, error)
	GetRegister(interface{}, testData) (bool, error)
	GetOrderStatusExtended(interface{}, testData) (bool, error)
}

type restData struct {
	success   bool
	username  string
	password  string
	amount    string
	returnUrl string
	formUrl   string
	rester
	testData
}

type testData struct {
	success bool
}

func (rd *restData) AddRegister(interface{}, testData) (bool, error) {
	if rd.success {
		return true, nil
	}
	return false, fmt.Errorf("AddRegister test error")
}

func TestAddRegisterUnit(t *testing.T) {
	type args struct {
		reg rester
	}
	tests := []struct {
		username   string
		password   string
		amount     string
		returnUrl  string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			username: "username exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			username: "username not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("username test error"),
			wantExists: false,
		},

		{
			password: "password exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			password: "password not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("password test error"),
			wantExists: false,
		},

		{
			amount: "amount exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			amount: "amount not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("amount test error"),
			wantExists: false,
		},

		{
			returnUrl: "returnUrl exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			returnUrl: "returnUrl not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("returnUrl test error"),
			wantExists: false,
		},
	}

	var f restData

	for _, tt := range tests {
		t.Run(tt.username, func(t *testing.T) {
			gotExists, gotErr := f.AddRegister(nil, f.testData)
			if gotExists != false {
				t.Errorf("Check func AddRegister() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func AddRegister() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
		t.Run(tt.password, func(t *testing.T) {
			gotExists, gotErr := f.AddRegister(tt.args.reg, f.testData)
			if gotExists != false {
				t.Errorf("Check func AddRegister() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func AddRegister() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
		t.Run(tt.amount, func(t *testing.T) {
			gotExists, gotErr := f.AddRegister(tt.args.reg, f.testData)
			if gotExists != false {
				t.Errorf("Check func AddRegister() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func AddRegister() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
		t.Run(tt.returnUrl, func(t *testing.T) {
			gotExists, gotErr := f.AddRegister(tt.args.reg, f.testData)
			if gotExists != false {
				t.Errorf("Check func AddRegister() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func AddRegister() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (rd *restData) GetRegister(interface{}, testData) (bool, error) {
	if rd.success {
		return true, nil
	}
	return false, fmt.Errorf("GetRegister test error")
}

func TestGetRegister(t *testing.T) {
	type args struct {
		reg rester
	}
	tests := []struct {
		formUrl    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			formUrl: "formUrl exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			formUrl: "formUrl not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("formUrl test error"),
			wantExists: false,
		},
	}

	var f restData

	for _, tt := range tests {
		t.Run(tt.formUrl, func(t *testing.T) {
			gotExists, gotErr := f.GetRegister(nil, f.testData)
			if gotExists != false {
				t.Errorf("Check func GetRegister() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func GetRegister() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (rd *restData) GetOrderStatusExtended(interface{}, testData) (bool, error) {
	if rd.success {
		return true, nil
	}
	return false, fmt.Errorf("GetOrderStatusExtended test error")
}

func TestGetOrderStatusExtendedUnit(t *testing.T) {
	type args struct {
		reg rester
	}
	tests := []struct {
		orderId    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			orderId: "orderId exists",
			args: args{
				reg: &restData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			orderId: "orderId not exists",
			args: args{
				reg: &restData{success: false},
			},
			wantErr:    fmt.Errorf("orderId test error"),
			wantExists: false,
		},
	}

	var f restData

	for _, tt := range tests {
		t.Run(tt.orderId, func(t *testing.T) {
			gotExists, gotErr := f.GetOrderStatusExtended(nil, f.testData)
			if gotExists != false {
				t.Errorf("Check func GetOrderStatusExtended() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func GetOrderStatusExtended() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
