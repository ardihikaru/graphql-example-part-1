package marshal

import "encoding/json"

// ToStruct casts interface to designated
func ToStruct(val interface{}, destType interface{}) error {
	claimValStr, _ := json.Marshal(val)
	err := json.Unmarshal(claimValStr, &destType)
	if err != nil {
		return err
	}

	return nil
}
