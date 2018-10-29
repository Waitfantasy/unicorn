package util

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnv(name string, typ string) (interface{}, error) {
	if env := os.Getenv(name); env != "" {
		switch typ {
		case "int":
			if v, err := strconv.Atoi(env); err != nil {
				return nil, fmt.Errorf("%s environment variable strconv.Atoi(%s) error: %v\n", name, env, err)
			} else {
				return v, nil
			}

		case "uint":
			if v, err := strconv.Atoi(env); err != nil {
				return nil, fmt.Errorf("%s environment variable strconv.Atoi(%s) error: %v\n", name, env, err)
			} else {
				return uint(v), nil
			}

		case "int64":
			if v, err := strconv.ParseInt(env, 10, 64); err != nil {
				return nil, fmt.Errorf("%s environment variable strconv.ParseInt(%s, 10, 64) error: %v\n", name, env, err)
			} else {
				return v, nil
			}

		case "uint64":
			if v, err := strconv.ParseUint(env, 10, 64); err != nil {
				return nil, fmt.Errorf("%s environment variable strconv.ParseUint(%s, 10, 64) error: %v\n", name, env, err)
			} else {
				return v, nil
			}

		case "bool":
			if v, err := strconv.ParseBool(env); err != nil {
				return nil, fmt.Errorf("%s environment variable strconv.ParseBool(%s) error: %v\n", name, env, err)
			} else {
				return v, nil
			}
		case "string":
			return env, nil
		default:
			return env, nil
		}
	} else {
		return nil, fmt.Errorf("missing %s parameter", name)
	}
}
