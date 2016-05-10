package tweethour

import (
	"errors"
)

func getValue(key string, v map[string]interface{}) (string, error) {
	valI, ok := v[key]
	if !ok {
		return "", errors.New("Unable to get key : " + key)
	}

	var val string
	val, ok = valI.(string)

	if !ok {
		return "", errors.New("Unable to get key : " + key)
	}

	return val, nil
}
