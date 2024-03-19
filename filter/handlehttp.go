package filter

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
	"time"
)

type Config struct {
	Logdebug  bool
	Whitelist []string //// If whitelist is not empty, only these hosts are allowed. Recommended, unless the proxy has restricted access
	Blacklist []string //// Ignored if whitelist is not empty. Otherwise, these hosts are blocked.
}

func (c *Config) Expand() {

}

type Filter struct {
	cfg   *Config
	proxy *goproxy.ProxyHttpServer
}

func NewFilter(proxy *goproxy.ProxyHttpServer, filters *Config) *Filter {

	f := &Filter{
		cfg:   filters,
		proxy: proxy,
	}

	if filters.Logdebug {
		proxy.OnRequest().DoFunc(f.LogRequest)
	}

	if len(filters.Whitelist) > 0 {
		for _, filter := range filters.Whitelist {
			proxy.OnRequest(goproxy.DstHostIs(filter)).DoFunc(f.Allow)
		}
		proxy.OnRequest().DoFunc(f.Block) //// FInal one blocks
	} else {
		for _, filter := range filters.Blacklist {
			proxy.OnRequest(goproxy.DstHostIs(filter)).DoFunc(f.Block)
		}
		/// If we are blacklisting and it's not there, let it through
	}
	return f
}

func (f *Filter) Block(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Host is not whitelisted (or is blacklisted)!")
}

func (f *Filter) Allow(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	return r, nil
}

func (f *Filter) LogRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	fmt.Println(time.Now(), ">", r.Host, r.RequestURI)
	return r, nil
}
