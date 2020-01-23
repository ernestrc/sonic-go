package sonic

import "fmt"

func getFloatInPayload(payloadMap map[string]interface{}, key string) (float64, error) {
	valueIfc, ok := payloadMap[key]
	if !ok {
		return 0, nil
	}

	value, ok := valueIfc.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid type for field %s: exptected float64", key)
	}
	return value, nil
}

func getStringInPayload(payloadMap map[string]interface{}, key string) (string, error) {
	valueIfc, ok := payloadMap[key]
	if !ok {
		return "", nil
	}

	value, ok := valueIfc.(string)
	if !ok {
		return "", fmt.Errorf("invalid type for field %s: expected string", key)
	}
	return value, nil
}
