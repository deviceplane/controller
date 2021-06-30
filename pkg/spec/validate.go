package spec

import (
	"fmt"

	"github.com/deviceplane/controller/pkg/validation"
	"gopkg.in/yaml.v2"
)

var (
	validators = map[string][]func(interface{}) error{
		"cap_add":          {validation.ValidateStringArray},
		"cap_drop":         {validation.ValidateStringArray},
		"command":          {validation.ValidateStringOrStringArray},
		"cpuset":           {validation.ValidateString},
		"cpu_shares":       {validation.ValidateStringOrInteger},
		"cpu_quota":        {validation.ValidateStringOrInteger},
		"devices":          {validation.ValidateStringArray},
		"dns":              {validation.ValidateStringOrStringArray},
		"dns_opt":          {validation.ValidateStringOrStringArray},
		"dns_search":       {validation.ValidateStringOrStringArray},
		"domainname":       {validation.ValidateString},
		"entrypoint":       {validation.ValidateStringOrStringArray},
		"environment":      {validation.ValidateArrayOrObject},
		"extra_hosts":      {validation.ValidateArrayOrObject},
		"group_add":        {validation.ValidateStringIntegerArray},
		"image":            {validation.ValidateString},
		"hostname":         {validation.ValidateString},
		"ipc":              {validation.ValidateString},
		"labels":           {validation.ValidateArrayOrObject},
		"mem_limit":        {validation.ValidateStringOrInteger},
		"mem_reservation":  {validation.ValidateStringOrInteger},
		"memswap_limit":    {validation.ValidateStringOrInteger},
		"network_mode":     {validation.ValidateString},
		"oom_kill_disable": {validation.ValidateBoolean},
		"oom_score_adj":    {validation.ValidateInteger},
		"pid":              {validation.ValidateString},
		"ports":            {validation.ValidateStringIntegerArray},
		"privileged":       {validation.ValidateBoolean},
		"read_only":        {validation.ValidateBoolean},
		"restart":          {validation.ValidateString},
		"runtime":          {validation.ValidateString},
		"security_opt":     {validation.ValidateStringArray},
		"shm_size":         {validation.ValidateStringOrInteger},
		"stop_signal":      {validation.ValidateString},
		"user":             {validation.ValidateString},
		"uts":              {validation.ValidateString},
		"volumes":          {validation.ValidateStringArray},
		"working_dir":      {validation.ValidateString},
	}
)

func Validate(c []byte) error {
	var m map[string]interface{}
	if err := yaml.Unmarshal(c, &m); err != nil {
		return err
	}

	for serviceName := range m {
		if len(serviceName) > 100 {
			return fmt.Errorf("service name '%s' is longer than 100 characters", serviceName)
		}
	}

	for serviceName, service := range m {
		service, ok := service.(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("service '%s' is not an object", serviceName)
		}

		for key := range service {
			typedKey, ok := key.(string)
			if !ok {
				return fmt.Errorf("service '%s': invalid key '%v'", serviceName, key)
			}
			if _, ok = validators[typedKey]; !ok {
				return fmt.Errorf("service '%s': invalid key '%s'", serviceName, typedKey)
			}
		}

		for key, validators := range validators {
			value, ok := service[key]
			if !ok {
				continue
			}
			for _, validator := range validators {
				if err := validator(value); err != nil {
					return fmt.Errorf("service '%s', key '%s': %v", serviceName, key, err)
				}
			}
		}
	}

	return nil
}
