package main

import (
	"log"
	"sort"
	"syscall/js"
	"time"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	pexu_promise "github.com/AnimusPEXUS/wasmtools/promise"
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

	upload_config_input *elementtreeconstructor.ElementMutator

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

	pager_element := etc.CreateElement("div")

	pager_settings := &widgetcollection.Pager00Settings{
		Pages: []*widgetcollection.Pager00Page{
			&widgetcollection.Pager00Page{
				PageId: 0,
				Element: etc.CreateElement("div").
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
			},
			&widgetcollection.Pager00Page{
				PageId:  1,
				Element: self.root_rules_editor.Element,
			},
			&widgetcollection.Pager00Page{
				PageId:  2,
				Element: rule_set_widget.Element,
			},
		},
		DisplayElement: pager_element,
	}

	pager := widgetcollection.NewPager00(etc, pager_settings)

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
												"Export Saved Config..",
												nil,
												etc,
												func() {
													// this goroutine is requred, else this will deadlock
													go func() {
														res, err := extension.GenerateSavedConfigJSON(true)
														if err != nil {
															log.Println("error", err)
															return
														}

														// TODO: move to separate package
														blob := js.Global().Get("Blob").New([]interface{}{res})
														objurl := js.Global().Get("URL").Call("createObjectURL", blob)
														// TODO: use downloads.onChanged to use revokeObjectURL
														// defer func() {
														// 	js.Global().Get("URL").Call("revokeObjectURL", objurl)
														// }()
														js.Global().
															Get("browser").
															Get("downloads").
															Call(
																"download",
																map[string]interface{}{
																	"url":      objurl,
																	"saveAs":   true,
																	"filename": "settings.json",
																},
															)
													}()
												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Export Active Config..",
												nil,
												etc,
												func() {
													res, err := extension.GenerateActiveConfigJSON(true)
													if err != nil {
														log.Println("error", err)
														return
													}

													// TODO: move to separate package
													blob := js.Global().Get("Blob").New([]interface{}{res})
													objurl := js.Global().Get("URL").Call("createObjectURL", blob)
													// TODO: use downloads.onChanged to use revokeObjectURL
													// defer func() {
													// 	js.Global().Get("URL").Call("revokeObjectURL", objurl)
													// }()
													js.Global().
														Get("browser").
														Get("downloads").
														Call(
															"download",
															map[string]interface{}{
																"url":      objurl,
																"saveAs":   true,
																"filename": "settings.json",
															},
														)
												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Import Active Config..",
												&[]string{"press Me, press Me hard (I love this)"}[0],
												etc,
												func() {
													i := etc.CreateElement("input").
														SetAttribute("type", "file")
													i.Set(
														"onchange",
														js.FuncOf(
															func(
																this js.Value,
																args []js.Value,
															) interface{} {
																go func() {
																	log.Println("file changed")

																	files := i.GetJsValue("files")

																	if files.Length() != 0 {
																		file := files.Index(0)

																		file_text_promise := file.Call("text")

																		file_text_promise_go, err := pexu_promise.NewPromiseFromJSValue(file_text_promise)
																		if err != nil {
																			return
																		}

																		file_text := ""

																		psucc := make(chan bool)
																		perr := make(chan bool)

																		file_text_promise_go.Then(
																			js.FuncOf(func(
																				this js.Value,
																				args []js.Value,
																			) interface{} {
																				file_text = args[0].String()
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
																			log.Println("file text received")
																		case <-perr:
																			log.Println("file text receive error")
																			return
																		case <-time.After(time.Duration(time.Minute)):
																			log.Println("file text receive timeout")
																			return
																		}

																		// TODO: error handeling
																		self.extension.UseActiveConfigJSON(file_text)

																	}
																	return
																}()
																return false
															},
														),
													)
													i.Call("click", nil)

												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Page 1",
												nil,
												etc,
												func() {
													pager.SwitchPage(0)
												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Page 2",
												nil,
												etc,
												func() {
													pager.SwitchPage(1)
												},
											).Element,

											etc.CreateTextNode("●"),

											widgetcollection.NewActiveLabel00(
												"Page 3",
												nil,
												etc,
												func() {
													pager.SwitchPage(2)
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
										Set("col-span", "3").
										SetStyle("border", "1px black solid").
										SetStyle("position", "relative").
										SetStyle("vertical-align", "top").
										// SetStyle("height", "100%").
										AppendChildren(pager.Element),
								),
						),
				),
		)

	self.ReloadProxyTargetList()

	return self
}

func (self *MainWindow) ReloadProxyTargetList() {
	self.proxy_targets_div.
		RemoveAllChildren()
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
