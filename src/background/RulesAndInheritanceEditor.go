package main

import (
	"syscall/js"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
	pexu_wsm_misc "github.com/AnimusPEXUS/wasmtools/misc"
)

type RulesAndInheritanceEditor struct {
	value_select *elementtreeconstructor.ElementMutator

	rules_inheritance_editor *RulesInheritanceEditor
	rules_editor             *RulesEditor

	RulesAndInheritance *RulesAndInheritance

	Element *pexu_dom.Element
}

func NewRulesAndInheritanceEditor(
	document *pexu_dom.Document,
	preset_rules_and_inheritance *RulesAndInheritance,
	onchange func(),
) *RulesAndInheritanceEditor {

	if preset_rules_and_inheritance == nil {
		preset_rules_and_inheritance = &RulesAndInheritance{}
	}

	self := &RulesAndInheritanceEditor{}

	self.RulesAndInheritance = preset_rules_and_inheritance

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	int_onchange := func() {
		onchange()
	}

	apply_to_subdomains_cb := etc.CreateElement("input").
		Set("type", "checkbox").
		Set("value", self.RulesAndInheritance.ApplyToSubdomains)

	apply_to_subdomains_cb.Set(
		"onchange",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				self.RulesAndInheritance.ApplyToSubdomains =
					pexu_wsm_misc.StrBoolToBool(
						apply_to_subdomains_cb.GetJsValue("value").String(),
					)
				int_onchange()
				return false
			},
		),
	)

	self.rules_inheritance_editor = NewRuleInheritanceEditor(
		document,
		preset_rules_and_inheritance.RulesInheritance,
		func() {
			// todo: maybe it's better not to do such assignment
			self.RulesAndInheritance.RulesInheritance = self.rules_inheritance_editor.RulesInheritance
			int_onchange()
		},
	)

	self.rules_editor = NewRulesEditor(
		document,
		preset_rules_and_inheritance.Rules,
		func() {
			self.RulesAndInheritance.Rules = self.rules_editor.Rules
			int_onchange()
		},
	)

	t := etc.CreateElement("div").
		SetStyle("border", "1px black dotted").
		SetStyle("border-left", "3px teal solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		SetStyle("display", "grid").
		SetStyle("gap", "3px").
		AppendChildren(
			etc.CreateElement("div").
				ExternalUse(applyBlackRoundedBoxInRuleEditor).
				AppendChildren(
					apply_to_subdomains_cb.Element,
					etc.CreateTextNode("Apply To Subdomains"),
				),
			self.rules_inheritance_editor.Element,
			self.rules_editor.Element,
		)

	self.Element = t.Element

	return self
}
