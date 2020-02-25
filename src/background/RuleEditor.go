package main

import (
	"strconv"

	"github.com/AnimusPEXUS/wasmtools/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection/select00"
)

type RulesEditor struct {
	etc *elementtreeconstructor.ElementTreeConstructor

	http_rule_select    *select00.Select00
	request_rule_select *select00.Select00
	proxy_rule_select   *select00.Select00
	proxy_target_select *select00.Select00

	Rules *Rules

	Element *elementtreeconstructor.ElementMutator
}

func NewRulesEditor(
	etc *elementtreeconstructor.ElementTreeConstructor,
	extension *ProxySwitcherExtension,
	preset_rules *Rules,
	onchange func(),
) *RulesEditor {

	self := &RulesEditor{}

	if preset_rules == nil {
		preset_rules = &Rules{}
	}

	self.Rules = preset_rules
	self.etc = etc

	int_onchange := func() {
		onchange()
	}

	self.http_rule_select = Select00FromMapUIntString(
		etc,
		HttpRuleStrings,
		preset_rules.HttpRule,
		func() {
			i, _ := strconv.Atoi(self.http_rule_select.Value)
			self.Rules.HttpRule = uint(i)
			int_onchange()
		},
	)

	self.request_rule_select = Select00FromMapUIntString(
		etc,
		RequestRuleStrings,
		preset_rules.RequestRule,
		func() {
			i, _ := strconv.Atoi(self.request_rule_select.Value)
			self.Rules.RequestRule = uint(i)
			int_onchange()
		},
	)

	self.proxy_rule_select = Select00FromMapUIntString(
		etc,
		ProxyRuleString,
		preset_rules.ProxyRule,
		func() {
			i, _ := strconv.Atoi(self.proxy_rule_select.Value)
			self.Rules.ProxyRule = uint(i)
			int_onchange()
		},
	)

	self.proxy_target_select = select00.NewSelect00(
		etc,
		extension.ProxyTargetList(),
		preset_rules.ProxyTarget,
		func() {
			self.Rules.ProxyTarget = self.proxy_target_select.Value
			int_onchange()
		},
	)

	self.Element = etc.CreateElement("table").
		SetStyle("border", "1px black dotted").
		SetStyle("border-left", "3px lime solid").
		SetStyle("border-radius", "5px").
		SetStyle("padding", "3px").
		AppendChildren(
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("http rule"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.http_rule_select.Element,
						),
				),
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("request rule"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.request_rule_select.Element,
						),
				),
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("proxy rule"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.proxy_rule_select.Element,
						),
				),
			etc.CreateElement("tr").
				AppendChildren(
					etc.CreateElement("td").
						AppendChildren(
							etc.CreateTextNode("proxy target"),
						),
					etc.CreateElement("td").
						AppendChildren(
							self.proxy_target_select.Element,
						),
				),
		)

	return self
}
