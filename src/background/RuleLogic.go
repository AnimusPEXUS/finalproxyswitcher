package main

// TODO: add cookies

type RuleInheritance = uint

const (
	RuleInheritanceNone RuleInheritance = iota
	RuleInheritanceParent
	RuleInheritanceGlobal
)

// fixme:
var RuleInheritanceStrings = map[RuleInheritance]string{
	//var RuleInheritanceStrings = map[uint]string{
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

// fixme:
var HttpRuleStrings = map[HttpRule]string{
	// var HttpRuleStrings = map[uint]string{
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

// fixme:
var RequestRuleStrings = map[RequestRule]string{
	// var RequestRuleStrings = map[uint]string{
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

// fixme:
var ProxyRuleString = map[ProxyRule]string{
	// var ProxyRuleString = map[uint]string{
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
