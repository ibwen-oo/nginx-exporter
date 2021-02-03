package collector

import (
	"exporter-demo/ngx"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"sync"
)

const (
	DefaultNameSpace = "nginx"
	NGINXUP     = 1
	NGINXDOWN   = 0
)

// NginxCollector 实现了 Collector 接口
type NginxCollector struct {
	Client   *ngx.NgxClient
	mutex       sync.Mutex
	connectionsActive *prometheus.Desc
	connectionsAccepted *prometheus.Desc
	connectionsHandled *prometheus.Desc
	connectionsReading *prometheus.Desc
	connectionsWriting *prometheus.Desc
	connectionsWaiting *prometheus.Desc
	requestsTotal *prometheus.Desc
	up prometheus.Gauge
}

// NewNginxCollector 初始化NginxCollector结构体
func NewNginxCollector(namespace string, client *ngx.NgxClient) *NginxCollector {
	return &NginxCollector{
		Client:              client,
		connectionsActive:   prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_active"),
			"Active client connections",
			[]string{"role"},
			nil),

		connectionsAccepted: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_accepted"),
			"Accepted client connections",
			[]string{"role"},
			nil),
		connectionsHandled:  prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_handled"),
			"Handled client connections",
			[]string{"role"},
			nil),
		connectionsReading:  prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_reading"),
			"Reading client connections",
			[]string{"role"},
			nil),
		connectionsWriting:  prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_writing"),
			"Writing client connections",
			[]string{"role"},
			nil),
		connectionsWaiting:  prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "connections_waiting"),
			"Waiting client connections",
			[]string{"role"},
			nil),
		requestsTotal:       prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "requests_total"),
			"Requests client connections",
			[]string{"role"},
			nil),
		up:                  prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name:"up", Help: "Nginx status(up or down)",
				Namespace: namespace, ConstLabels: map[string]string{"role": "web"}},
				),
	}
}

// Describe Collector 接口的 Describe 方法
func (n *NginxCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.up.Desc()
	ch <- n.connectionsActive
	ch <- n.connectionsAccepted
	ch <- n.connectionsHandled
	ch <- n.connectionsReading
	ch <- n.connectionsWriting
	ch <- n.connectionsWaiting
	ch <- n.requestsTotal
}

// Collect Collector 接口的 Collect 方法
func (n *NginxCollector) Collect(ch chan<- prometheus.Metric) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	ms, err := n.Client.QueryNgxStatus()
	if err != nil {
		n.up.Set(NGINXDOWN)
		ch <- n.up
		zap.L().Error("Get metrics failde, nginx may be down.", zap.Error(err))
	}
	n.up.Set(NGINXUP)
	ch <- n.up
	ch <- prometheus.MustNewConstMetric(n.connectionsActive,
		prometheus.GaugeValue, float64(ms.Active), "web")
	ch <- prometheus.MustNewConstMetric(n.connectionsAccepted,
		prometheus.CounterValue, float64(ms.Accepted), "web")
	ch <- prometheus.MustNewConstMetric(n.connectionsHandled,
		prometheus.CounterValue, float64(ms.Handled), "web")
	ch <- prometheus.MustNewConstMetric(n.connectionsReading,
		prometheus.GaugeValue, float64(ms.Reading), "web")
	ch <- prometheus.MustNewConstMetric(n.connectionsWriting,
		prometheus.GaugeValue, float64(ms.Writing), "web")
	ch <- prometheus.MustNewConstMetric(n.connectionsWaiting,
		prometheus.GaugeValue, float64(ms.Waiting), "web")
	ch <- prometheus.MustNewConstMetric(n.requestsTotal,
		prometheus.CounterValue, float64(ms.Requests), "web")
}
