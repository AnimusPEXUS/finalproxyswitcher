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

	"github.com/AnimusPEXUS/gojstools/elementtreeconstructor"
	pexu_dom "github.com/AnimusPEXUS/gojswebapi/dom"
	pexu_promise "github.com/AnimusPEXUS/gojswebapi/promise"
	"github.com/AnimusPEXUS/utils/domainname"
)

// TODO: currently, rule sets are based only on domain names. it is possible,
//       site destinguishing by ports/schemas/etc are also required

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
		config:          NewConfigModel(),
	}

	self.LoadConfig()

	if self.config.ProxyTargets == nil {
		self.config.ProxyTargets = map[string]*ProxyTarget{}
	}

	if self.config.RuleSet == nil {
		self.config.RuleSet = map[string]*DomainSettings{}
	}

	return self
}

func (self *ProxySwitcherExtension) GetStorageLocalValue(name string) (string, error) {
	// TODO: move this to wasmtools

	g := js.Global()

	config_promise_js := g.Get("browser").Get("storage").Get("local").Call(
		"get",
		name,
	)

	config_promise, err := pexu_promise.NewPromiseFromJSValue(&config_promise_js)
	if err != nil {
		return "", err
	}

	var res_string string

	psucc := make(chan bool)
	perr := make(chan bool)

	config_promise.Then(
		&[]js.Func{js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			if len(args) == 0 {
				perr <- true
				return false
			}
			res_string = args[0].Get(name).String()
			psucc <- true
			return false
		},
		)}[0],
		&[]js.Func{js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			perr <- true
			return false
		},
		)}[0],
	)

	select {
	case <-psucc:
		log.Println("storage data loaded ok")
		return res_string, nil
	case <-perr:
		log.Println("error loading storage data from browser")
		return "", errors.New("error loading storage data from browser")
	case <-time.After(time.Duration(time.Minute)):
		log.Println("timeout loading storage data from browser")
		return "", errors.New("timeout loading storage data from browser")
	}

	return "", errors.New("unknown error")
}

func (self *ProxySwitcherExtension) UseActiveConfigJSON(config_string string) error {

	err := json.Unmarshal([]byte(config_string), &self.config)
	if err != nil {
		return err
	}

	self.config.Fix()

	self.Changed()

	return nil
}

func (self *ProxySwitcherExtension) GenerateActiveConfigJSON(use_indent bool) (string, error) {
	return self.GenerateJSON(self.config, use_indent)
}

func (self *ProxySwitcherExtension) GenerateSavedConfigJSON(use_indent bool) (string, error) {

	config_string, err := self.GetStorageLocalValue("config")
	if err != nil {
		return "", err
	}

	var t ConfigModel

	err = json.Unmarshal([]byte(config_string), &t)
	if err != nil {
		return "", err
	}

	ret, err := self.GenerateJSON(t, use_indent)
	if err != nil {
		return "", err
	}

	return ret, nil
}

func (self *ProxySwitcherExtension) GenerateJSON(data interface{}, use_indent bool) (string, error) {

	var b []byte
	var err error

	if use_indent {
		b, err = json.MarshalIndent(self.config, "  ", "  ")
		if err != nil {
			return "", err
		}
	} else {
		b, err = json.Marshal(self.config)
		if err != nil {
			return "", err
		}
	}

	return string(b), nil
}

func (self *ProxySwitcherExtension) LoadConfig() error {

	config_string, err := self.GetStorageLocalValue("config")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(config_string), self.config)
	if err != nil {
		return err
	}

	self.config.Fix()

	return nil
}

