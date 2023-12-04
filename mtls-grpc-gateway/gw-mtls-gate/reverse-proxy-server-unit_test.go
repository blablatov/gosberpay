// Модульное (unit) тестирование прокси-сервера grpc-gateway.
// Определены тестовые сигнатуры реальных методов прокси-сервера.
// Объявлены тестовые типы соответствующие интерфейсу, как его экземпляры.
// go test -v reverse-proxy-server-unit_test.go

package main

import (
	"fmt"
	"reflect"
	"testing"
)

type prxServer interface {
	NewOauthAccess(string) bool
	LoadX509KeyPair(interface{}, string) (bool, error)
	NewCertPool() interface{}
	ReadFile(string) (bool, error)
	AppendCertsFromPEM(CertPool) (ok bool)

	orderUnaryClientInterceptor(interface{}, string, ClientConn, ...interface{}) error
	NewServeMux(DialOptions) interface{}
}

type prxData struct {
	prxServer
	success bool
	CertPool
	ClientConn
}

type DialOption interface {
	WithPerRPCCredentials(interface{}) bool
	WithUnaryInterceptor(interface{}) bool
	WithTransportCredentials(interface{}) bool
	RegisterRestRequestsHandlerFromEndpoint(interface{}, string, *DialOptions) error
}

type DialOptions struct {
	DialOption
	success bool
}

type CertPool struct {
	okCert bool
}

type ClientConn struct {
	clconn bool
}

func (p *prxData) NewOauthAccess(string) bool {
	if p.success {
		return true
	}
	return false
}

