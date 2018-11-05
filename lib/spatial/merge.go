package spatial

// MergeFeatures aggregates features that have the same properties, if possible.
func MergeFeatures(fts []Feature) []Feature {
	if len(fts) == 1 {
		return fts
	}
	for {
		startLen := len(fts)
		fts = searchAndMerge(fts)
		if startLen == len(fts) {
			return fts
		}
	}
}

func searchAndMerge(fts []Feature) []Feature {
	ref := &fts[0]
	for i, ft := range fts {
		if i == 0 {
			continue
		}
		if ft.Geometry.typ == ref.Geometry.typ {
			if equalProps(ft.Props, ref.Props) {
				switch ft.Geometry.typ {
				case GeomTypeLineString:
					l, merged := mergeLines(ref.Geometry.MustLineString(), ft.Geometry.MustLineString())
					if merged {
						ref.Geometry = MustNewGeom(l)
						fts = append(fts[:i], fts[i+1:]...)
					}
				}
			}
		}
	}
	return fts
}

func mergeLines(l1, l2 Line) (Line, bool) {
	if l1[len(l1)-1] == l2[0] {
		return append(l1, l2[1:]...), true
	}
	return l1, false
}

func equalProps(p1, p2 map[string]interface{}) bool {
	if len(p1) != len(p2) {
		return false
	}
	for k, v1 := range p1 {
		if v2, ok := p2[k]; !ok {
			return false
		} else {
			if v1 != v2 {
				return false
			}
		}
	}
	return true
}
