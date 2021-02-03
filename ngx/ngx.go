package ngx

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NgxClientParams 请求nginx status相关参数
type NgxClientParams struct {
	EndPoint  *string
	UserAgent *string
	Timeout   time.Duration
}


// NgxClient nginx 客户端信息
type NgxClient struct {
	endPoint   string
	httpClient *http.Client
}

// NgxMetrics 采集指标
type NgxMetrics struct {
	Active   int64 // 活跃的连接数
	Accepted int64 // 总共处理了多少个连接
	Handled  int64 // 成功创建了多少次握手
	Reading  int64 // 读取客户端连接数
	Writing  int64 // 响应数据到客户端的数量
	Waiting  int64 // 开启 keep-alive 的情况下,这个值等于 active – (reading+writing),意思就是 Nginx 已经处理完正在等候下一次请求指令的驻留连接.
	Requests int64 // 总共处理了多少个请求
}

// QueryNgxStatus 请求nginx status,获取监控指标
func (n *NgxClient) QueryNgxStatus() (metrics *NgxMetrics, err error) {
	response, err := n.httpClient.Get(n.endPoint)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	metrics, err = parseStatusData(bodyBytes)
	if err != nil {
		return nil, err
	}
	return metrics, err
}

// parseStatusData 解析监控指标
func parseStatusData(data []byte) (metrics *NgxMetrics, err error) {
	dataStr := string(data)
	dataSlice := strings.Split(dataStr, "\n")
	dataErr := fmt.Errorf("invalid value: %v\n", dataStr)
	if len(dataSlice) != 5 {
		return nil, dataErr
	}
	ac := strings.TrimSpace(strings.Split(dataSlice[0], ":")[1])

	activeConn, err := strconv.ParseInt(ac, 10, 64)
	if err != nil {
		return nil, dataErr
	}

	// accepts/handled/requests
	ahr := strings.Split(strings.TrimSpace(dataSlice[2]), " ")
	if len(ahr) != 3 {
		return nil, dataErr
	}

	accepts, err := strconv.ParseInt(strings.TrimSpace(ahr[0]), 10, 64)
	if err != nil {
		return nil, dataErr
	}
	handled, err := strconv.ParseInt(strings.TrimSpace(ahr[1]), 10, 64)
	if err != nil {
		return nil, dataErr
	}
	requests, err := strconv.ParseInt(strings.TrimSpace(ahr[2]), 10, 64)
	if err != nil {
		return nil, dataErr
	}

	// Reading/Writing/Waiting
	rww := strings.Split(strings.TrimSpace(dataSlice[3]), " ")
	reading, err := strconv.ParseInt(strings.TrimSpace(rww[1]), 10, 64)
	if err != nil {
		return nil, dataErr
	}
	writing, err := strconv.ParseInt(strings.TrimSpace(rww[1]), 10, 64)
	if err != nil {
		return nil, dataErr
	}
	waiting, err := strconv.ParseInt(strings.TrimSpace(rww[1]), 10, 64)
	if err != nil {
		return nil, dataErr
	}

	metrics = &NgxMetrics{
		Active:   activeConn,
		Accepted: accepts,
		Handled:  handled,
		Reading:  reading,
		Writing:  writing,
		Waiting:  waiting,
		Requests: requests,
	}

	return metrics, nil
}

//
type userAgentRoundTripper struct {
	ua string
	rt http.RoundTripper
}

// 设置http client的header,添加User-Agent.
func (u *userAgentRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36")
	return u.rt.RoundTrip(req)
}

// InitHttpClient 初始化请求nginx status 的http client,并发送请求.
func InitHttpClient(p NgxClientParams) (client *NgxClient, err error) {
	transport := &http.Transport{}
	ut := &userAgentRoundTripper{
		ua: *p.UserAgent,
		rt: transport,
	}
	httpClient := &http.Client{
		Timeout:   p.Timeout,
		Transport: ut,
	}

	client = &NgxClient{
		endPoint:   *p.EndPoint,
		httpClient: httpClient,
	}
	return client, nil
}
