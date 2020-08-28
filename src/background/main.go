package main

import (
	"log"
	"runtime/debug"
	"syscall/js"
	// "github.com/therecipe/qt/widgets"
)

var functi js.Func

var DEFAULT_PROXY = DIRECT_PROXY

var TOR_PROXY = map[string]interface{}{
	"type":     "socks",
	"host":     "127.0.0.1",
	"port":     "9050",
	"proxyDNS": true,
}

var I2P_PROXY = map[string]interface{}{
	"type":     "socks",
	"host":     "127.0.0.1",
	"port":     "9052",
	"proxyDNS": true,
}

var DIRECT_PROXY = map[string]interface{}{
	"type": "direct",
}

var pse *ProxySwitcherExtension

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Print("extension main iteration have crushed:")
			log.Print("      message: ", r)
			log.Print("        stack: \n", string(debug.Stack()))
		}
	}()

	log.SetFlags(log.Lshortfile)

	log.Println("ProxySwitcherExtension init")

	pse = NewProxySwitcherExtension()

	g := js.Global()

	b := g.Get("browser")

	b.Get("webRequest").Get("onBeforeRequest").Call(
		"addListener",
		js.FuncOf(pse.BrowserWebRequestOnBeforeRequestHandler),
		map[string]interface{}{
			"urls": []interface{}{"<all_urls>"},
		},
		[]interface{}{"blocking"},
	)

	b.Get("proxy").Get("onRequest").Call(
		"addListener",
		js.FuncOf(pse.BrowserProxyOnRequestHandler),
		map[string]interface{}{
			"urls": []interface{}{"<all_urls>"},
		},
	)

	b.Get("tabs").Get("onActivated").Call(
		"addListener",
		js.FuncOf(pse.BrowserTabsOnActivatedHandler),
	)

	b.Get("browserAction").Get("onClicked").Call(
		"addListener",
		js.FuncOf(pse.ShowMainWindow),
	)

	g.Set(
		"pse",
		map[string]interface{}{
			"renderMainWindow": js.FuncOf(pse.RenderMainWindow),
		},
	)

	// without this, code will become unavailable

	// tik := true
	// for {
	// 	time.Sleep(time.Second)
	// 	if tik {
	// 		log.Print("tik")
	// 	} else {
	// 		log.Print("tak")
	// 	}
	// 	tik = !tik
	// }
	c := make(chan bool)
	<-c

}
