package main

import (
	"syscall/js"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
)

//	pexu_wsm_misc "github.com/AnimusPEXUS/wasmtools/misc"

type RulesAndInheritanceEditor struct {
	value_select *elementtreeconstructor.ElementMutator

	rules_inheritance_editor *RulesInheritanceEditor
	rules_editor             *RulesEditor

	RulesAndInheritance *RulesAndInheritance

	Element *elementtreeconstructor.ElementMutator
}

func NewRulesAndInheritanceEditor(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
	preset_rules_and_inheritance *RulesAndInheritance,
	onchange func(),
) *RulesAndInheritanceEditor {

	if preset_rules_and_inheritance == nil {
		preset_rules_and_inheritance = &RulesAndInheritance{}
	}

	self := &RulesAndInheritanceEditor{}

	self.RulesAndInheritance = preset_rules_and_inheritance

	int_onchange := func() {
		onchange()
	}

	apply_to_subdomains_cb := etc.CreateElement("input").
		Set("type", "checkbox").
		Set("checked", self.RulesAndInheritance.ApplyToSubdomains)

	apply_to_subdomains_cb.Set(
		"onchange",
		js.FuncOf(
			func(this js.Value, args []js.Value) interface{} {
				self.RulesAndInheritance.ApplyToSubdomains =
					apply_to_subdomains_cb.GetJsValue("checked").Bool()
					// pexu_wsm_misc.StrBoolToBool(

					// )

				// log.Println("self.RulesAndInheritance.ApplyToSubdomains as", self.RulesAndInheritance.ApplyToSubdomains)

				int_onchange()
				return false
			},
		),
	)

	self.rules_inheritance_editor = NewRuleInheritanceEditor(
		etc,
		extension,
		self.RulesAndInheritance.RulesInheritance.Copy(),
		func() {
			self.RulesAndInheritance.RulesInheritance = self.rules_inheritance_editor.RulesInheritance
			int_onchange()
		},
	)

	self.rules_editor = NewRulesEditor(
		etc,
		extension,
		self.RulesAndInheritance.Rules.Copy(),
		func() {
			self.RulesAndInheritance.Rules = self.rules_editor.Rules
			int_onchange()
		},
	)

	self.Element = etc.CreateElement("div").
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

	return self
}
