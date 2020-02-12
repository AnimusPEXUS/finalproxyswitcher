package main

type (
	ConfigModel struct {
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
