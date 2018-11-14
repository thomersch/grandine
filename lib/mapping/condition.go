package mapping

import "github.com/thomersch/grandine/lib/spatial"

type tagMapFn func(map[string]interface{}) map[string]interface{}
type geomOp func(spatial.Geom) []spatial.Geom

type Condition struct {
	// TODO: make it possible to specify Condition type (node/way/rel)
	key    string
	value  string
	mapper tagMapFn
	op     geomOp
}

func (c *Condition) Matches(kv map[string]interface{}) bool {
	if v, ok := kv[c.key]; ok {
		if len(c.value) == 0 || c.value == v {
			return true
		}
	}
	return false
}

// Map converts an incoming key-value map using the given mapping.
func (c *Condition) Map(kv map[string]interface{}) map[string]interface{} {
	return c.mapper(kv)
}

// Transform applies property mapping and performs geometry operations.
// Can emit multiple features, depending on the operation.
func (c *Condition) Transform(f spatial.Feature) []spatial.Feature {
	if c.op == nil {
		return []spatial.Feature{
			{Props: c.Map(f.Props), Geometry: f.Geometry},
		}
	}

	var (
		fts   []spatial.Feature
		props = c.Map(f.Props)
	)
	for _, ng := range c.op(f.Geometry) {
		fts = append(fts, spatial.Feature{Props: props, Geometry: ng})
	}
	return fts
}
