package main

import (
	"exporter-demo/collector"
	"exporter-demo/logger"
	"exporter-demo/ngx"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36"
)

func main() {
	logPath := flag.String("logPath", "./exporter.log", "exporter error log")
	ngxStatusPath := flag.String("ngxStatusPath", "http://123.57.93.249/status", "Nginx status url path")
	httpClientTimeout := flag.Duration("httpClientTimeout", time.Duration(time.Second*3), "The time the request nginx timed out")
	namespace := flag.String("namespace", "null", "exporter namespace")

	flag.Parse()

	// 0.初始化Logger
	if err := logger.Init(*logPath); err != nil {
		log.Fatalln("Init logger failed, error:", err)
	}

	// 1.构建http client,用于请求nginx status
	params := ngx.NgxClientParams{
		EndPoint:  ngxStatusPath,
		UserAgent: &UserAgent,
		Timeout:   *httpClientTimeout,
	}
	client, err := ngx.InitHttpClient(params)
	if err != nil {
		zap.L().Error("InitHttpClient failed", zap.Error(err))
		os.Exit(1)
	}

	// 2.prometheus client
	register := prometheus.NewRegistry()

	if *namespace == "null" {
		*namespace = collector.DefaultNameSpace
	}

	ngxCollector := collector.NewNginxCollector(*namespace, client)

	register.MustRegister(ngxCollector)

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		register,
	}

	handler := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{
		ErrorLog:      &log.Logger{},
		ErrorHandling: promhttp.ContinueOnError,
	})

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	err = http.ListenAndServe(":8888", nil)
	zap.L().Error("Start export http server failed", zap.Error(err))
}
