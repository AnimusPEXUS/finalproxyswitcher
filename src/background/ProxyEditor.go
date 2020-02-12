package main

import (
	"log"
	"syscall/js"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
)

type ProxyTargetEditor struct {
	setting_name string

	edit_b, cancel_b, save_b, load_b, delete_b *elementtreeconstructor.ElementMutator
	changed_asterisk                           *elementtreeconstructor.ElementMutator

	type_select                        *elementtreeconstructor.ElementMutator
	name_text                          *elementtreeconstructor.ElementMutator
	host_text_cb                       *elementtreeconstructor.ElementMutator
	host_text                          *elementtreeconstructor.ElementMutator
	port_text_cb                       *elementtreeconstructor.ElementMutator
	port_text                          *elementtreeconstructor.ElementMutator
	username_text_cb                   *elementtreeconstructor.ElementMutator
	username_text                      *elementtreeconstructor.ElementMutator
	password_text_cb                   *elementtreeconstructor.ElementMutator
	password_text                      *elementtreeconstructor.ElementMutator
	proxydns_cb_cb                     *elementtreeconstructor.ElementMutator
	proxydns_cb                        *elementtreeconstructor.ElementMutator
	failover_timeout_text_cb           *elementtreeconstructor.ElementMutator
	failover_timeout_text              *elementtreeconstructor.ElementMutator
	proxy_authorization_header_text_cb *elementtreeconstructor.ElementMutator
	proxy_authorization_header_text    *elementtreeconstructor.ElementMutator
	connection_isolation_key_text_cb   *elementtreeconstructor.ElementMutator
	connection_isolation_key_text      *elementtreeconstructor.ElementMutator

	view_div *elementtreeconstructor.ElementMutator
	edit_div *elementtreeconstructor.ElementMutator

	Element   *pexu_dom.Element
	extension *ProxySwitcherExtension
}

func (self *ProxySwitcherExtension) ProxyTargetEditor(
	setting_name string,
	editing_switch_possible bool,
	editing_mode bool,
	document *pexu_dom.Document,
	delete_cb func(),
) *ProxyTargetEditor {
	ret := NewProxyTargetEditor(
		setting_name,
		editing_switch_possible,
		editing_mode,
		document,
		delete_cb,
		self,
	)
	return ret
}

