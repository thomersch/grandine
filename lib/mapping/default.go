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
		{"aeroway", "aerodrome", aerowayMapFn, nil},
		{"aeroway", "apron", aerowayMapFn, nil},
		{"aeroway", "heliport", aerowayMapFn, nil},
		{"aeroway", "runway", aerowayMapFn, nil},
		{"aeroway", "helipad", aerowayMapFn, nil},
		{"aeroway", "taxiway", aerowayMapFn, nil},
		{"highway", "motorway", transportationMapFn, nil},
		{"highway", "primary", transportationMapFn, nil},
		{"highway", "trunk", transportationMapFn, nil},
		{"highway", "secondary", transportationMapFn, nil},
		{"highway", "tertiary", transportationMapFn, nil},
		{"building", "", buildingMapFn, nil},
		{"landuse", "forest", landuseMapFn, nil},
		{"railway", "rail", transportationMapFn, nil},
		{"waterway", "river", waterwayMapFn, nil},
	}
)
