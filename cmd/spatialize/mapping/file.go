package mapping

import (
	"fmt"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type fileMapKV struct {
	Key   string      `yaml:"key"`
	Value interface{} `yaml:"value"`
}

type fileMap struct {
	Src  fileMapKV   `yaml:"src"`
	Dest []fileMapKV `yaml:"dest"`
}

type fileMappings []fileMap

func ParseMapping(r io.Reader) ([]Condition, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var fms fileMappings
	err = yaml.UnmarshalStrict(buf, &fms)
	if err != nil {
		return nil, err
	}

	var conds []Condition
	for _, fm := range fms {
		sv, ok := fm.Src.Value.(string)
		if !ok {
			return nil, fmt.Errorf("source key %s must be of type string (has: %v)", fm.Src.Key, fm.Src.Value)
		}
		if sv == "*" {
			sv = ""
		}

		var (
			staticKV  = map[string]interface{}{}
			dynamicKV = map[string]string{}
		)
		for _, kvm := range fm.Dest {
			if dv, ok := kvm.Value.(string); !ok {
				staticKV[kvm.Key] = dv
			} else {
				if dv != "$" {
					staticKV[kvm.Key] = dv
				} else {
					dynamicKV[kvm.Key] = dv
				}
			}
		}

		conds = append(conds, Condition{
			key:   fm.Src.Key,
			value: sv,
			mapper: func(map[string]string) map[string]interface{} {
				// TODO: auto-detect if there are any dynamic mappings and apply more specific (and faster) functions
				return nil
			},
		})
	}
	return conds, nil
}
