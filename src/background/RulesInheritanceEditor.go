package main

import (
	"strconv"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection/select00"
)

type RulesInheritanceEditor struct {
	etc *elementtreeconstructor.ElementTreeConstructor

	http_rule_inheritance_select    *select00.Select00
	request_rule_inheritance_select *select00.Select00
	proxy_rule_inheritance_select   *select00.Select00

	RulesInheritance *RulesInheritance

	Element *elementtreeconstructor.ElementMutator
}

func NewRuleInheritanceEditor(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
	preset_rules_inheritance *RulesInheritance,
	onchange func(),
) *RulesInheritanceEditor {

	self := &RulesInheritanceEditor{}

	if preset_rules_inheritance == nil {
		preset_rules_inheritance = &RulesInheritance{}
	}

	self.RulesInheritance = preset_rules_inheritance
	self.etc = etc

	int_onchange := func() {
		onchange()
	}

	self.http_rule_inheritance_select = Select00FromMapUIntString(
		etc,
		RuleInheritanceStrings,
		self.RulesInheritance.HttpRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.http_rule_inheritance_select.Value)
			self.RulesInheritance.HttpRuleInheritance = uint(i)
			int_onchange()
		},
	)

	self.request_rule_inheritance_select = Select00FromMapUIntString(
		etc,
		RuleInheritanceStrings,
		self.RulesInheritance.RequestRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.request_rule_inheritance_select.Value)
			self.RulesInheritance.RequestRuleInheritance = uint(i)
			int_onchange()
		},
	)

	self.proxy_rule_inheritance_select = Select00FromMapUIntString(
		etc,
		RuleInheritanceStrings,
		self.RulesInheritance.ProxyRuleInheritance,
		func() {
			i, _ := strconv.Atoi(self.proxy_rule_inheritance_select.Value)
			self.RulesInheritance.ProxyRuleInheritance = uint(i)
			int_onchange()
		},
	)

	self.Element = etc.CreateElement("div").
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

	return self
}
