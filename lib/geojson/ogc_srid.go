package geojson

var (
	ogcSRID = map[string]string{
		"urn:ogc:def:crs:OGC:1.3:CRS84": "4326",
	}
	sridOGC = map[string]string{
		"4326": "urn:ogc:def:crs:OGC:1.3:CRS84",
	}
)