func NewProxyTargetEditor(
	setting_name string,
	editing_switch_possible bool,
	editing_mode bool,
	document *pexu_dom.Document,
	delete_cb func(),
	extension *ProxySwitcherExtension,
) *ProxyTargetEditor {

	self := &ProxyTargetEditor{}

	self.setting_name = setting_name
	self.extension = extension

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	self_Element_Mutator := etc.CreateElement("div").
		ExternalUse(applyBorder).
		SetStyle("margin", "1px").
		SetStyle("padding", "3px").
		AppendChildren(

			etc.CreateElement("form").
				AppendChildren(

					etc.CreateElement("span").
						AssignSelf(&self.view_div).
						AppendChildren(

							etc.CreateElement("button").
								ExternalUse(applyButtonStyle).
								AppendChildren(
									etc.CreateTextNode("Edit"),
								).
								AssignSelf(&self.edit_b),
						),

					etc.CreateElement("span").
						// SetStyle("overflow-wrap", " break-word").
						AssignSelf(&self.edit_div).
						AppendChildren(

							etc.CreateElement("button").
								ExternalUse(applyButtonStyle).
								AppendChildren(
									etc.CreateTextNode("Cancel"),
								).
								AssignSelf(&self.cancel_b),

							etc.CreateElement("button").
								ExternalUse(applyButtonStyle).
								AppendChildren(
									etc.CreateTextNode("Load"),
								).
								AssignSelf(&self.load_b),

							etc.CreateElement("button").
								ExternalUse(applyButtonStyle).
								AppendChildren(
									etc.CreateTextNode("Save"),
									etc.CreateElement("span").
										ExternalUse(applySpanChangedAsterisk).
										AppendChildren(
											etc.CreateTextNode("*"),
										).
										AssignSelf(&self.changed_asterisk),
								).
								AssignSelf(&self.save_b),

							etc.CreateElement("button").
								ExternalUse(applyButtonStyle).
								AppendChildren(
									etc.CreateTextNode("Delete"),
								).
								AssignSelf(&self.delete_b),

							etc.CreateElement("select").
								ExternalUse(applyButtonStyle).
								Set("title", "type").
								AppendChildren(

									etc.CreateElement("option").
										AppendChildren(
											etc.CreateTextNode("direct"),
										).
										Set("value", "direct"),

									etc.CreateElement("option").
										AppendChildren(
											etc.CreateTextNode("http"),
										).
										Set("value", "http"),

									etc.CreateElement("option").
										AppendChildren(
											etc.CreateTextNode("https"),
										).
										Set("value", "https"),

									etc.CreateElement("option").
										AppendChildren(
											etc.CreateTextNode("socks"),
										).
										Set("value", "socks"),

									etc.CreateElement("option").
										AppendChildren(
											etc.CreateTextNode("socks4"),
										).
										Set("value", "socks4"),
								).
								AssignSelf(&self.type_select),

							etc.CreateElement("span").
								Set("title", "Name").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "text").
										SetStyle("border", "none").
										AssignSelf(&self.name_text),
								),

							etc.CreateElement("span").
								Set("title", "Host").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.host_text_cb),

									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.host_text),
								),

							etc.CreateElement("span").
								Set("title", "Port").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.port_text_cb),

									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.port_text),
								),

							etc.CreateElement("span").
								Set("title", "UserName").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.username_text_cb),

									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.username_text),
								),

							etc.CreateElement("span").
								Set("title", "User Password").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.password_text_cb),

									etc.CreateElement("input").
										Set("type", "password").
										AssignSelf(&self.password_text),
								),

							etc.CreateElement("span").
								Set("title", "Use Proxy for DNS").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.proxydns_cb_cb),

									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.proxydns_cb),
								),

							etc.CreateElement("span").
								Set("title", "Failover Timeout").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.failover_timeout_text_cb),

									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.failover_timeout_text),
								),

							etc.CreateElement("span").
								Set("title", "Proxy Authorization Header").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.proxy_authorization_header_text_cb),

									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.proxy_authorization_header_text),
								),

							etc.CreateElement("span").
								Set("title", "Connection Isolation Key").
								ExternalUse(applyProxyEditorSpanStyle).
								AppendChildren(
									etc.CreateElement("input").
										Set("type", "checkbox").
										AssignSelf(&self.connection_isolation_key_text_cb),
									etc.CreateElement("input").
										Set("type", "text").
										AssignSelf(&self.connection_isolation_key_text),
								),
						),
				),
		)

	self.Element = self_Element_Mutator.Element

	loadTarget := func() {

		name := self.setting_name

		info, ok := self.extension.config.ProxyTargets[name]
		if !ok {
			if name != "" {
				log.Printf("proxy target named '%s' is not found in config\n", name)
			}
			return
		}

		self.name_text.Set("value", name)

		self.type_select.Set("value", info.Type)

		self.host_text_cb.Set("checked", info.Host != nil)
		if info.Host != nil {
			self.host_text.Set("value", *info.Host)
		}

		self.port_text_cb.Set("checked", info.Port != nil)
		if info.Port != nil {
			self.port_text.Set("value", *info.Port)
		}

		self.username_text_cb.Set("checked", info.Username != nil)
		if info.Username != nil {
			self.username_text.Set("value", *info.Username)
		}

		self.password_text_cb.Set("checked", info.Password != nil)
		if info.Password != nil {
			self.password_text.Set("value", *info.Password)
		}

		self.proxydns_cb_cb.Set("checked", info.ProxyDNS != nil)
		if info.ProxyDNS != nil {
			self.proxydns_cb.Set("checked", *info.ProxyDNS)
		}

		self.failover_timeout_text_cb.Set("checked", info.FailoverTimeout != nil)
		if info.FailoverTimeout != nil {
			self.failover_timeout_text.Set("value", *info.FailoverTimeout)
		}

		self.proxy_authorization_header_text_cb.Set("checked", info.ProxyAuthorizationHeader != nil)
		if info.ProxyAuthorizationHeader != nil {
			self.proxy_authorization_header_text.Set("value", *info.ProxyAuthorizationHeader)
		}

		self.connection_isolation_key_text_cb.Set("checked", info.ConnectionIsolationKey != nil)
		if info.ConnectionIsolationKey != nil {
			self.connection_isolation_key_text.Set("value", *info.ConnectionIsolationKey)
		}

	}

	loadTarget_j := func(
		this js.Value,
		args []js.Value,
	) interface{} {
		loadTarget()
		return false
	}

	self.load_b.Set(
		"onclick",
		js.FuncOf(loadTarget_j),
	)

	self.save_b.Set(
		"onclick",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				name := self.name_text.Element.Get("value").String()
				self.setting_name = name

				info := &ProxyTarget{}
				info.Type = self.type_select.Element.Get("value").String()

				if self.host_text_cb.GetJsValue("checked").Bool() {
					info.Host = &[]string{self.host_text.GetJsValue("value").String()}[0]
				}

				if self.port_text_cb.GetJsValue("checked").Bool() {
					info.Port = &[]string{self.port_text.GetJsValue("value").String()}[0]
				}

				if self.username_text_cb.GetJsValue("checked").Bool() {
					info.Username = &[]string{self.username_text.GetJsValue("value").String()}[0]
				}

				if self.password_text_cb.GetJsValue("checked").Bool() {
					info.Password = &[]string{self.password_text.GetJsValue("value").String()}[0]
				}

				if self.proxydns_cb_cb.GetJsValue("checked").Bool() {
					info.ProxyDNS = &[]bool{self.proxydns_cb.GetJsValue("checked").Bool()}[0]
				}

				if self.failover_timeout_text_cb.GetJsValue("checked").Bool() {
					info.FailoverTimeout = &[]int{self.failover_timeout_text.GetJsValue("value").Int()}[0]
				}

				if self.proxy_authorization_header_text_cb.GetJsValue("checked").Bool() {
					info.ProxyAuthorizationHeader = &[]string{self.proxy_authorization_header_text.GetJsValue("value").String()}[0]
				}

				if self.connection_isolation_key_text_cb.GetJsValue("checked").Bool() {
					info.ConnectionIsolationKey = &[]string{self.connection_isolation_key_text.GetJsValue("value").String()}[0]
				}

				self.extension.config.ProxyTargets[name] = info

				go self.extension.Changed()
				self.changed_asterisk.SetStyle("display", "none")

				return false
			},
		),
	)

	self.edit_b.Set(
		"onclick",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				self.view_div.SetStyle("display", "none")
				self.edit_div.SetStyle("display", "inline")
				return false
			},
		),
	)

	self.cancel_b.Set(
		"onclick",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				self.view_div.SetStyle("display", "inline")
				self.edit_div.SetStyle("display", "none")
				return false
			},
		),
	)

	self.delete_b.Set(
		"onclick",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				delete(self.extension.config.ProxyTargets, self.setting_name)
				go self_Element_Mutator.RemoveFromParent()
				go self.extension.Changed()
				return false
			},
		),
	)

	loadTarget()

	for _, i := range [][2]*elementtreeconstructor.ElementMutator{
		[2]*elementtreeconstructor.ElementMutator{
			self.host_text_cb,
			self.host_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.port_text_cb,
			self.port_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.username_text_cb,
			self.username_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.password_text_cb,
			self.password_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.proxydns_cb_cb,
			self.proxydns_cb,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.failover_timeout_text_cb,
			self.failover_timeout_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.proxy_authorization_header_text_cb,
			self.proxy_authorization_header_text,
		},
		[2]*elementtreeconstructor.ElementMutator{
			self.connection_isolation_key_text_cb,
			self.connection_isolation_key_text,
		},
	} {
		i[0].Set(
			"onchange",
			func(cb, target *elementtreeconstructor.ElementMutator) js.Func {
				t2 := func() {
					checked := cb.Element.Get("checked").Bool()
					target.Set("disabled", !checked)
					if checked {
						target.SetStyle("display", "inline")
					} else {
						target.SetStyle("display", "none")
					}
				}
				t := func(this js.Value, args []js.Value) interface{} {
					t2()
					go self.Changed()
					return false
				}
				ret := js.FuncOf(t)
				t2()

				return ret
			}(i[0], i[1]),
		)

		i[1].Set(
			"onchange",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				go self.Changed()
				return false
			},
			),
		)

		i[0].
			SetStyle("border", "none").
			Set("title", "Check to Enable")

		if i[1] != self.proxydns_cb {
			applyEditorStyle(i[1])
		}
	}

	return self
}

func (self *ProxyTargetEditor) Changed() {
	self.changed_asterisk.SetStyle("display", "inline")
}

// func (self *ProxyTargetEditor) Save() {
// 	self.save_asterisk.SetStyle("display", "none")
// }
