package main

import (
	"sort"
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection"
)

type DomainSettingsEditor struct {
	etc       *elementtreeconstructor.ElementTreeConstructor
	extension *ProxySwitcherExtension

	DomainSettings *DomainSettings
	Element        *elementtreeconstructor.ElementMutator

	value_select *elementtreeconstructor.ElementMutator

	rules_and_inheritance_editor *RulesAndInheritanceEditor

	domain_settings_subrequest_defaults_editor *RulesAndInheritanceEditor
	domain_settings_subrequests_editors        map[string]*DomainSubrequestSettingsEditor

	domain_input     *elementtreeconstructor.ElementMutator
	changed_asterisk *elementtreeconstructor.ElementMutator

	editors *elementtreeconstructor.ElementMutator

	// onchange func()
	ondelete func(domain string)
	onrename func(domain0, domain1 string)
	onapply  func(domain string)
}

func NewDomainSettingsEditor(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
	settings *DomainSettings,
	// onchange func(),
	ondelete func(domain string),
	onrename func(domain0, domain1 string),
	onapply func(domain string),
) *DomainSettingsEditor {

	self := &DomainSettingsEditor{
		etc:            etc,
		extension:      extension,
		DomainSettings: settings,
		ondelete:       ondelete,
		onrename:       onrename,
		onapply:        onapply,
	}

	if self.DomainSettings == nil {
		self.DomainSettings = &DomainSettings{}
	}

	if self.DomainSettings.DomainSubrequestSettings == nil {
		self.DomainSettings.DomainSubrequestSettings = map[string]*DomainSubrequestSettings{}
	}

	if self.domain_settings_subrequests_editors == nil {
		self.domain_settings_subrequests_editors =
			map[string]*DomainSubrequestSettingsEditor{}
	}

	self.domain_input = etc.CreateElement("input").
		Set("type", "text").
		Set("value", self.DomainSettings.Domain).
		Set(
			"onchange",
			js.FuncOf(
				func(this js.Value, args []js.Value) interface{} {
					self.Changed()
					return false
				},
			),
		)

	{
		// rai := (*RulesAndInheritance)(nil)
		// if self.DomainSettings != nil && self.DomainSettings.RulesAndInheritance != nil {
		// 	rai = self.DomainSettings.RulesAndInheritance.Copy()
		// }

		self.rules_and_inheritance_editor = NewRulesAndInheritanceEditor(
			etc,
			self.extension,
			self.DomainSettings.RulesAndInheritance.Copy(),
			func() {
				self.DomainSettings.RulesAndInheritance =
					self.rules_and_inheritance_editor.RulesAndInheritance
				self.Changed()
			},
		)

		// aasd := (*RulesAndInheritance)(nil)
		// if self.DomainSettings != nil &&
		// 	self.DomainSettings.DomainSubrequestSettingsDefaults != nil {
		// 	aasd = self.DomainSettings.DomainSubrequestSettingsDefaults.Copy()
		// }

		self.domain_settings_subrequest_defaults_editor = NewRulesAndInheritanceEditor(
			etc,
			self.extension,
			self.DomainSettings.DomainSubrequestSettingsDefaults.Copy(),
			func() {
				self.DomainSettings.DomainSubrequestSettingsDefaults =
					self.domain_settings_subrequest_defaults_editor.RulesAndInheritance
				self.Changed()
			},
		)
	}

	self.editors = etc.CreateElement("div").
		SetStyle("padding-left", "6px").
		SetStyle("display", "grid").
		SetStyle("gap", "1px")

	add_subrequest := widgetcollection.NewActiveLabel00(
		"Add",
		nil,
		etc,
		func() {
			if _, ok := self.domain_settings_subrequests_editors[""]; ok {
				return
			}

			e := NewDomainSubrequestSettingsEditor(
				etc,
				self.extension,
				nil,
				// self.OnSubEditorChanged,
				self.OnSubEditorDelete,
				self.OnSubEditorRename,
				self.OnSubEditorApply,
			)

			self.addEditor(e)
		},
	)

	reload_button := widgetcollection.NewActiveLabel00(
		"Reload",
		nil,
		etc,
		func() {
			self.Reload()
		},
	)

	remove_btn := widgetcollection.NewActiveLabel00(
		"Remove",
		nil,
		etc,
		func() {
			self.ondelete(self.DomainSettings.Domain)
		},
	)

	apply_settings_btn := widgetcollection.NewActiveLabel00(
		"Apply",
		nil,
		etc,
		func() {

			old_name := self.DomainSettings.Domain
			new_name := self.domain_input.GetJsValue("value").String()

			self.onapply(old_name)

			if old_name != new_name {
				self.onrename(old_name, new_name)
				self.DomainSettings.Domain = new_name
			}

			self.Unchanged()
		},
	)

	self.Element = etc.CreateElement("div").
		AppendChildren(
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					self.domain_input,
					remove_btn.Element,
					etc.CreateTextNode(" "),
					apply_settings_btn.Element,
					etc.CreateElement("span").
						AppendChildren(
							etc.CreateTextNode("*"),
						).
						ExternalUse(applySpanChangedAsterisk).
						AssignSelf(&self.changed_asterisk),
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
					self.domain_settings_subrequest_defaults_editor.Element,
				),
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					etc.CreateTextNode("Subrequest Domain Setting "),
					etc.CreateTextNode(" "),
					add_subrequest.Element,
					etc.CreateTextNode(" "),
					reload_button.Element,
				),
			self.editors,
		).
		SetStyle("border", "1px black solid").
		SetStyle("border-left", "3px magenta solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		SetStyle("display", "grid").
		SetStyle("gap", "3px")

	self.Reload()

	return self
}

