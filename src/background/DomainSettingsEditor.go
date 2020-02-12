package main

import (
	"syscall/js"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
)

type DomainSettingsEditor struct {
	value_select *elementtreeconstructor.ElementMutator

	rules_and_inheritance_editor *RulesAndInheritanceEditor

	as_a_subrequest_defaults_editor    *RulesAndInheritanceEditor
	as_a_subrequest_per_domain_editors map[string]*DomainSettingsSubrequestEditor

	domain_input     *elementtreeconstructor.ElementMutator
	changed_asterisk *elementtreeconstructor.ElementMutator

	DomainSettings *DomainSettings

	Element *pexu_dom.Element

	onchange func()
	ondelete func(domain string)
	onrename func(domain0, domain1 string)
	onapply  func(domain string)
}

func NewDomainSettingsEditor(
	domain string,
	document *pexu_dom.Document,
	settings *DomainSettings,
	onchange func(),
	ondelete func(domain string),
	onrename func(domain0, domain1 string),
	onapply func(domain string),
) *DomainSettingsEditor {

	if settings == nil {
		settings = &DomainSettings{}
	}

	self := &DomainSettingsEditor{
		DomainSettings: settings,
		// parent:         parent,

		onchange: onchange,
		ondelete: ondelete,
		onrename: onrename,
		onapply:  onapply,
	}

	if self.as_a_subrequest_per_domain_editors == nil {
		self.as_a_subrequest_per_domain_editors = map[string]*DomainSettingsSubrequestEditor{}
	}

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	self.domain_input = etc.CreateElement("input").
		Set("type", "text").
		Set("value", domain)

	apply_to_subdomains_cb := etc.CreateElement("input").
		Set("type", "checkbox")

	{
		rai := (*RulesAndInheritance)(nil)
		if settings != nil {
			rai = settings.RulesAndInheritance
		}

		self.rules_and_inheritance_editor = NewRulesAndInheritanceEditor(
			document,
			rai,
			func() {},
		)

		aasd := (*RulesAndInheritance)(nil)
		if settings != nil {
			aasd = settings.AsASubrequestDefaults
		}

		self.as_a_subrequest_defaults_editor = NewRulesAndInheritanceEditor(
			document,
			aasd,
			func() {},
		)
	}

	subeditors_div := etc.CreateElement("div").
		SetStyle("padding-left", "6px").
		SetStyle("display", "grid").
		SetStyle("gap", "1px")

	add_subrequest := etc.CreateElement("a").
		ExternalUse(applyAStyle).
		Set(
			"onclick",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					e := NewDomainSettingsSubrequestEditor(
						"",
						document,
						nil,
						self.OnSubEditorChange,
						self.OnSubEditorDelete,
						self.OnSubEditorRename,
					)
					self.as_a_subrequest_per_domain_editors[""] = e
					subeditors_div.AppendChildren(e.Element)
					return false
				},
			),
		).
		AppendChildren(
			etc.CreateTextNode("Add"),
		)

	rename_btn := etc.CreateElement("a").
		ExternalUse(applyAStyle).
		Set(
			"onclick",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					// TODO: add checks
					current := self.DomainSettings.Domain
					new_one := self.domain_input.GetJsValue("value").String()
					onrename(current, new_one)
					self.DomainSettings.Domain = new_one
					self.Unchanged()
					return false
				},
			),
		).
		AppendChildren(
			etc.CreateTextNode("Apply Renaming"),
			etc.CreateElement("span").
				AppendChildren(
					etc.CreateTextNode("*"),
				).
				ExternalUse(applySpanChangedAsterisk).
				AssignSelf(&self.changed_asterisk),
		)

	remove_btn := etc.CreateElement("a").
		ExternalUse(applyAStyle).
		Set(
			"onclick",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					ondelete(domain)
					return false
				},
			),
		).
		AppendChildren(
			etc.CreateTextNode("Remove"),
		)

	apply_settings_btn := etc.CreateElement("a").
		ExternalUse(applyAStyle).
		Set(
			"onclick",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					onapply(self.DomainSettings.Domain)
					return false
				},
			),
		).
		AppendChildren(
			etc.CreateTextNode("Apply"),
		)

	div := etc.CreateElement("div").
		AppendChildren(
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					self.domain_input,
					rename_btn,
					etc.CreateTextNode(" "),
					remove_btn,
					etc.CreateTextNode(" "),
					apply_settings_btn,
				),
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					apply_to_subdomains_cb,
				),
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					etc.CreateTextNode("This Domain Settings"),
					self.rules_and_inheritance_editor.Element,
				),
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					etc.CreateTextNode("Default Settings for Subrequests"),
					self.as_a_subrequest_defaults_editor.Element,
				),
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					etc.CreateTextNode("Subrequest Domain Setting "),
					etc.CreateTextNode(" "),
					add_subrequest,
				),
			subeditors_div,
		).
		SetStyle("border", "1px black solid").
		SetStyle("border-left", "3px magenta solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		SetStyle("display", "grid").
		SetStyle("gap", "3px")

	self.Element = div.Element

	if settings != nil && settings.AsASubrequestPerDomain != nil {
		for k, v := range settings.AsASubrequestPerDomain {
			e := NewDomainSettingsSubrequestEditor(
				k,
				document,
				v,
				self.OnSubEditorChange,
				self.OnSubEditorDelete,
				self.OnSubEditorRename,
			)
			self.as_a_subrequest_per_domain_editors[""] = e
			subeditors_div.AppendChildren(e.Element)
		}
	}

	return self
}

func (self *DomainSettingsEditor) OnSubEditorChange() {
	self.Changed()
}

func (self *DomainSettingsEditor) OnSubEditorDelete(domain string) {

	defer func() { self.Changed() }()

	if t, ok := self.as_a_subrequest_per_domain_editors[domain]; ok {
		delete(self.as_a_subrequest_per_domain_editors, domain)
		elementtreeconstructor.NewElementMutatorFromElement(t.Element).RemoveFromParent()
	}

	if _, ok := self.DomainSettings.AsASubrequestPerDomain[domain]; ok {
		delete(self.DomainSettings.AsASubrequestPerDomain, domain)
	}

}

func (self *DomainSettingsEditor) OnSubEditorRename(domain0, domain1 string) {

	defer func() { self.Changed() }()

	if _, ok := self.as_a_subrequest_per_domain_editors[domain1]; ok {
		// TODO: Show Message. ask confirmation
		self.OnSubEditorDelete(domain1)
	}

	if t, ok := self.as_a_subrequest_per_domain_editors[domain0]; ok {
		delete(self.as_a_subrequest_per_domain_editors, domain0)
		self.as_a_subrequest_per_domain_editors[domain1] = t
	}

	if t, ok := self.DomainSettings.AsASubrequestPerDomain[domain0]; ok {
		delete(self.DomainSettings.AsASubrequestPerDomain, domain0)
		self.DomainSettings.AsASubrequestPerDomain[domain1] = t
	}

}

func (self *DomainSettingsEditor) Changed() {
	self.onchange()
	self.domain_input.SetStyle("display", "inline")
}

func (self *DomainSettingsEditor) Unchanged() {
	self.domain_input.SetStyle("display", "none")
}
