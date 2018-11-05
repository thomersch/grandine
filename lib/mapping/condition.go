package mapping

type tagMapFn func(map[string]interface{}) map[string]interface{}

type Condition struct {
	// TODO: make it possible to specify Condition type (node/way/rel)
	key    string
	value  string
	mapper tagMapFn
}

func (c *Condition) Matches(kv map[string]interface{}) bool {
	if v, ok := kv[c.key]; ok {
		if len(c.value) == 0 || c.value == v {
			return true
		}
	}
	return false
}

func (c *Condition) Map(kv map[string]interface{}) map[string]interface{} {
	return c.mapper(kv)
}
