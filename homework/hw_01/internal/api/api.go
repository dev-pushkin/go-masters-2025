package api

import (
	"log/slog"
	"net/http"

	"github.com/go_course_master/homework/hw_01/internal/app"
)

type Api struct {
	l      *slog.Logger
	server *http.Server
	app    *app.App
}

func NewApi(l *slog.Logger, s *http.Server, app *app.App) *Api {
	return &Api{
		l:      l,
		server: s,
		app:    app,
	}
}

func (a *Api) RegisterRouteAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/encrypt", a.encryptHandler)
	mux.HandleFunc("/decrypt", a.decryptHandler)

	a.server.Handler = mux
	a.server.Addr = addr

	a.l.Info("Starting server", "addr", addr)
	return a.server.ListenAndServe()
}

func (a *Api) encryptHandler(w http.ResponseWriter, r *http.Request) {
	data := r.FormValue("data")
	secret := r.FormValue("secret")

	resData, err := a.app.Encrypt(secret, data)
	if err != nil {
		responseError(a.l, w, err)
		return
	}
	responseSuccess(w, resData)
}

func (a *Api) decryptHandler(w http.ResponseWriter, r *http.Request) {
	secretKey := r.FormValue("secretKey")
	data := r.FormValue("data")

	resData, err := a.app.Decrypt(secretKey, data)
	if err != nil {
		responseError(a.l, w, err)
		return
	}

	responseSuccess(w, resData)
}