func TestNewOauthAccess(t *testing.T) {
	type args struct {
		tok prxServer
	}
	tests := []struct {
		token      string
		args       args
		wantExists bool
	}{
		{
			token: "token exists",
			args: args{
				tok: &prxData{success: true},
			},
			wantExists: true,
		}, {
			token: "token not exists",
			args: args{
				tok: &prxData{success: false},
			},
			wantExists: false,
		},
	}

	var f prxData

	for _, tt := range tests {
		t.Run(tt.token, func(t *testing.T) {
			gotExists := f.NewOauthAccess("auth-token")
			if gotExists != false {
				t.Errorf("Check func NewOauthAccess() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func (p *prxData) LoadX509KeyPair(interface{}, string) (bool, error) {
	if p.success {
		return true, nil
	}
	return false, fmt.Errorf("LoadX509KeyPair test error")
}

func TestLoadX509KeyPair(t *testing.T) {
	type args struct {
		key prxServer
	}
	tests := []struct {
		crtFile    string
		keyFile    string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			crtFile: "crtFile exists",
			args: args{
				key: &prxData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			crtFile: "crtFile not exists",
			args: args{
				key: &prxData{success: false},
			},
			wantErr:    fmt.Errorf("crtFile test error"),
			wantExists: false,
		},

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

func (p *prxData) NewCertPool() interface{} {
	if p.success {
		return true
	}
	return false
}

func TestNewCertPool(t *testing.T) {
	type args struct {
		cer prxServer
	}
	tests := []struct {
		cert       string
		args       args
		wantExists bool
	}{
		{
			cert: "cert exists",
			args: args{
				cer: &prxData{success: true},
			},
			wantExists: true,
		}, {
			cert: "cert not exists",
			args: args{
				cer: &prxData{success: false},
			},
			wantExists: false,
		},
	}

	var f prxData

	for _, tt := range tests {
		t.Run(tt.cert, func(t *testing.T) {
			gotExists := f.NewCertPool()
			if gotExists != false {
				t.Errorf("Check func NewCertPool() gotExists = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}

func (p *prxData) ReadFile(string) (bool, error) {
	if p.success {
		return true, nil
	}
	return false, fmt.Errorf("ReadFile test error")
}

func TestReadFile(t *testing.T) {
	type args struct {
		read prxServer
	}
	tests := []struct {
		file       string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			file: "file exists",
			args: args{
				read: &prxData{success: true},
			},
			wantErr:    nil,
			wantExists: true,
		}, {
			file: "file not exists",
			args: args{
				read: &prxData{success: false},
			},
			wantErr:    fmt.Errorf("crtFile test error"),
			wantExists: false,
		},
	}

	var f prxData

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			gotExists, gotErr := f.ReadFile("local_file")
			if gotExists != false {
				t.Errorf("Check func ReadFile() gotExists = %v, want %v", gotExists, tt.wantExists)
			}

			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func ReadFile() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (p *prxData) AppendCertsFromPEM(CertPool) (ok bool) {
	if p.okCert {
		return true
	}
	return false
}

func TestAppendCertsFromPEM(t *testing.T) {
	type args struct {
		app prxServer
	}
	tests := []struct {
		apc  string
		args args
		ok   bool
	}{
		{
			apc: "AppendCerts exists",
			args: args{
				app: &prxData{success: true},
			},
			ok: true,
		}, {
			apc: "AppendCerts not exists",
			args: args{
				app: &prxData{success: false},
			},
			ok: false,
		},
	}

	var f prxData
	for _, tt := range tests {
		t.Run(tt.apc, func(t *testing.T) {
			ok := f.AppendCertsFromPEM(f.CertPool)
			if ok != false {
				t.Errorf("Check func AppendCertsFromPEM() gotExists = %v, want %v", ok, tt.ok)
			}
		})
	}
}

func (p *prxData) orderUnaryClientInterceptor(interface{}, string, ClientConn, ...interface{}) error {
	if p.clconn {
		return nil
	}
	return fmt.Errorf("orderUnaryClientInterceptor test error")
}

func TestOrderUnaryClientInterceptor(t *testing.T) {
	type args struct {
		inc prxServer
	}
	tests := []struct {
		method     string
		args       args
		wantErr    error
		wantExists bool
	}{
		{
			method: "method exists",
			args: args{
				inc: &prxData{success: true},
			},
			wantErr: nil,
		}, {
			method: "method not exists",
			args: args{
				inc: &prxData{success: false},
			},
			wantErr: fmt.Errorf("crtFile test error"),
		},
	}

	var f prxData

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			gotErr := f.orderUnaryClientInterceptor(nil, "POST", f.ClientConn, nil, nil)
			if reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("Check func orderUnaryClientInterceptor() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func (d *DialOptions) NewServeMux(DialOptions) interface{} {
	opts := DialOptions{
		d.DialOption,
		d.success,
	}
	return opts
}

func (d *DialOptions) WithPerRPCCredentials(interface{}) bool {
	if d.success {
		return true
	}
	return false
}

func (d *DialOptions) WithUnaryInterceptor(interface{}) bool {
	if d.success {
		return true
	}
	return false
}

func (d *DialOptions) WithTransportCredentials(interface{}) bool {
	if d.success {
		return true
	}
	return false
}

func (d *DialOptions) RegisterRestRequestsHandlerFromEndpoint(interface{}, string, *DialOptions) error {
	if d.success {
		return nil
	}
	return fmt.Errorf("RegisterRestRequestsHandlerFromEndpoint test error")
}

func TestRegisterRestRequestsHandlerFromEndpoint(t *testing.T) {
	type args struct {
		point DialOption
	}
	tests := []struct {
		endpoint string
		args     args
		wantErr  error
	}{
		{
			endpoint: "endpoint exists",
			args: args{
				point: &DialOptions{success: true},
			},
			wantErr: nil,
		}, {
			endpoint: "endpoint not exists",
			args: args{
				point: &DialOptions{success: false},
			},
			wantErr: fmt.Errorf("endpoint test error"),
		},
	}

	var f DialOptions

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			gotErr := f.RegisterRestRequestsHandlerFromEndpoint(nil, "endpoint", nil)
			if gotErr == nil {
				t.Errorf("Check func RegisterRestRequestsHandlerFromEndpoint() gotExists = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
