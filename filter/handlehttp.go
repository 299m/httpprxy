package filter

import (
	"fmt"
	"github.com/elazarl/goproxy"
	"net/http"
	"regexp"
	"time"
)

type Config struct {
	Logdebug       bool
	Whitelist      []string //// If whitelist is not empty, only these hosts are allowed. Recommended, unless the proxy has restricted access
	Blacklist      []string //// Ignored if whitelist is not empty. Otherwise, these hosts are blocked.
	whitelistPorts []string
	WhitelistPorts []int
}

func (c *Config) Expand() {
	c.whitelistPorts = make([]string, 0, len(c.WhitelistPorts))
	for _, port := range c.WhitelistPorts {
		c.whitelistPorts = append(c.whitelistPorts, fmt.Sprintf("%d", port))
	}
}

type Filter struct {
	cfg      *Config
	proxy    *goproxy.ProxyHttpServer
	logdebug bool
}

func NewFilter(proxy *goproxy.ProxyHttpServer, filters *Config) *Filter {

	f := &Filter{
		cfg:      filters,
		proxy:    proxy,
		logdebug: filters.Logdebug,
	}

	if filters.Logdebug {
		proxy.OnRequest().DoFunc(f.LogRequest)
	}
	if len(filters.Whitelist) > 0 {
		for _, filter := range filters.Whitelist {
			proxy.OnRequest(goproxy.DstHostIs(filter)).DoFunc(f.Allow)
			proxy.OnRequest(goproxy.DstHostIs(filter)).HandleConnectFunc(f.AllowOnConnect)
		}
		proxy.OnRequest().HandleConnectFunc(f.BlockOnConnect) //// FInal one blocks
	} else {
		for _, filter := range filters.Blacklist {
			proxy.OnRequest(goproxy.DstHostIs(filter)).DoFunc(f.Block)
			proxy.OnRequest(goproxy.DstHostIs(filter)).HandleConnectFunc(f.BlockOnConnect)
		}
		/// If we are blacklisting and it's not there, let it through
	}

	if len(filters.WhitelistPorts) > 0 {
		for _, filter := range filters.whitelistPorts {
			proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(".*:" + filter))).DoFunc(f.Allow)
			proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile(".*:" + filter))).HandleConnectFunc(f.AllowOnConnect)
		}
		proxy.OnRequest().DoFunc(f.Block)                     //// FInal one blocks
		proxy.OnRequest().HandleConnectFunc(f.BlockOnConnect) //// FInal one blocks
	}
	return f
}

func (f *Filter) LogDebug(msg ...string) {
	if f.logdebug {
		fmt.Println(msg)
	}
}

func (f *Filter) AllowOnConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	f.LogDebug("Allowing on connect ", host)
	return goproxy.OkConnect, host
}

func (f *Filter) BlockOnConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	f.LogDebug("Blocking on connect ", host)
	return goproxy.RejectConnect, host
}

func (f *Filter) Block(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	f.LogDebug("Blocking ", r.Host, r.RequestURI)
	return r, goproxy.NewResponse(r,
		goproxy.ContentTypeText, http.StatusForbidden,
		"Host is not whitelisted (or is blacklisted)!")
}

func (f *Filter) Allow(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	f.LogDebug("Allowing ", r.Host, r.RequestURI)
	return r, nil
}

func (f *Filter) LogRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	fmt.Println(time.Now(), ">", r.Host, r.RequestURI)
	return r, nil
}
