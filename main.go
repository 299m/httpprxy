package main

import (
	"flag"
	"fmt"
	"github.com/299m/util/util"
	"github.com/elazarl/goproxy"
	"httpprxy/filter"
	"net/http"
	"os"
)

type TlsConfig struct {
	Cert string
	Key  string
	Port string
}

func (t *TlsConfig) Expand() {
	t.Cert = os.ExpandEnv(t.Cert)
	t.Key = os.ExpandEnv(t.Key)
}

func main() {
	defer util.OnPanicFunc()
	basedir := flag.String("cfgdir", "", "The directory where the configuration files are located")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	tlscfg := &TlsConfig{}
	filters := &filter.Config{}

	cfg := map[string]util.Expandable{
		"tls":    tlscfg,
		"filter": filters,
	}
	util.ReadConfig(*basedir, cfg)
	filter.NewFilter(proxy, filters)

	fmt.Println("Starting proxy on port "+tlscfg.Port, " with cert ", tlscfg.Cert, " and key ", tlscfg.Key)

	util.CheckError(http.ListenAndServeTLS(":"+tlscfg.Port, tlscfg.Cert, tlscfg.Key, proxy))
}
