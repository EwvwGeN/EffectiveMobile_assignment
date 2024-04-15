package config

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type ValidatorConfig struct {
	RegisterNumberRegex  string `yaml:"reg_num"`
	MarkRegex            string `yaml:"mark"`
	ModelRegex           string `yaml:"model"`
	OwnerNameRegex       string `yaml:"owner_name"`
	OwnerSurnameRegex    string `yaml:"owner_surname"`
	OwnerPatronymicRegex string `yaml:"owner_patronymic"`
}

func (v *ValidatorConfig) checkRegex() error {
	b, _ := json.Marshal(v)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	for i, v := range m {
		_, err := regexp.Compile(v.(string))
		if err != nil {
			return fmt.Errorf("incorrect %s", i)
		}
	}
	return nil
}