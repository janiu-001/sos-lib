package util

import "encoding/json"

func ConvertStructData(src interface{}, dst ...interface{}) error {
	jsonStr, err := json.Marshal(src)
	if err != nil {
		return err
	}

	for _, v := range dst {
		if err := json.Unmarshal(jsonStr, v); err != nil {
			return err
		}
	}
	return nil
}
