package main

import (
	"sort"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection"
)

type MainWindow struct {
	etc       *elementtreeconstructor.ElementTreeConstructor
	extension *ProxySwitcherExtension

	add_new_proxy_target_button     *elementtreeconstructor.ElementMutator
	proxy_targets_div               *elementtreeconstructor.ElementMutator
	reload_proxy_target_list_button *elementtreeconstructor.ElementMutator

	save_settings_button *elementtreeconstructor.ElementMutator
	save_asterisk        *elementtreeconstructor.ElementMutator

	// export_saved_settings_button  *elementtreeconstructor.ElementMutator
	// export_active_settings_button *elementtreeconstructor.ElementMutator
	// import_active_settings_button *elementtreeconstructor.ElementMutator

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

	self := &MainWindow{etc: etc}

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
											widgetcollection.NewActiveLabel00(
												"Save",
												nil,
												etc,
												func() {
													self.Save()
												},
											).Element,

											etc.CreateElement("span").
												ExternalUse(applySpanChangedAsterisk).
												AppendChildren(
													etc.CreateTextNode("*"),
												).
												AssignSelf(&self.save_asterisk),

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Export Saved Settings",
												nil,
												etc,
												func() {

												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Export Active Settings",
												nil,
												etc,
												func() {

												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Import Active Settings",
												nil,
												etc,
												func() {

												},
											).Element,
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

															widgetcollection.NewActiveLabel00(
																"Add",
																&[]string{"add new proxy target"}[0],
																etc,
																func() {
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
																},
															).Element,

															etc.CreateTextNode("●"),

															widgetcollection.NewActiveLabel00(
																"Reload",
																&[]string{"add new proxy target"}[0],
																etc,
																func() {
																	self.ReloadProxyTargetList()
																},
															).Element,
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

	self.ReloadProxyTargetList()

	return self
}

func (self *MainWindow) ReloadProxyTargetList() {
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
					self.etc,
					func() {},
				).Element,
			)
	}
}

func (self *MainWindow) Changed() {
	self.save_asterisk.SetStyle("display", "inline")
}

func (self *MainWindow) Save() {
	self.save_asterisk.SetStyle("display", "none")
	go self.extension.SaveConfig()
}
