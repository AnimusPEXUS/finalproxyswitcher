package main

// TODO: add cookies

type RuleInheritance = uint

const (
	RuleInheritanceNone RuleInheritance = iota
	RuleInheritanceParent
	RuleInheritanceGlobal
)

var RuleInheritanceStrings = map[RuleInheritance]string{
	RuleInheritanceNone:   "None",
	RuleInheritanceParent: "Parent",
	RuleInheritanceGlobal: "Global",
}

type HttpRule = uint

const (
	HttpRuleUndefined HttpRule = iota
	HttpRuleBlock
	HttpRuleConvertToHttps
	HttpRulePass
)

var HttpRuleStrings = map[HttpRule]string{
	HttpRuleUndefined:      "Undefined",
	HttpRuleBlock:          "Block",
	HttpRuleConvertToHttps: "Convert To Https",
	HttpRulePass:           "Pass",
}

type RequestRule = uint

const (
	RequestRuleUndefined RequestRule = iota
	RequestRuleBlock
	RequestRuleAllow
)

var RequestRuleStrings = map[RequestRule]string{
	RequestRuleUndefined: "Undefined",
	RequestRuleBlock:     "Block",
	RequestRuleAllow:     "Allow",
}

type ProxyRule = uint

const (
	ProxyRuleUndefined ProxyRule = iota
	ProxyRuleUseTarget
	ProxyRulePassUnchanged
)

var ProxyRuleString = map[ProxyRule]string{
	ProxyRuleUndefined:     "Undefined",
	ProxyRuleUseTarget:     "Use Target",
	ProxyRulePassUnchanged: "Pass Unchanged",
}

type (
	RuleSet struct {
		DefaultHighRequestRule *Rules
		DefaultSubRequestRule  *Rules
		HigherRequestRules     map[string]*DomainSettings
	}

	DomainSettings struct {
		Domain string // e.g. org, onion, i2p, com, net ... etc

		ApplyToSubdomains bool

		RulesAndInheritance *RulesAndInheritance

		AsASubrequestDefaults  *RulesAndInheritance
		AsASubrequestPerDomain map[string]*RulesAndInheritance
	}

	RulesAndInheritance struct {
		RulesInheritance *RulesInheritance
		Rules            *Rules
	}

	RulesInheritance struct {
		HttpRuleInheritance    RuleInheritance
		RequestRuleInheritance RuleInheritance
		ProxyRuleInheritance   RuleInheritance
	}

	Rules struct {
		HttpRule    HttpRule
		RequestRule RequestRule
		ProxyRule   ProxyRule
	}
)
