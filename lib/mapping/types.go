package mapping

import "fmt"

type mapType uint8

const (
	mapTypeString mapType = iota
	mapTypeInt
)

func (mt *mapType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ts string
	err := unmarshal(&ts)
	if err != nil {
		return err
	}
	switch ts {
	case "int":
		*mt = mapTypeInt
	case "string":
		*mt = mapTypeString
	default:
		return fmt.Errorf("unknown type: %s (allowed values: int, string)", ts)
	}
	return nil
}

func InterfaceMap(i map[string]string) (om map[string]interface{}) {
	for k, v := range i {
		om[k] = v
	}
	return
}
