// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

type Config struct {
	Period       time.Duration `config:"period"`
	AccessToken  string        `config:"access_token"`
	GitlabAdress string        `config:"gitlab_address"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
}
