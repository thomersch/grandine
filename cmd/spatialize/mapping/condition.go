package mapping

type tagMapFn func(map[string]string) map[string]interface{}

type Condition struct {
	// TODO: make it possible to specify Condition type (node/way/rel)
	key    string
	value  string
	mapper tagMapFn
}

func (c *Condition) Matches(kv map[string]string) bool {
	if v, ok := kv[c.key]; ok {
		if len(c.value) == 0 || c.value == v {
			return true
		}
	}
	return false
}

func (c *Condition) Map(kv map[string]string) map[string]interface{} {
	return c.mapper(kv)
}
