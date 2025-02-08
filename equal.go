package jsontools

import (
	"encoding/json"
	"reflect"
)

func JsonEqual(a, b []byte) (bool, error) {
	filter := NewJsonNullFilter(false)
	a, err := filter.Filter(a)
	if err != nil {
		return false, err
	}
	b, err = filter.Filter(b)
	if err != nil {
		return false, err
	}

	var o1, o2 any
	err = json.Unmarshal(a, &o1)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, &o2)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(o1, o2), nil
}
