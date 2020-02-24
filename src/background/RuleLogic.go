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

		RulesAndInheritance *RulesAndInheritance

		DomainSubrequestSettingsDefaults *RulesAndInheritance
		DomainSubrequestSettings         map[string]*DomainSubrequestSettings
	}

	DomainSubrequestSettings struct {
		Domain string // e.g. org, onion, i2p, com, net ... etc

		RulesAndInheritance *RulesAndInheritance
	}

	RulesAndInheritance struct {
		ApplyToSubdomains bool
		RulesInheritance  *RulesInheritance
		Rules             *Rules
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
		ProxyTarget string
	}
)

func (self *DomainSettings) Copy() *DomainSettings {

	if self == nil {
		return nil
	}

	var new_self DomainSettings

	new_self = *self
	new_self.RulesAndInheritance = new_self.RulesAndInheritance.Copy()
	new_self.DomainSubrequestSettingsDefaults = new_self.DomainSubrequestSettingsDefaults.Copy()

	old_DomainSubrequestSettings := new_self.DomainSubrequestSettings

	new_self.DomainSubrequestSettings = make(map[string]*DomainSubrequestSettings)

	for k, v := range old_DomainSubrequestSettings {
		new_self.DomainSubrequestSettings[k] = v.Copy()
	}

	return &new_self
}

func (self *DomainSubrequestSettings) Copy() *DomainSubrequestSettings {

	if self == nil {
		return nil
	}

	var new_self DomainSubrequestSettings

	new_self = *self
	new_self.RulesAndInheritance = new_self.RulesAndInheritance.Copy()

	return &new_self
}

func (self *RulesAndInheritance) Copy() *RulesAndInheritance {
	if self == nil {
		return nil
	}

	var new_self RulesAndInheritance

	new_self = *self

	new_self.RulesInheritance = new_self.RulesInheritance.Copy()
	new_self.Rules = new_self.Rules.Copy()

	return &new_self
}

func (self *RulesInheritance) Copy() *RulesInheritance {

	if self == nil {
		return nil
	}

	var new_self RulesInheritance
	new_self = *self
	return &new_self
}

func (self *Rules) Copy() *Rules {

	if self == nil {
		return nil
	}

	var new_self Rules
	new_self = *self
	return &new_self
}
