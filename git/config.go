package git

import "strings"

type GitConfig struct {
	UserName       string
	UserEmail      string
	UserSigningKey string
	CommitGpgSign  bool
}

func (g *Git) GetConfigLocal() (GitConfig, error) {
	return g.getParsedConfig("local")
}

func (g *Git) GetConfigGlobal() (GitConfig, error) {
	return g.getParsedConfig("global")
}

func (g *Git) GetConfigSystem() (GitConfig, error) {
	return g.getParsedConfig("system")
}

func (g *Git) getParsedConfig(list string) (GitConfig, error) {
	cfg, err := g.getConfig(list)
	if err != nil {
		return GitConfig{}, err
	}
	return g.parseConfig(cfg), nil
}

func (g *Git) getConfig(list string) (map[string]string, error) {
	c, err := g.RunCwd("--no-pager", "config", "--list", "--"+list, "-z")
	if err != nil {
		return nil, err
	}

	cfg := make(map[string]string)
	opts := strings.Split(c, "\x00")
	for _, opt := range opts {
		kv := strings.Split(opt, "\n")
		if len(kv) == 2 {
			cfg[kv[0]] = kv[1]
		}
	}

	return cfg, nil
}

func (g *Git) parseConfig(cfg map[string]string) GitConfig {
	config := GitConfig{}

	for k, v := range cfg {
		switch k {
		case "user.name":
			config.UserName = v
		case "user.email":
			config.UserEmail = v
		case "user.signingkey":
			config.UserSigningKey = v
		case "commit.gpgsign":
			config.CommitGpgSign = v == "true"
		}
	}

	return config
}

func (g *Git) UpdateUserName(list string, value string) error {
	return g.updateConfig(list, "user.name", value)
}

func (g *Git) UpdateUserEmail(list string, value string) error {
	return g.updateConfig(list, "user.email", value)
}

func (g *Git) UpdateUserSigningKey(list string, value string) error {
	return g.updateConfig(list, "user.signingkey", value)
}

func (g *Git) UpdateCommitGpgSign(list string, value bool) error {
	var value_s string
	if value {
		value_s = "true"
	} else {
		value_s = "false"
	}
	return g.updateConfig(list, "commit.gpgsign", value_s)
}

func (g *Git) updateConfig(list string, key string, value string) error {
	// list "system" not currently supported
	if list != "global" {
		list = "local"
	}
	if value == "" {
		return g.clearConfig(list, key)
	}
	_, err := g.RunCwd("config", "--"+list, key, value)
	return err
}

func (g *Git) clearConfig(list string, key string) error {
	_, err := g.RunCwd("config", "--"+list, "--unset-all", key)
	return err
}
