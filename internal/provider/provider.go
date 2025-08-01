package provider

type ProviderConfig struct {
	User     string
	Hostname string
}

func GetProviderConfig(provider string) (ProviderConfig, bool) {
	providers := map[string]ProviderConfig{
		"github": {
			User:     "git",
			Hostname: "github.com",
		},
		"gitlab": {
			User:     "git",
			Hostname: "gitlab.com",
		},
		"bitbucket": {
			User:     "git",
			Hostname: "bitbucket.org",
		},
	}

	config, exists := providers[provider]
	return config, exists
}

func GetSupportedProviders() []string {
	return []string{"github", "gitlab", "bitbucket", "generic"}
}
