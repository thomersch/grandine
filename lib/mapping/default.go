package mapping

var (
	transportationMapFn = func(kv map[string]interface{}) map[string]interface{} {
		var cl string
		if class, ok := kv["highway"]; ok {
			cl = class.(string)
		}
		return map[string]interface{}{
			"@layer": "transportation",
			"class":  cl,
		}
	}

	landuseMapFn = func(kv map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"__type": "area",
			"@layer": "landcover",
			"class":  "wood",
		}
	}

	aerowayMapFn = func(kv map[string]interface{}) map[string]interface{} {
		var cl string
		if class, ok := kv["aeroway"]; ok {
			cl = class.(string)
		}
		return map[string]interface{}{
			"@layer": "aeroway",
			"class":  cl,
		}
	}

	buildingMapFn = func(kv map[string]interface{}) map[string]interface{} {
		return map[string]interface{}{
			"@layer":    "building",
			"@zoom:min": 14,
		}
	}

	waterwayMapFn = func(kv map[string]interface{}) map[string]interface{} {
		var cl string
		if class, ok := kv["waterway"]; ok {
			cl = class.(string)
		}
		return map[string]interface{}{
			"@layer": "waterway",
			"class":  cl,
		}
	}

	Default = []Condition{
		{"aeroway", "aerodrome", aerowayMapFn},
		{"aeroway", "apron", aerowayMapFn},
		{"aeroway", "heliport", aerowayMapFn},
		{"aeroway", "runway", aerowayMapFn},
		{"aeroway", "helipad", aerowayMapFn},
		{"aeroway", "taxiway", aerowayMapFn},
		{"highway", "motorway", transportationMapFn},
		{"highway", "primary", transportationMapFn},
		{"highway", "trunk", transportationMapFn},
		{"highway", "secondary", transportationMapFn},
		{"highway", "tertiary", transportationMapFn},
		{"building", "", buildingMapFn},
		{"landuse", "forest", landuseMapFn},
		{"railway", "rail", transportationMapFn},
		{"waterway", "river", waterwayMapFn},
	}
)
