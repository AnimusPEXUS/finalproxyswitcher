package main

import (
	"strconv"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection/select00"
)

type RulesInheritanceEditor struct {
	document  *pexu_dom.Document
	extension *ProxySwitcherExtension

	http_rule_inheritance_select    *select00.Select00
	request_rule_inheritance_select *select00.Select00
	proxy_rule_inheritance_select   *select00.Select00

	RulesInheritance *RulesInheritance

	Element *pexu_dom.Element
}

func NewRuleInheritanceEditor(
	document *pexu_dom.Document,
	preset_rules_inheritance *RulesInheritance,
	onchange func(),
) *RulesInheritanceEditor {

	self := &RulesInheritanceEditor{}

	if preset_rules_inheritance == nil {
		preset_rules_inheritance = &RulesInheritance{}
	}

	self.RulesInheritance = preset_rules_inheritance

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	int_onchange := func() {
		onchange()
	}

	self.http_rule_inheritance_select = Select00FromMapUIntString(
		document,
		RuleInheritanceStrings,
		preset_rules_inheritance.HttpRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.http_rule_inheritance_select.Value)
			self.RulesInheritance.HttpRuleInheritance = uint(i)
			int_onchange()
		},
	)

	self.request_rule_inheritance_select = Select00FromMapUIntString(
		document,
		RuleInheritanceStrings,
		preset_rules_inheritance.RequestRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.request_rule_inheritance_select.Value)
			self.RulesInheritance.RequestRuleInheritance = uint(i)
			int_onchange()
		},
	)

	self.proxy_rule_inheritance_select = Select00FromMapUIntString(
		document,
		RuleInheritanceStrings,
		preset_rules_inheritance.ProxyRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.proxy_rule_inheritance_select.Value)
			self.RulesInheritance.ProxyRuleInheritance = uint(i)
			int_onchange()
		},
	)

	t := etc.CreateElement("div").
		SetStyle("border", "1px black dotted").
		SetStyle("border-left", "3px gold solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		AppendChildren(
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("http rule inheritance"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.http_rule_inheritance_select.Element,
						),
				),
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("request rule inheritance"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.request_rule_inheritance_select.Element,
						),
				),
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("proxy rule inheritance"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.proxy_rule_inheritance_select.Element,
						),
				),
		)

	self.Element = t.Element

	return self
}
