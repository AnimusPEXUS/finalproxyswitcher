package main

import (
	"strconv"

	pexu_dom "github.com/AnimusPEXUS/wasmtools/dom"
	"github.com/AnimusPEXUS/wasmtools/dom/elementtreeconstructor"
	"github.com/AnimusPEXUS/wasmtools/widgetcollection/select00"
)

type RulesEditor struct {
	document *pexu_dom.Document

	http_rule_select    *select00.Select00
	request_rule_select *select00.Select00
	proxy_rule_select   *select00.Select00

	Rules *Rules

	Element *pexu_dom.Element
}

func NewRulesEditor(
	document *pexu_dom.Document,
	preset_rules *Rules,
	onchange func(),
) *RulesEditor {

	self := &RulesEditor{}

	if preset_rules == nil {
		preset_rules = &Rules{}
	}

	self.Rules = preset_rules

	etc := elementtreeconstructor.NewElementTreeConstructor(document)

	int_onchange := func() {
		onchange()
	}

	self.http_rule_select = Select00FromMapUIntString(
		document,
		HttpRuleStrings,
		preset_rules.HttpRule,
		func() {
			i, _ := strconv.Atoi(self.http_rule_select.Value)
			self.Rules.HttpRule = uint(i)
			int_onchange()
		},
	)

	self.request_rule_select = Select00FromMapUIntString(
		document,
		RequestRuleStrings,
		preset_rules.RequestRule,
		func() {
			i, _ := strconv.Atoi(self.request_rule_select.Value)
			self.Rules.RequestRule = uint(i)
			int_onchange()
		},
	)

	self.proxy_rule_select = Select00FromMapUIntString(
		document,
		ProxyRuleString,
		preset_rules.ProxyRule,
		func() {
			i, _ := strconv.Atoi(self.proxy_rule_select.Value)
			self.Rules.ProxyRule = uint(i)
			int_onchange()
		},
	)

	t := etc.CreateElement("table").
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
		)

	self.Element = t.Element

	return self
}
