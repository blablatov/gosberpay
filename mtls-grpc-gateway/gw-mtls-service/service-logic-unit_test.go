// Модульное (unit) тестирование сервиса gPRC.
// Определены тестовые сигнатуры реальных методов сервиса.
// Тестовые типы соответствующие интерфейсу, как его экземпляры.
// go test -v service-logic-unit_test.go

package main

import (
	"fmt"
	"reflect"
	"testing"
)

type rester interface {
	AddRegister(interface{}, *restData, ...interface{}) (bool, error)
	NewRestRequestsClient(interface{}) bool
}

type restData struct {
	success   bool
	username  string `json:"UserName`
	password  string
	amount    string
	returnUrl string
	rester
}

func (rd *restData) AddRegister(interface{}, *restData, ...interface{}) (bool, error) {
	if rd.success {
		return true, nil
	}
	return false, fmt.Errorf("AddRegister test error")
}

func TestAddRegister(t *testing.T) {
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
		////////////////////////////////////////////////
		{
			keyFile: "keyFile exists",
			args: args{
				key: &prxData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			keyFile: "keyFile not exists",
			args: args{
				key: &prxData{success: false},
			},
			wantErr:    fmt.Errorf("keyFile test error"),
			wantExists: false,
		},
	}

	var f prxData

	for _, tt := range tests {
		t.Run(tt.crtFile, func(t *testing.T) {
			gotExists, gotErr := f.LoadX509KeyPair(tt.wantErr, "crtFile")
			if gotExists != false {
				t.Errorf("Check func LoadX509KeyPair() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if gotErr == tt.wantErr {
				t.Errorf("Check func LoadX509KeyPair() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
		t.Run(tt.keyFile, func(t *testing.T) {
			gotExists, gotErr := f.LoadX509KeyPair(tt.args.key, "keyFile")
			if gotExists != false {
				t.Errorf("Check func LoadX509KeyPair() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if gotErr == tt.wantErr {
				t.Errorf("Check func LoadX509KeyPair() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
