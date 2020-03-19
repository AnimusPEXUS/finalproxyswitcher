package main

import (
	"github.com/AnimusPEXUS/utils/domainname"
)

type (
	ConfigModel struct {
		RootRules    *DomainSettings
		ProxyTargets map[string]*ProxyTarget
		RuleSet      map[string]*DomainSettings
	}

	ProxyTarget struct {
		Type                     string
		Host                     *string
		Port                     *string
		Username                 *string
		Password                 *string
		ProxyDNS                 *bool
		FailoverTimeout          *int
		ProxyAuthorizationHeader *string
		ConnectionIsolationKey   *string
	}
)

func NewConfigModel() *ConfigModel {
	self := &ConfigModel{
		ProxyTargets: map[string]*ProxyTarget{},
		RuleSet:      map[string]*DomainSettings{},
	}
	return self
}

func (self *ConfigModel) Fix() {
	self.FixDomains()
}

func (self *ConfigModel) FixDomains() {

	for k1, _ := range self.RuleSet {
		if self.RuleSet[k1].Domain == nil {
			self.RuleSet[k1].Domain = domainname.NewDomainNameFromString(k1)
		}
		for k2, _ := range self.RuleSet[k1].DomainSubrequestSettings {
			if self.RuleSet[k1].DomainSubrequestSettings[k2].Domain == nil {
				self.RuleSet[k1].DomainSubrequestSettings[k2].Domain = domainname.NewDomainNameFromString(k2)
			}
		}
	}

loo0:
	for k1, _ := range self.RuleSet {
		if k1 != self.RuleSet[k1].Domain.String() {
			if k1 == "" && self.RuleSet[k1].Domain.String() != "" {
				z := self.RuleSet[k1]
				delete(self.RuleSet, k1)
				self.RuleSet[z.Domain.String()] = z
				goto loo0
			}

			if (k1 != "" || self.RuleSet[k1].Domain.String() != "") ||
				(self.RuleSet[k1].Domain.String() == "" && k1 != "") {
				self.RuleSet[k1].Domain.SetFromString(k1)
				goto loo0
			}
		}

	loo1:
		for k2, _ := range self.RuleSet[k1].DomainSubrequestSettings {
			if k2 != self.RuleSet[k1].DomainSubrequestSettings[k2].Domain.String() {
				if k2 == "" && self.RuleSet[k1].DomainSubrequestSettings[k2].Domain.String() != "" {
					z := self.RuleSet[k1].DomainSubrequestSettings[k2]
					delete(self.RuleSet[k1].DomainSubrequestSettings, k2)
					self.RuleSet[k1].DomainSubrequestSettings[z.Domain.String()] = z
					goto loo0
				}

				if (k2 != "" || self.RuleSet[k1].DomainSubrequestSettings[k2].Domain.String() != "") ||
					(self.RuleSet[k1].DomainSubrequestSettings[k2].Domain.String() == "" && k2 != "") {
					self.RuleSet[k1].DomainSubrequestSettings[k2].Domain.SetFromString(k2)
					goto loo1
				}
			}
		}
	}
}
