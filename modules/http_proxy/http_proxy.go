package http_proxy

import (
	"github.com/bettercap/bettercap/session"
)

type HttpProxy struct {
	session.SessionModule
	proxy *HTTPProxy
}

func NewHttpProxy(s *session.Session) *HttpProxy {
	mod := &HttpProxy{
		SessionModule: session.NewSessionModule("http.proxy", s),
		proxy:         NewHTTPProxy(s),
	}

	mod.AddParam(session.NewIntParameter("http.port",
		"80",
		"HTTP port to redirect when the proxy is activated."))

	mod.AddParam(session.NewStringParameter("http.proxy.address",
		session.ParamIfaceAddress,
		session.IPv4Validator,
		"Address to bind the HTTP proxy to."))

	mod.AddParam(session.NewIntParameter("http.proxy.port",
		"8080",
		"Port to bind the HTTP proxy to."))

	mod.AddParam(session.NewStringParameter("http.proxy.script",
		"",
		"",
		"Path of a proxy JS script."))

	mod.AddParam(session.NewStringParameter("http.proxy.injectjs",
		"",
		"",
		"URL, path or javascript code to inject into every HTML page."))

	mod.AddParam(session.NewBoolParameter("http.proxy.sslstrip",
		"false",
		"Enable or disable SSL stripping."))

	mod.AddHandler(session.NewModuleHandler("http.proxy on", "",
		"Start HTTP proxy.",
		func(args []string) error {
			return mod.Start()
		}))

	mod.AddHandler(session.NewModuleHandler("http.proxy off", "",
		"Stop HTTP proxy.",
		func(args []string) error {
			return mod.Stop()
		}))

	return mod
}

func (mod *HttpProxy) Name() string {
	return "http.proxy"
}

func (mod *HttpProxy) Description() string {
	return "A full featured HTTP proxy that can be used to inject malicious contents into webpages, all HTTP traffic will be redirected to it."
}

func (mod *HttpProxy) Author() string {
	return "Simone Margaritelli <evilsocket@gmail.com>"
}

func (mod *HttpProxy) Configure() error {
	var err error
	var address string
	var proxyPort int
	var httpPort int
	var scriptPath string
	var stripSSL bool
	var jsToInject string

	if mod.Running() {
		return session.ErrAlreadyStarted
	} else if err, address = mod.StringParam("http.proxy.address"); err != nil {
		return err
	} else if err, proxyPort = mod.IntParam("http.proxy.port"); err != nil {
		return err
	} else if err, httpPort = mod.IntParam("http.port"); err != nil {
		return err
	} else if err, scriptPath = mod.StringParam("http.proxy.script"); err != nil {
		return err
	} else if err, stripSSL = mod.BoolParam("http.proxy.sslstrip"); err != nil {
		return err
	} else if err, jsToInject = mod.StringParam("http.proxy.injectjs"); err != nil {
		return err
	}

	return mod.proxy.Configure(address, proxyPort, httpPort, scriptPath, jsToInject, stripSSL)
}

func (mod *HttpProxy) Start() error {
	if err := mod.Configure(); err != nil {
		return err
	}

	return mod.SetRunning(true, func() {
		mod.proxy.Start()
	})
}

func (mod *HttpProxy) Stop() error {
	return mod.SetRunning(false, func() {
		mod.proxy.Stop()
	})
}