func (self *ProxySwitcherExtension) SaveConfig() error {
	g := js.Global()

	active_settings_json, err := self.GenerateActiveConfigJSON(true)
	if err != nil {
		return err
	}

	log.Println(active_settings_json)

	config_promise_js := g.Get("browser").Get("storage").Get("local").Call(
		"set",
		map[string]interface{}{"config": active_settings_json},
	)

	psucc := make(chan bool)
	perr := make(chan bool)

	config_promise, err := pexu_promise.NewPromiseFromJSValue(&config_promise_js)
	if err != nil {
		return err
	}

	config_promise.Then(
		&[]js.Func{js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			psucc <- true
			return false
		},
		)}[0],
		&[]js.Func{js.FuncOf(func(
			this js.Value,
			args []js.Value,
		) interface{} {
			perr <- true
			return false
		},
		)}[0],
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

		".thepiratebay.org",

		".protonmail.com",

		".google.com",
		".googleapis.com",

		".telegra.ph",
		".t.me",
		".telegram.org",

		".linkedin.com",

		".animespirit.ru",
		".animespirit.online",
		".animespirit.cc",

		".anilibria.tv",

		".golang.org",
		"design.firefox.com",

		".intel.com",
		".drone.io",
		".allelectronics.com",

		".dub.pm",

		".github.com",
		".githubassets.com",
		".githubusercontent.com",

		".slack.com",
		".slack-edge.com",
		".slack-redir.net",

		".origin.com",
		".ea.com",

		".facebook.com",
		".fbcdn.net",

		".xvideos.com",
		".hentai4manga.com",
		".xuk.life",

		".lurkmore.to",
		".vk.com",
		".mail.ru",

		".opennet.ru",
		".linux.org.ru",
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

	for _, i := range []string{
		".i2p",
	} {

		c := i

		subdomains := c[0] == '.'

		if subdomains {
			c = c[1:]
		}

		if c == url_p.Host || (subdomains && strings.HasSuffix(url_p.Host, "."+c)) {
			return I2P_PROXY
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

	record, _ := self.request_history.AddFromMozillaObject(args[0])

	log.Println(record.String())

	// if strings.HasPrefix(record.RequestId, "fakeRequest-") {
	// 	return ret
	// }

	u, err := url.Parse(record.URL)
	if err != nil {
		// TODO: block request
	}

	request_domain := ""
	subrequest_domain := ""

	if record.DocumentURL == nil {
		request_domain = u.Hostname()
	} else {
		// record_main := self.request_history.TabIdGetMainRequestEntry(record.TabId)
		// if record_main == nil {
		// 	log.Println("record_main is nil")
		// 	return ret
		// }
		u2, err := url.Parse(*record.DocumentURL)
		if err != nil {
			// TODO: block request
		}
		request_domain = u2.Hostname()
		subrequest_domain = u.Hostname()
	}

	// if record.DocumentURL == nil {
	// 	log.Println("url", record.URL, "no doc_url")
	// } else {
	// 	log.Println("url", record.URL, "doc_url", *record.DocumentURL)
	// }

	self.CalculateCurrentRules(request_domain, subrequest_domain)

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

	// TODO: somehow use window.open (or something) to hide window's titlebar
	// prom_js := g.Get("browser").Get("windows").Call(
	// 	"create",
	// 	map[string]interface{}{
	// 		// "url":  "mainwindow.html",
	// 		"type": "popup",
	// 	},
	// )

	// prom, err := pexu_promise.NewPromiseFromJSValue(prom_js)
	// if err != nil {
	// 	return false
	// }

	// prom.Then(
	// 	js.FuncOf(
	// 		func(
	// 			this js.Value,
	// 			args []js.Value,
	// 		) interface{} {
	// 			if len(args) == 0 {
	// 				return false
	// 			}
	// 			args[0].Call(
	// 				"open",
	// 				"mainwindow.html",
	// 				"",
	// 				map[string]interface{}{
	// 					"titlebar": "0",
	// 				},
	// 			)
	// 			return false
	// 		},
	// 	),
	// )

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

		pexu_dom_document := &pexu_dom.Document{&window_document}

		etc := elementtreeconstructor.NewElementTreeConstructor(pexu_dom_document)

		window := self.MainWindow(etc)

		etc.ReplaceChildren([]pexu_dom.ToNodeConvertable{window.Element.Element})

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

func (self *ProxySwitcherExtension) BrowserTabsOnActivatedHandler(
	this js.Value,
	args []js.Value,
) interface{} {
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

// func (self *ProxySwitcherExtension) ProxyTargetListKysOnly() []string {

// 	// TODO: is this optimal? gues not

// 	ret := make([]string, 0)

// 	for k, _ := range self.config.ProxyTargets {
// 		ret = append(ret, k)
// 	}

// 	sort.Strings(ret)

// 	return ret
// }

func (self *ProxySwitcherExtension) ProxyTargetList() [][2]string {

	ret := make([][2]string, 0)

	ret = append(ret, [2]string{"", "Undefined"})

	for k, _ := range self.config.ProxyTargets {
		ret = append(ret, [2]string{k, k})
	}

	return ret
}

type SettingsStruct struct {
	// TODO: try using this shortcut if preformance needed
	// domain           string
	domainDomainName *domainname.DomainName
	rules            *DomainSettings
}

func (self *ProxySwitcherExtension) CalculateCurrentRules(
	request_host string,
	subrequest_host string,
) *Rules {

	// TODO: treat IDN

	request_host_s := domainname.NewDomainNameFromString(request_host)
	subrequest_host_s := domainname.NewDomainNameFromString(subrequest_host)

	log.Println("CalculateCurrentRules")
	log.Println("  main request", request_host_s.String(), "sub request", subrequest_host_s.String())

	matching_settings := make([]*SettingsStruct, 0)

	for _, v := range self.config.RuleSet {
		if request_host_s.IsEqualTo(v.Domain) ||
			(v.RulesAndInheritance.ApplyToSubdomains &&
				request_host_s.IsSubdomainTo(v.Domain)) {
			matching_settings = append(
				matching_settings,
				&SettingsStruct{
					domainDomainName: v.Domain,
					rules:            v,
				},
			)
		}
	}

	if len(matching_settings) != 0 {
		for i := 0; i != len(matching_settings)-1; i++ {
			for j := i + 1; j != len(matching_settings); j++ {
				if matching_settings[i].domainDomainName.CompareTo(matching_settings[j].domainDomainName) > 0 {
					z := matching_settings[i]
					matching_settings[i] = matching_settings[j]
					matching_settings[j] = z
				}
			}
		}
	}

	log.Printf("matching settings %d:", len(matching_settings))
	for _, i := range matching_settings {
		log.Println("   ", i.domainDomainName.String())
	}

	ret := &Rules{}

	if len(matching_settings) != 0 {
		// TODO: add error handeling
		self.CalculateCurrentRulesRulePart(
			request_host,
			subrequest_host,
			ret,
			0,
			matching_settings,
			// matching_settings[len(matching_settings)-1].domainDomainName.String(),
		)

		// TODO: add error handeling
		self.CalculateCurrentRulesRulePart(
			request_host,
			subrequest_host,
			ret,
			1,
			matching_settings,
			// matching_settings[len(matching_settings)-1].domainDomainName.String(),
		)

		// TODO: add error handeling
		self.CalculateCurrentRulesRulePart(
			request_host,
			subrequest_host,
			ret,
			2,
			matching_settings,
			// matching_settings[len(matching_settings)-1].domainDomainName.String(),
		)
	}

	b, _ := json.MarshalIndent(ret, "  ", "  ")

	log.Printf("Calculated Rule: %s", string(b))
	return ret
}

func (self *ProxySwitcherExtension) CalculateCurrentRulesRulePart(
	request_host string,
	subrequest_host string,
	rules_structure *Rules,
	mode int, // TODO: use named constants here
	matched_domain_setting_structs_slice []*SettingsStruct,
	// start_with_domain string,
) error {

	var (
		str_in_q   *SettingsStruct
		str_in_q_i int
	)

	str_in_q_i = len(matched_domain_setting_structs_slice)

loop0:
	str_in_q_i--

	if str_in_q_i < 0 {
		switch mode {
		default:
			log.Println("programming error")
			return errors.New("programming error")
		case 0:
			rules_structure.HttpRule = self.config.RootRules.RulesAndInheritance.Rules.HttpRule
		case 1:
			rules_structure.RequestRule = self.config.RootRules.RulesAndInheritance.Rules.RequestRule
		case 2:
			rules_structure.ProxyRule = self.config.RootRules.RulesAndInheritance.Rules.ProxyRule
			rules_structure.ProxyTarget = self.config.RootRules.RulesAndInheritance.Rules.ProxyTarget
		}
		return nil
	}

	str_in_q = matched_domain_setting_structs_slice[str_in_q_i]

	if (mode == 0 && str_in_q.rules.RulesAndInheritance.Rules.HttpRule == HttpRuleUndefined) ||
		(mode == 1 && str_in_q.rules.RulesAndInheritance.Rules.RequestRule == RequestRuleUndefined) ||
		(mode == 2 && str_in_q.rules.RulesAndInheritance.Rules.ProxyRule == ProxyRuleUndefined) {

		sw_value := RuleInheritance(0)
		switch mode {
		default:
			log.Println("programming error")
			return errors.New("programming error")
		case 0:
			sw_value = str_in_q.rules.RulesAndInheritance.RulesInheritance.HttpRuleInheritance
		case 1:
			sw_value = str_in_q.rules.RulesAndInheritance.RulesInheritance.RequestRuleInheritance
		case 2:
			sw_value = str_in_q.rules.RulesAndInheritance.RulesInheritance.ProxyRuleInheritance
		}

		switch sw_value {
		default:
			log.Println("programming error")
			return errors.New("programming error")
		case RuleInheritanceNone:
			fallthrough
		case RuleInheritanceGlobal:
			str_in_q_i = 0
			fallthrough
		case RuleInheritanceParent:
			goto loop0
		}

	}

	switch mode {
	default:
		log.Println("programming error")
		return errors.New("programming error")
	case 0:
		rules_structure.HttpRule = str_in_q.rules.RulesAndInheritance.Rules.HttpRule
	case 1:
		rules_structure.RequestRule = str_in_q.rules.RulesAndInheritance.Rules.RequestRule
	case 2:
		rules_structure.ProxyRule = str_in_q.rules.RulesAndInheritance.Rules.ProxyRule
		rules_structure.ProxyTarget = str_in_q.rules.RulesAndInheritance.Rules.ProxyTarget
	}

	return nil
}
