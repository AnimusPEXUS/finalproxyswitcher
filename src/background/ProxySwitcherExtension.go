package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/url"
	"runtime/debug"
	"strings"
	"syscall/js"
	"time"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
	pexu_promise "github.com/AnimusPEXUS/wasmtools/promise"
)

type ProxySwitcherExtension struct {
	request_history *RequestHistory
	config          *ConfigModel
	current_tab_id  int
	changed         bool
	main_window     *MainWindow
}

func NewProxySwitcherExtension() *ProxySwitcherExtension {
	self := &ProxySwitcherExtension{
		request_history: NewRequestHistory(),
		current_tab_id:  -1,
		config:          &ConfigModel{},
	}
	self.LoadConfig()
	if self.config.ProxyTargets == nil {
		self.config.ProxyTargets = map[string]*ProxyInfo{}
	}
	return self
}

func (self *ProxySwitcherExtension) LoadConfig() error {
	g := js.Global()

	config_promise_js := g.Get("browser").Get("storage").Get("local").Call(
		"get",
		"config",
	)

	psucc := make(chan bool)
	perr := make(chan bool)

	config_promise, err := pexu_promise.NewPromiseFromJSValue(config_promise_js)
	if err != nil {
		return err
	}

	var config_bytes string

	config_promise.Then(
		js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			ts := args[0].Get("config")
			if ts == js.Undefined() {
				log.Println("Then succ - got undefined")
				perr <- true
				return false
			}
			s := ts.String()
			log.Println("Then succ", s)
			config_bytes = s
			psucc <- true
			return false
		},
		),
		js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			perr <- true
			return false
		},
		),
	)

	select {
	case <-psucc:

	case <-perr:
		log.Println("error loading config from browser")
		return errors.New("error loading config from browser")
	case <-time.After(time.Duration(time.Minute)):
		log.Println("timeout loading config from browser")
		return errors.New("timeout loading config from browser")
	}

	err = json.Unmarshal([]byte(config_bytes), self.config)
	if err != nil {
		return err
	}

	return nil
}

func (self *ProxySwitcherExtension) SaveConfig() error {
	g := js.Global()

	b, err := json.Marshal(self.config)
	if err != nil {
		return err
	}

	config_promise_js := g.Get("browser").Get("storage").Get("local").Call(
		"set",
		map[string]interface{}{"config": string(b)},
	)

	psucc := make(chan bool)
	perr := make(chan bool)

	config_promise, err := pexu_promise.NewPromiseFromJSValue(config_promise_js)
	if err != nil {
		return err
	}

	config_promise.Then(
		js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			psucc <- true
			return false
		},
		),
		js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			perr <- true
			return false
		},
		),
	)

	select {
	case <-psucc:
		log.Println("config saved")
		self.changed = false
	case <-perr:
		log.Println("config saving error")
		return errors.New("config saving error")
	case <-time.After(time.Duration(time.Minute)):
		log.Println("timeout config saving")
		return errors.New("timeout config saving")
	}

	return nil
}

func (self *ProxySwitcherExtension) BrowserProxyOnRequestHandler(
	this js.Value,
	args []js.Value,
) interface{} {

	if len(args) != 1 {
		return DEFAULT_PROXY
	}

	url_s := args[0].Get("url").String()

	url_p, err := url.Parse(url_s)
	if err != nil {
		return DEFAULT_PROXY
	}

	for _, i := range []string{
		".onion",
		".animespirit.ru",
		".anilibria.tv",
		".slack.com",
		".slack-edge.com",
		"design.firefox.com",
		".xvideos.com",
		"hentai4manga.com",
		".golang.org",
		".origin.com",
		".ea.com",
		".facebook.com",
		".fbcdn.net",
		".thepiratebay.org",
		".xuk.life",
		".lurkmore.to",
	} {

		c := i

		subdomains := c[0] == '.'

		if subdomains {
			c = c[1:]
		}

		if c == url_p.Host || (subdomains && strings.HasSuffix(url_p.Host, "."+c)) {
			return TOR_PROXY
		}
	}

	return DIRECT_PROXY
}

func (self *ProxySwitcherExtension) BrowserWebRequestOnBeforeRequestHandler(
	this js.Value,
	args []js.Value,
) interface{} {

	ret := map[string]interface{}{}

	if len(args) != 1 {
		ret["cancel"] = true
		return ret
	}

	self.request_history.AddFromMozillaObject(args[0], true)

	return ret
}

func (self *ProxySwitcherExtension) ShowMainWindow(
	this js.Value,
	args []js.Value,
) interface{} {

	defer func() {
		if r := recover(); r != nil {
			log.Print("main window renderer have crushed:")
			log.Print("      message: ", r)
			log.Print("        stack: \n", string(debug.Stack()))
		}
	}()

	g := js.Global()

	g.Get("browser").Get("windows").Call(
		"create",
		map[string]interface{}{
			"url":  "mainwindow.html",
			"type": "popup",
		},
	)

	return false

}

func (self *ProxySwitcherExtension) RenderMainWindow(
	this js.Value,
	args []js.Value,
) interface{} {

	defer func() {
		if r := recover(); r != nil {
			log.Print("main window renderer have crushed:")
			log.Print("      message: ", r)
			log.Print("        stack: \n", string(debug.Stack()))
		}
	}()

	if len(args) != 0 {

		window_document := args[0].Get("document")

		pexu_dom_document := &pexu_dom.Document{window_document}

		etc := elementtreeconstructor.NewElementTreeConstructor(pexu_dom_document)

		window := self.MainWindow(pexu_dom_document)

		etc.ReplaceChildren([]pexu_dom.ToNodeConvertable{window.Element})

	}

	if self.current_tab_id > -1 {
		tab_id := self.current_tab_id
		log.Println("active tab id", tab_id)

		hosts := self.request_history.ComputeTabHosts(tab_id)

		b, err := json.MarshalIndent(hosts, "  ", "  ")
		if err != nil {
			log.Println("err", err)
			return nil
		}

		log.Println("tab hosts", string(b))
	}

	return nil
}

func (self *ProxySwitcherExtension) BrowserTabsOnActivatedHandler(this js.Value, args []js.Value) interface{} {
	t := args[0].Get("tabId").Int()
	if t > -1 {
		self.current_tab_id = t
	}
	return nil
}

func (self *ProxySwitcherExtension) Changed() {
	self.changed = true
	go self.main_window.Changed()
}
