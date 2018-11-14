package mapping

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	Op   string      `yaml:"op"`
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

		switch fm.Op {
		case "lines":
			cond.op = polyToLines
		}

		conds = append(conds, cond)
	}
	return conds, nil
}

type staticMapper struct {
	staticElems map[string]interface{}
}

func (sm *staticMapper) Map(_ map[string]interface{}) map[string]interface{} {
	return sm.staticElems
}

type dynamicMapper struct {
	staticElems  map[string]interface{}
	dynamicElems map[string]typedField
}

func (dm *dynamicMapper) Map(src map[string]interface{}) map[string]interface{} {
	var (
		vals = map[string]interface{}{}
		err  error
	)
	for k, v := range dm.staticElems {
		vals[k] = v
	}
	for keyName, field := range dm.dynamicElems {
		if srcV, ok := src[field.Name]; ok {
			switch field.Typ {
			case mapTypeInt:
				vals[keyName], err = dm.toInt(srcV)
			default:
				vals[keyName] = srcV
			}
			if err != nil {
				log.Println(err) // Not sure if this won't get too verbose. Let's keep it here for some time.
				vals[keyName] = srcV
				err = nil
			}
		}
	}
	return vals
}

func (dm *dynamicMapper) toInt(i interface{}) (int, error) {
	switch v := i.(type) {
	case string:
		k, err := strconv.Atoi(v)
		if err == nil {
			return k, nil
		}
		if v == "yes" {
			return 1, nil
		}
		if v == "no" {
			return 0, nil
		}
		return 0, err
	case int:
		return v, nil
	default:
		return 0, fmt.Errorf("cannot convert %v (type %T) to int", v, v)
	}
}
