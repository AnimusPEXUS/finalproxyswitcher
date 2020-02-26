package main

import (
	"sort"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection"
)

type RuleSetWidget struct {
	extension *ProxySwitcherExtension

	etc *elementtreeconstructor.ElementTreeConstructor

	domain_settings_editors map[string]*DomainSettingsEditor

	// root     *elementtreeconstructor.ElementMutator
	// controls *elementtreeconstructor.ElementMutator
	editors *elementtreeconstructor.ElementMutator

	Element *elementtreeconstructor.ElementMutator
}

func NewRuleSetWidget(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
) *RuleSetWidget {
	self := &RuleSetWidget{
		etc:                     etc,
		extension:               extension,
		domain_settings_editors: make(map[string]*DomainSettingsEditor),
	}

	root := etc.CreateElement("div")
	controls := etc.CreateElement("div")
	self.editors = etc.CreateElement("div").
		SetStyle("display", "grid").
		SetStyle("gap", "3px")

	root.AppendChildren(
		controls.Element,
		self.editors.Element,
	)

	self.Element = root

	add_button := widgetcollection.NewActiveLabel00(
		"Add",
		nil,
		etc,
		func() {

			if _, ok := self.domain_settings_editors[""]; ok {
				return
			}

			ed := NewDomainSettingsEditor(
				self.etc,
				self.extension,
				nil,
				// self.OnSubEditorChanged,
				self.OnSubEditorDelete,
				self.OnSubEditorRename,
				self.OnSubEditorApply,
			)

			self.addEditor(ed)

			// adding unnamed editor should not mean need to save
			// self.extension.Changed()

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

	controls.AppendChildren(
		etc.CreateTextNode("Domain Rules"),
		etc.CreateTextNode(" "),
		add_button.Element,
		etc.CreateTextNode(" "),
		reload_button.Element,
	)

	self.Reload()

	return self
}

func (self *RuleSetWidget) addEditor(ed *DomainSettingsEditor) {
	self.editors.AppendChildren(ed.Element)
	// TODO: avoid adding if domain == "" ?
	self.domain_settings_editors[ed.DomainSettings.Domain] = ed
}

func (self *RuleSetWidget) rmEditor(ed *DomainSettingsEditor) {
	delete(self.domain_settings_editors, ed.DomainSettings.Domain)
	ed.Element.RemoveFromParent()
}

func (self *RuleSetWidget) Reload() {

	for _, v := range self.domain_settings_editors {
		self.rmEditor(v)
	}

	keys := make([]string, 0)
	for k, _ := range self.extension.config.RuleSet {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		ed := NewDomainSettingsEditor(
			self.etc,
			self.extension,
			self.extension.config.RuleSet[k].Copy(),
			// self.OnSubEditorChanged,
			self.OnSubEditorDelete,
			self.OnSubEditorRename,
			self.OnSubEditorApply,
		)

		self.addEditor(ed)
	}

	return
}

func (self *RuleSetWidget) SubEditorDelete(domain string) {
	if t, ok := self.domain_settings_editors[domain]; ok {
		self.rmEditor(t)
	}

	if _, ok := self.extension.config.RuleSet[domain]; ok {
		delete(self.extension.config.RuleSet, domain)
	}

	self.extension.Changed()
}

func (self *RuleSetWidget) SubEditorRename(old_name, new_name string) {

	if _, ok := self.domain_settings_editors[new_name]; ok {
		// TODO: Show Message. ask confirmation
		self.SubEditorDelete(new_name)
	}

	if t, ok := self.domain_settings_editors[old_name]; ok {
		delete(self.domain_settings_editors, old_name)
		self.domain_settings_editors[new_name] = t

	}

	if t, ok := self.extension.config.RuleSet[old_name]; ok {
		delete(self.extension.config.RuleSet, old_name)
		self.extension.config.RuleSet[new_name] = t
	}

	self.extension.Changed()
}

// func (self *RuleSetWidget) OnSubEditorChanged() {
// 	self.extension.Changed()
// }

func (self *RuleSetWidget) OnSubEditorDelete(domain string) {
	self.SubEditorDelete(domain)
}

func (self *RuleSetWidget) OnSubEditorRename(old_name, new_name string) {
	self.SubEditorRename(old_name, new_name)
}

func (self *RuleSetWidget) OnSubEditorApply(domain string) {
	self.extension.config.RuleSet[domain] =
		self.domain_settings_editors[domain].DomainSettings
	self.extension.Changed()
}
