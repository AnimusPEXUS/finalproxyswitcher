package main

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection"
)

type DomainSubrequestSettingsEditor struct {
	extension *ProxySwitcherExtension

	DomainSubrequestSettings *DomainSubrequestSettings
	Element                  *elementtreeconstructor.ElementMutator

	hRulesAndInheritanceEditor *RulesAndInheritanceEditor

	// parent *DomainSettingsEditor

	domain_input     *elementtreeconstructor.ElementMutator
	changed_asterisk *elementtreeconstructor.ElementMutator

	// onchange func()
	ondelete func(domain string)
	onrename func(domain0, domain1 string)
	onapply  func(domain string)
}

func NewDomainSubrequestSettingsEditor(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
	settings *DomainSubrequestSettings,
	ondelete func(domain string),
	onrename func(old_name, new_name string),
	onapply func(domain string),
) *DomainSubrequestSettingsEditor {

	self := &DomainSubrequestSettingsEditor{
		extension: extension,

		DomainSubrequestSettings: settings,
		ondelete:                 ondelete,
		onrename:                 onrename,
		onapply:                  onapply,
	}

	if self.DomainSubrequestSettings == nil {
		self.DomainSubrequestSettings = &DomainSubrequestSettings{}
	}

	// if self.DomainSubrequestSettings.RulesAndInheritance == nil {
	// 	self.DomainSubrequestSettings.RulesAndInheritance = &RulesAndInheritance{}
	// }

	self.domain_input = etc.CreateElement("input").
		Set("type", "text").
		Set("value", self.DomainSubrequestSettings.Domain).
		Set(
			"onchange",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					self.Changed()
					return false
				},
			),
		)

	self.hRulesAndInheritanceEditor = NewRulesAndInheritanceEditor(
		etc,
		self.extension,
		self.DomainSubrequestSettings.RulesAndInheritance, // TODO: make copy? yes!
		func() {
			self.DomainSubrequestSettings.RulesAndInheritance =
				self.hRulesAndInheritanceEditor.RulesAndInheritance
			self.Changed()
		},
	)

	apply_btn := widgetcollection.NewActiveLabel00(
		"Apply",
		nil,
		etc,
		func() {

			old_name := self.DomainSubrequestSettings.Domain
			new_name := self.domain_input.GetJsValue("value").String()

			self.onapply(old_name)

			if old_name != new_name {
				self.onrename(old_name, new_name)
				self.DomainSubrequestSettings.Domain = new_name
			}

			self.Unchanged()

		},
	)

	remove_btn := widgetcollection.NewActiveLabel00(
		"Remove",
		nil,
		etc,
		func() {
			ondelete(self.DomainSubrequestSettings.Domain)

		},
	)

	self.Element = etc.CreateElement("div").
		SetStyle("border", "1px black solid").
		SetStyle("border-left", "3px blue solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		SetStyle("display", "grid").
		SetStyle("gap", "3px").
		AppendChildren(
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					self.domain_input.Element,
					remove_btn.Element,
					etc.CreateTextNode(" "),
					apply_btn.Element,
					etc.CreateElement("span").
						AppendChildren(
							etc.CreateTextNode("*"),
						).
						ExternalUse(applySpanChangedAsterisk).
						AssignSelf(&self.changed_asterisk),
				),
			self.hRulesAndInheritanceEditor.Element,
		)

	return self
}

func (self *DomainSubrequestSettingsEditor) Changed() {
	// self.onchange()
	self.changed_asterisk.SetStyle("display", "inline")
}

func (self *DomainSubrequestSettingsEditor) Unchanged() {
	self.changed_asterisk.SetStyle("display", "none")
}
