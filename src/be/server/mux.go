package server

import (
	"net/http"
	"reflect"
	"runtime"
	"runtime/debug"

	"be/common"
	xe "be/common/error"
	"be/common/log"

	"github.com/gorilla/mux"
)

type WWWMux struct {
	r *mux.Router
}

func New() *WWWMux {
	return &WWWMux{r: mux.NewRouter()}
}

func (m *WWWMux) GetRouter() *mux.Router {
	return m.r
}

// 记录日志，同时包裹没有捕获到的异常
func AccessLogHandler(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Errorln(err)
				log.Errorf("HTTP Handler异常，stack信息 %s", debug.Stack())
				common.ResMsg(w, 500, xe.HandleRequestError().Error())
				return
			}
		}()
		log.Infof("%s - %s", r.Method, r.RequestURI)
		h(w, r)
	}
}

func (m *WWWMux) RegistURLMapping(path string, method string, handle func(http.ResponseWriter, *http.Request)) {
	log.WithFields(log.Fields{
		"path":   path,
		"method": method,
		"handle": runtime.FuncForPC(reflect.ValueOf(handle).Pointer()).Name(),
	}).Info("注册URL映射")
	handle = AccessLogHandler(handle)
	m.r.HandleFunc(path, handle).Methods(method)
}