func (self *DomainSettingsEditor) addEditor(ed *DomainSubrequestSettingsEditor) {
	self.editors.AppendChildren(ed.Element)
	// TODO: avoid adding if domain == "" ?
	self.domain_settings_subrequests_editors[ed.DomainSubrequestSettings.Domain] = ed
}

func (self *DomainSettingsEditor) rmEditor(ed *DomainSubrequestSettingsEditor) {
	delete(self.domain_settings_subrequests_editors, ed.DomainSubrequestSettings.Domain)
	ed.Element.RemoveFromParent()
}

func (self *DomainSettingsEditor) Reload() {

	for _, v := range self.domain_settings_subrequests_editors {
		self.rmEditor(v)
	}

	keys := make([]string, 0)
	for k, _ := range self.DomainSettings.DomainSubrequestSettings {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		ed := NewDomainSubrequestSettingsEditor(
			self.etc,
			self.extension,
			self.DomainSettings.DomainSubrequestSettings[k].Copy(),
			// self.OnSubEditorChanged,
			self.OnSubEditorDelete,
			self.OnSubEditorRename,
			self.OnSubEditorApply,
		)

		self.addEditor(ed)
	}

	return
}

func (self *DomainSettingsEditor) SubEditorDelete(domain string) {

	if t, ok := self.domain_settings_subrequests_editors[domain]; ok {
		self.rmEditor(t)
	}

	if _, ok := self.DomainSettings.DomainSubrequestSettings[domain]; ok {
		delete(self.DomainSettings.DomainSubrequestSettings, domain)
	}

	// self.onchange()

}

func (self *DomainSettingsEditor) SubEditorRename(old_name, new_name string) {

	if _, ok := self.domain_settings_subrequests_editors[new_name]; ok {
		// TODO: Show Message. ask confirmation
		self.SubEditorDelete(new_name)
	}

	if t, ok := self.domain_settings_subrequests_editors[old_name]; ok {
		delete(self.domain_settings_subrequests_editors, old_name)
		self.domain_settings_subrequests_editors[new_name] = t
	}

	if t, ok := self.DomainSettings.DomainSubrequestSettings[old_name]; ok {
		delete(self.DomainSettings.DomainSubrequestSettings, old_name)
		self.DomainSettings.DomainSubrequestSettings[new_name] = t
	}

	self.Changed()

}

// func (self *DomainSettingsEditor) OnSubEditorChanged() {
// 	self.Changed()
// }

func (self *DomainSettingsEditor) OnSubEditorDelete(domain string) {
	self.SubEditorDelete(domain)
}

func (self *DomainSettingsEditor) OnSubEditorRename(old_name, new_name string) {
	self.SubEditorRename(old_name, new_name)
}

func (self *DomainSettingsEditor) OnSubEditorApply(domain string) {
	self.DomainSettings.DomainSubrequestSettings[domain] =
		self.domain_settings_subrequests_editors[domain].DomainSubrequestSettings

	self.Changed()
}

func (self *DomainSettingsEditor) Changed() {
	self.changed_asterisk.SetStyle("display", "inline")
}

func (self *DomainSettingsEditor) Unchanged() {
	self.changed_asterisk.SetStyle("display", "none")
}
