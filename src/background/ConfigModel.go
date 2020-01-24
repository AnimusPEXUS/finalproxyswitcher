package main

type (
	ConfigModel struct {
		ProxyTargets map[string]*ProxyInfo
	}

	SettingInheritance struct {
		None                     bool
		GetFromParent            bool
		GetFromDefaultSubdomain  bool
		GetFromDefaultSubrequest bool
		GetFromDefaultDomain     bool
		GetFromDefaultRequest    bool
	}

	HttpSetting struct {
		Block       bool
		Allow       bool
		Inheritance SettingInheritance
	}

	ExtensionSettings struct {
		// DefaultHttpSettings    HttpSettings
		// DefaultRequestSettings RequestSettings
		// DefaultProxySettings   ProxySettings
	}

	DomainRequestSetting struct {
	}

	DomainSubRequestSetting struct {
	}

	ProxyInfo struct {
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
