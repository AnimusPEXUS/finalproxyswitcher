package main

import (
	"sort"
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

type MainWindow struct {
	extension *ProxySwitcherExtension

	add_new_proxy_target_button     *elementtreeconstructor.ElementMutator
	proxy_targets_div               *elementtreeconstructor.ElementMutator
	reload_proxy_target_list_button *elementtreeconstructor.ElementMutator

	save_settings_button *elementtreeconstructor.ElementMutator
	save_asterisk        *elementtreeconstructor.ElementMutator

	export_saved_settings_button  *elementtreeconstructor.ElementMutator
	export_active_settings_button *elementtreeconstructor.ElementMutator
	import_active_settings_button *elementtreeconstructor.ElementMutator

	root_rules_editor *RulesEditor

	Element *elementtreeconstructor.ElementMutator
}

func (self *ProxySwitcherExtension) MainWindow(
	etc *elementtreeconstructor.ElementTreeConstructor,
) *MainWindow {
	ret := NewMainWindow(etc, self)
	self.main_window = ret
	if self.changed {
		ret.Changed()
	}
	return ret
}

func NewMainWindow(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
) *MainWindow {

	self := &MainWindow{}

	self.extension = extension

	rule_set_widget := NewRuleSetWidget(
		etc,
		extension,
	)

	self.root_rules_editor = NewRulesEditor(
		etc,
		extension,
		self.extension.config.RootRules.Copy(),
		func() {
			self.extension.config.RootRules = self.root_rules_editor.Rules
			self.Changed()
		},
	)

	self.Element = etc.CreateElement("html").
		SetStyle("position", "absolute").
		SetStyle("top", "0px").
		SetStyle("bottom", "0px").
		SetStyle("left", "0px").
		SetStyle("right", "0px").
		SetStyle("margin", "0px").
		SetStyle("padding", "0px").
		AppendChildren(
			etc.CreateElement("head").
				AppendChildren(
					etc.CreateElement("title").
						AppendChildren(
							etc.CreateTextNode("main title text"),
						),
				),
			etc.CreateElement("body").
				SetStyle("position", "absolute").
				SetStyle("top", "0px").
				SetStyle("bottom", "0px").
				SetStyle("left", "0px").
				SetStyle("right", "0px").
				SetStyle("margin", "0px").
				SetStyle("padding", "0px").
				SetStyle("font-size", "10px").
				AppendChildren(
					etc.CreateElement("table").
						SetStyle("table-layout", "fixed").
						SetStyle("position", "absolute").
						// SetStyle("display", "block").
						SetStyle("top", "0px").
						SetStyle("bottom", "0px").
						SetStyle("left", "0px").
						SetStyle("right", "0px").
						AppendChildren(
							etc.CreateElement("tr").
								AppendChildren(
									etc.CreateElement("td").
										AppendChildren(
											etc.CreateElement("a").
												ExternalUse(applyAStyle).
												AssignSelf(&self.save_settings_button).
												AppendChildren(
													etc.CreateTextNode("Save"),
													etc.CreateElement("span").
														ExternalUse(applySpanChangedAsterisk).
														AppendChildren(
															etc.CreateTextNode("*"),
														).
														AssignSelf(&self.save_asterisk),
												),
										),
									etc.CreateElement("td").
										SetAttribute("colspan", "2").
										AppendChildren(
											etc.CreateTextNode("Final Proxy Switcher - by AnimusPEXUS"),
										),
								),
							etc.CreateElement("tr").
								AppendChildren(

									etc.CreateElement("td").
										SetStyle("border", "1px black solid").
										SetStyle("position", "relative").
										SetStyle("vertical-align", "top").
										// SetStyle("height", "100%").
										AppendChildren(
											etc.CreateElement("div").
												SetStyle("overflow-y", "scroll").
												AppendChildren(

													etc.CreateElement("div").
														AppendChildren(

															etc.CreateElement("a").
																Set("title", "add new proxy target").
																ExternalUse(applyAStyle).
																AssignSelf(&self.add_new_proxy_target_button).
																AppendChildren(
																	etc.CreateTextNode("add"),
																),

															etc.CreateTextNode("●"),

															etc.CreateElement("a").
																Set("title", "reload list").
																ExternalUse(applyAStyle).
																AssignSelf(&self.reload_proxy_target_list_button).
																AppendChildren(
																	etc.CreateTextNode("reload"),
																),

															etc.CreateTextNode("●"),

															etc.CreateElement("a").
																ExternalUse(applyAStyle).
																AssignSelf(&self.export_saved_settings_button).
																AppendChildren(
																	etc.CreateTextNode("export saved settings"),
																),

															etc.CreateTextNode("●"),

															etc.CreateElement("a").
																ExternalUse(applyAStyle).
																AssignSelf(&self.export_active_settings_button).
																AppendChildren(
																	etc.CreateTextNode("export active settings"),
																),

															etc.CreateTextNode("●"),

															etc.CreateElement("a").
																ExternalUse(applyAStyle).
																AssignSelf(&self.import_active_settings_button).
																AppendChildren(
																	etc.CreateTextNode("import active settings"),
																),
														),

													etc.CreateElement("div").
														AssignSelf(&self.proxy_targets_div),
												),
										),
									etc.CreateElement("td").
										SetStyle("border", "1px black solid").
										AppendChildren(
											self.root_rules_editor.Element,
										),
									etc.CreateElement("td").
										SetStyle("border", "1px black solid").
										AppendChildren(
											rule_set_widget.Element,
										),
								),
						),
				),
		)

	addNewProxyTarget := func(
		this js.Value,
		args []js.Value,
	) interface{} {
		self.proxy_targets_div.
			AppendChildren(
				self.extension.ProxyTargetEditor(
					"",
					true,
					true,
					etc,
					func() {},
				).Element,
			)
		return true
	}

	reloadProxyTargetList := func() {
		self.proxy_targets_div.
			RemoveChildren()
		lst := []string{}
		for k, _ := range self.extension.config.ProxyTargets {
			lst = append(lst, k)
		}
		sort.Strings(lst)
		for _, i := range lst {
			self.proxy_targets_div.
				AppendChildren(
					self.extension.ProxyTargetEditor(
						i,
						true,
						true,
						etc,
						func() {},
					).Element,
				)
		}
	}

	reloadProxyTargetList_j := func(
		this js.Value,
		args []js.Value,
	) interface{} {
		reloadProxyTargetList()
		return true
	}

	self.add_new_proxy_target_button.Set("onclick", js.FuncOf(addNewProxyTarget))
	self.reload_proxy_target_list_button.Set("onclick", js.FuncOf(reloadProxyTargetList_j))
	self.save_settings_button.Set(
		"onclick",
		js.FuncOf(
			func(
				this js.Value,
				args []js.Value,
			) interface{} {
				go self.Save()
				return true
			},
		),
	)

	go reloadProxyTargetList()

	return self
}

func (self *MainWindow) Changed() {
	self.save_asterisk.SetStyle("display", "inline")
}

func (self *MainWindow) Save() {
	self.save_asterisk.SetStyle("display", "none")
	go self.extension.SaveConfig()
}
