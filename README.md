# 包介绍
## collector
实现prometheus.Collector接口,定义采集指标,采集数据.

## logger
log相关.

## ngx
定义采集指标的方法.


# 执行参数介绍
- logPath: 日志路径,默认当前目录下(./exporter.log)
- ngxStatusPath: nginx status 路径,默认(http://127.0.0.1/status)
- httpClientTimeout: 请求nginx status 超时时间,默认3秒
- namespace: exporter namespace,默认nginx
