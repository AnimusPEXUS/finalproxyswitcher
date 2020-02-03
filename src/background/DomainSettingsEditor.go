package main

import (
	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
)

type DomainSettingsEditor struct {
	value_select *elementtreeconstructor.ElementMutator

	domain_input           *elementtreeconstructor.ElementMutator
	apply_to_subdomains_cb *elementtreeconstructor.ElementMutator

	rules_and_inheritance *RulesAndInheritanceEditor

	as_a_subrequest_defaults   *RulesAndInheritanceEditor
	as_a_subrequest_per_domain map[string]*RulesAndInheritanceEditor

	Element *pexu_dom.Element
}

func NewDomainSettingsEditor(
	document *pexu_dom.Document,
	domain_settings *DomainSettings,
	onchange func(),
) *DomainSettingsEditor {

	self := &DomainSettingsEditor{}

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	form := etc.CreateElement("form").
		AppendChildren(
			self.domain_input.Element,
			self.apply_to_subdomains_cb.Element,
			self.rules_and_inheritance.Element,
			self.as_a_subrequest_defaults.Element,
		)

	for _, i := range self.as_a_subrequest_per_domain {
		form.AppendChildren(i)
	}

	self.Element = form.Element

	return self
}
