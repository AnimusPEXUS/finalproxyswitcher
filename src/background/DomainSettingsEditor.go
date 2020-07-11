package main

import (
	"sort"
	"syscall/js"

	"github.com/AnimusPEXUS/gojstools/elementtreeconstructor"
	"github.com/AnimusPEXUS/gojstools/widgetcollection"
	"github.com/AnimusPEXUS/utils/domainname"
)

type DomainSettingsEditor struct {
	etc       *elementtreeconstructor.ElementTreeConstructor
	extension *ProxySwitcherExtension

	root_mode      bool
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
	settings *DomainSettings,
	root_mode bool,
	extension *ProxySwitcherExtension,
	etc *elementtreeconstructor.ElementTreeConstructor,
	// onchange func(),
	ondelete func(domain string),
	onrename func(domain0, domain1 string),
	onapply func(domain string),
) *DomainSettingsEditor {

	self := &DomainSettingsEditor{
		root_mode:      root_mode,
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

	if self.DomainSettings.Domain == nil {
		self.DomainSettings.Domain = domainname.NewDomainNameFromString("")
	}

	if self.domain_settings_subrequests_editors == nil {
		self.domain_settings_subrequests_editors =
			map[string]*DomainSubrequestSettingsEditor{}
	}

	self.domain_input = etc.CreateElement("input").
		Set("type", "text").
		Set("value", self.DomainSettings.Domain.String()).
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
			self.DomainSettings.RulesAndInheritance.Copy(),
			root_mode,
			self.extension,
			etc,
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
			self.DomainSettings.DomainSubrequestSettingsDefaults.Copy(),
			false,
			self.extension,
			etc,
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
				nil,
				self.extension,
				etc,
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
			self.ondelete(self.DomainSettings.Domain.String())
		},
	)

	apply_settings_btn := widgetcollection.NewActiveLabel00(
		"Apply",
		nil,
		etc,
		func() {
			self.ApplySettings()
		},
	)

	title := "This Domain Settings"
	if root_mode {
		title = "Root Settings"
	}

	self.Element = etc.CreateElement("div")

	domain_input_element := (*elementtreeconstructor.ElementMutator)(nil)

	self.Element.
		AppendChildren(
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AssignSelf(&domain_input_element).
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
					etc.CreateTextNode(title),
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

	if root_mode {
		domain_input_element.SetStyle("display", "none")
	}

	self.Reload()

	return self
}

func (self *DomainSettingsEditor) addEditor(ed *DomainSubrequestSettingsEditor) {
	self.editors.AppendChildren(ed.Element)
	// TODO: avoid adding if domain == "" ?
	self.domain_settings_subrequests_editors[ed.DomainSubrequestSettings.Domain.String()] = ed
}

func (self *DomainSettingsEditor) rmEditor(ed *DomainSubrequestSettingsEditor) {
	delete(self.domain_settings_subrequests_editors, ed.DomainSubrequestSettings.Domain.String())
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
			self.DomainSettings.DomainSubrequestSettings[k].Copy(),
			self.extension,
			self.etc,
			// self.OnSubEditorChanged,
			self.OnSubEditorDelete,
			self.OnSubEditorRename,
			self.OnSubEditorApply,
		)

		self.addEditor(ed)
	}

	return
}

func (self *DomainSettingsEditor) ApplySettings() {
	old_name := self.DomainSettings.Domain.String()
	new_name := self.domain_input.GetJsValue("value").String()

	self.onapply(old_name)

	if old_name != new_name {
		self.onrename(old_name, new_name)
		self.DomainSettings.Domain.SetFromString(new_name)
	}

	self.Unchanged()
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
	if self.root_mode {
		self.ApplySettings()
	}
}

func (self *DomainSettingsEditor) Unchanged() {
	self.changed_asterisk.SetStyle("display", "none")
}
