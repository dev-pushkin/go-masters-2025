package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_fibonacci(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test 1 -  Zero value",
			args: args{n: 0},
			want: 0,
		},
		{
			name: "Test 2 - Negative value",
			args: args{n: -1},
			want: 0,
		},
		{
			name: "Test 4 - Positive, valid value",
			args: args{n: 10},
			want: 55,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fibonacci(tt.args.n); got != tt.want {
				t.Errorf("fibonacci() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler(t *testing.T) {
	type args struct {
		value string
	}
	type want struct {
		statusCode  int
		body        string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test 1 - empty value",
			args: args{value: ""},
			want: want{
				statusCode:  http.StatusBadRequest,
				body:        "invalid parameter\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test 2 - big value",
			args: args{value: "100"},
			want: want{
				statusCode:  http.StatusBadRequest,
				body:        "parameter too large\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test 3 - valid value",
			args: args{value: "10"},
			want: want{
				statusCode:  http.StatusOK,
				body:        `{"result": 55}`,
				contentType: "application/json",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/fibonacci?N="+tt.args.value, nil)
			rr := httptest.NewRecorder()
			handler(rr, req)
			res := rr.Result()

			if res.StatusCode != int(tt.want.statusCode) {
				t.Errorf("Expected status code %d, got %d", tt.want.statusCode, res.StatusCode)
			}
			if string(rr.Body.String()) != tt.want.body {
				t.Errorf("Expected body %s, got %s", tt.want.body, string(rr.Body.String()))
			}
			if res.Header.Get("Content-Type") != tt.want.contentType {
				t.Errorf("Expected Content-Type %s, got %s", tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}
