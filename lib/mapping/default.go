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
		{"aeroway", []string{"aerodrome"}, aerowayMapFn, nil},
		{"aeroway", []string{"apron"}, aerowayMapFn, nil},
		{"aeroway", []string{"heliport"}, aerowayMapFn, nil},
		{"aeroway", []string{"runway"}, aerowayMapFn, nil},
		{"aeroway", []string{"helipad"}, aerowayMapFn, nil},
		{"aeroway", []string{"taxiway"}, aerowayMapFn, nil},
		{"highway", []string{"motorway"}, transportationMapFn, nil},
		{"highway", []string{"primary"}, transportationMapFn, nil},
		{"highway", []string{"trunk"}, transportationMapFn, nil},
		{"highway", []string{"secondary"}, transportationMapFn, nil},
		{"highway", []string{"tertiary"}, transportationMapFn, nil},
		{"building", []string{""}, buildingMapFn, nil},
		{"landuse", []string{"forest"}, landuseMapFn, nil},
		{"railway", []string{"rail"}, transportationMapFn, nil},
		{"waterway", []string{"river"}, waterwayMapFn, nil},
	}
)
