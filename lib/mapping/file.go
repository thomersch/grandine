package mapping

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type fileMapKV struct {
	Key   string      `yaml:"key"`
	Value interface{} `yaml:"value"`
	Typ   mapType     `yaml:"type"`
}

type fileMap struct {
	Src  fileMapKV   `yaml:"src"`
	Dest []fileMapKV `yaml:"dest"`
}

type fileMappings []fileMap

type typedField struct {
	Name string
	Typ  mapType
}

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
			dynamicKV = map[string]typedField{}
		)
		for _, kvm := range fm.Dest {
			if dv, ok := kvm.Value.(string); !ok {
				staticKV[kvm.Key] = kvm.Value
			} else {
				if dv[0:1] != "$" {
					staticKV[kvm.Key] = dv
				} else {
					// TODO: this can probably be optimized by generating more specific methods at parse time
					dynamicKV[kvm.Key] = typedField{Name: dv[1:], Typ: kvm.Typ}
				}
			}
		}

		cond := Condition{
			key:   fm.Src.Key,
			value: sv,
		}
		if len(dynamicKV) == 0 {
			sm := staticMapper{staticElems: staticKV}
			cond.mapper = sm.Map
		} else {
			dm := dynamicMapper{staticElems: staticKV, dynamicElems: dynamicKV}
			cond.mapper = dm.Map
		}
		conds = append(conds, cond)
	}
	return conds, nil
}

type staticMapper struct {
	staticElems map[string]interface{}
}

func (sm *staticMapper) Map(_ map[string]string) map[string]interface{} {
	return sm.staticElems
}

type dynamicMapper struct {
	staticElems  map[string]interface{}
	dynamicElems map[string]typedField
}

func (dm *dynamicMapper) Map(src map[string]string) map[string]interface{} {
	var vals = map[string]interface{}{}
	for k, v := range dm.staticElems {
		vals[k] = v
	}
	for keyName, field := range dm.dynamicElems {
		if srcV, ok := src[field.Name]; ok {
			switch field.Typ {
			case mapTypeInt:
				v, err := strconv.Atoi(srcV)
				if err == nil {
					vals[keyName] = v
				}
			default:
				vals[keyName] = srcV
			}
		}
	}
	return vals
}
