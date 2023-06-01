package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/rs/cors"
)

func NewMultipleHostsProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		//代理模式下 tracker 需要 header 增加x-real-ip 以获取真实ip
		var addr string
		if len(req.RemoteAddr) > 0 {
			tmp := strings.Split(req.RemoteAddr, ":")
			if len(tmp) > 0 {
				addr = tmp[0]
			}
		}
		req.Header.Add("x-real-ip", addr)
		log.Println("debug:", req.URL, req.RemoteAddr)

		if strings.HasPrefix(req.URL.Path, "/user") {
			req.URL.Host = fmt.Sprintf("%s:%d", "127.0.0.1", 8001)
		} else if strings.HasPrefix(req.URL.Path, "/file") {
			req.URL.Host = fmt.Sprintf("%s:%d", "127.0.0.1", 8002)
		} else {
			req.URL.Host = fmt.Sprintf("%s:%d", "127.0.0.1", 8003)
		}
		//req.URL.Path = innerURL.Path
		req.URL.Scheme = "http"
	}

	tr := http.DefaultTransport.(*http.Transport)
	tr2 := tr.Clone()
	tr2.DisableKeepAlives = false
	tr2.MaxIdleConns = 0
	tr2.MaxConnsPerHost = 0
	tr2.MaxIdleConnsPerHost = 2
	tr2.IdleConnTimeout = time.Second * 3

	return &httputil.ReverseProxy{Director: director, Transport: tr2}
}

func main() {
	proxy := NewMultipleHostsProxy()
	http.ListenAndServe(":9898", cors.AllowAll().Handler(proxy))
}
