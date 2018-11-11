package main

import (
	go_log "log"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"be/common/log"
	"be/handle"
	"be/mysql"
	"be/options"
	"be/server"
	"be/session"
)

func doServe() {
	defer func() {
		if err := recover(); err != nil {
			doServe()
		}
	}()

	// todo 启动project和service的校对比较服务

	// 初始化DB
	mysql.DB.InitConn()
	// 初始化服务,并启动服务
	mux := server.New()
	// URL映射
	handle.InitHandle(mux)
	if options.Options.EnableTls {
		// HTTP服务重定向到HTTPS
		go http.ListenAndServe(options.Options.HTTPAddress+":"+strconv.FormatUint(options.Options.HTTPPort, 10), http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			http.Redirect(res, req, "https://"+req.Host+req.URL.String(), http.StatusMovedPermanently)
		}))

		srv := &http.Server{
			Handler:      mux.GetRouter(),
			Addr:         options.Options.HTTPSAddr,
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
			ErrorLog:     go_log.New(log.StandardLogger().Writer(), "", 0),
		}

		// 启动主服务
		log.Fatal(srv.ListenAndServeTLS(options.Options.CertFile, options.Options.KeyFile))
	} else {
		srv := &http.Server{
			Handler:      mux.GetRouter(),
			Addr:         options.Options.HTTPAddress + ":" + strconv.FormatUint(options.Options.HTTPPort, 10),
			WriteTimeout: 60 * time.Second,
			ReadTimeout:  60 * time.Second,
			ErrorLog:     go_log.New(log.StandardLogger().Writer(), "", 0),
		}

		// 启动主服务
		log.Fatal(srv.ListenAndServe())
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	// 从命令行、配置文件初始化配置
	options.Options.InitOptions()

	// 初始化Log
	log.InitLog()
	// 可以使用log了
	log.Infoln("日志文件初始化成功")
	// 初始化cookie
	session.InitCM()
	// 启动服务
	doServe()
}
