package options

import (
	"flag"

	"github.com/vharitonsky/iniflags"
)

type WWWOptions struct {
	// 默认日志级别
	LogLevel string
	LogFile  string

	// HTTP服务监听地址
	HTTPAddress string
	HTTPPort    uint64

	// 是否使用HTTPS
	EnableTls bool

	// HTTPS服务监听地址
	HTTPSAddr string
	// TLS
	CertFile string
	KeyFile  string

	// 静态文件路径
	StaticFilePath string
	// 页面模板文件路径
	TemplateFilePath string

	// 密码加密盐
	PasswordSalt string
	// Cookie01
	Cookie01 string
	// Cookie02
	Cookie02 string

	// MySQL dataSourceName
	DataSourceName string
	// MySQL 最大连接数
	DBMaxOpenConn int
	// MySQL 最大闲置连接数
	DBMaxIdleConn int

	// 小程序信息
	TPAppId  string
	TPSecret string

	// Service信息
	ServiceHTTPAddress string
	ServiceHTTPPort    uint64

	// 源码路径
	CodesRootPath string

	// 下载源码超时时间
	CodeFetchTimeout int64
}

var Options WWWOptions

func (o *WWWOptions) InitOptions() {
	flag.StringVar(&o.LogLevel, "log_level", "DEBUG", "Log Level")
	flag.StringVar(&o.LogFile, "log_file", "D:\\logs\\sweetfish.log", "Log File")
	flag.StringVar(&o.HTTPAddress, "http_address", "0.0.0.0", "HTTP Address")
	flag.Uint64Var(&o.HTTPPort, "http_port", 8888, "HTTP Port")
	flag.StringVar(&o.DataSourceName, "dsn", "root:rootroot@tcp(127.0.0.1:3306)/sweetfish?autocommit=0&collation=utf8_general_ci", "MySQL DataSourceName")
	flag.IntVar(&o.DBMaxOpenConn, "max_open_conn", 32, "MySQL Max Open Connections")
	flag.IntVar(&o.DBMaxIdleConn, "max_idle_conn", 16, "MySQL Max Idle Connections")
	flag.StringVar(&o.StaticFilePath, "static_file_path", "", "StaticFilePath")
	flag.StringVar(&o.TemplateFilePath, "template_file_path", "", "TemplateFilePath")
	flag.StringVar(&o.PasswordSalt, "password_salt", "", "PasswordSalt")
	flag.StringVar(&o.Cookie01, "cookie_01", "", "Cookie01")
	flag.StringVar(&o.Cookie02, "cookie_02", "", "Cookie02")
	flag.StringVar(&o.HTTPSAddr, "https_addr", "0.0.0.0:443", "HTTPSAddr")
	flag.StringVar(&o.CertFile, "cert_file", "./key/server.pem", "CertFile")
	flag.StringVar(&o.KeyFile, "key_file", "./key/server.key", "KeyFile")
	flag.BoolVar(&o.EnableTls, "enable_tls", false, "EnableTls")
	flag.StringVar(&o.TPAppId, "tp_appid", "", "TPAppId")
	flag.StringVar(&o.TPSecret, "tp_secret", "", "TPSecret")
	flag.StringVar(&o.ServiceHTTPAddress, "service_http_address", "127.0.0.1", "Service HTTP Address")
	flag.Uint64Var(&o.ServiceHTTPPort, "service_http_port", 19512, "Service HTTP Port")
	flag.StringVar(&o.CodesRootPath, "codes_root_path", "D:\\logs", "CodesRootPath")
	flag.Int64Var(&o.CodeFetchTimeout, "code_fetch_timeout", 600, "CodeFetchTimeout")

	iniflags.Parse()
}
