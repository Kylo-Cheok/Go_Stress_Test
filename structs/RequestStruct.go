package structs

import (
	"io"
	"strings"
)

/**
 * @Author: kylo_cheok
 * @Email:  maggic0816@gmail.com
 * @Date:   2022/2/19 17:13
 * @Desc:   Grace under pressure
 */

type Request struct {
	URL     string            // URL
	Method  string            // 方法 GET/POST
	Headers map[string]string // Headers
	Body    string            // body
	MaxCon  int               // 每个连接的请求数
}

// GetBody 将body转换成request可用的格式
func (r *Request) GetBody() (body io.Reader) {
	return strings.NewReader(r.Body)
}

type RequestResults struct {
	ID       int     // 消息ID
	RT       float64 // 请求时间
	RespData string  //响应文本
	Succeed  bool    // 是否请求成功
}
