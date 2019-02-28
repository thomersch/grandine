package spatial

import "sort"

// MergeFeatures aggregates features that have the same properties, if possible.
func MergeFeatures(fts []Feature) []Feature {
	if len(fts) == 1 {
		return fts
	}

	out := fts[:0]
	buckets := tagBuckets(fts)
	for bid := range buckets {
		for {
			startLen := len(buckets[bid])
			if startLen == 1 {
				break
			}
			buckets[bid] = searchAndMerge(buckets[bid])
			if len(buckets[bid]) == startLen {
				break
			}
		}
	}

	for _, bucket := range buckets {
		out = append(out, bucket...)
	}
	return out
}

func tagBuckets(fts []Feature) [][]Feature {
	var buckets [][]Feature

Outer:
	for _, ft := range fts {
		for bID := range buckets {
			if equalProps(buckets[bID][0].Props, ft.Props) {
				buckets[bID] = append(buckets[bID], ft)
				continue Outer
			}
		}
		buckets = append(buckets, []Feature{ft})
	}
	return buckets
}

type ignoreList []int

func (il ignoreList) search(i int) int {
	return sort.Search(len(il), func(pos int) bool { return il[pos] >= i })
}

func (il ignoreList) Has(i int) bool {
	res := il.search(i)
	if res < len(il) && il[res] == i {
		return true
	}
	return false
}

func (il *ignoreList) Add(i int) {
	r := il.search(i)
	iil := *il
	*il = append(iil[:r], append(ignoreList{i}, iil[r:]...)...)
}

func searchAndMerge(fts []Feature) []Feature {
	if len(fts) == 0 {
		return fts
	}
	var ignore = make(ignoreList, 0, len(fts)/10)

	for refID := range fts {
		if ignore.Has(refID) {
			continue
		}
		for i, ft := range fts {
			if ignore.Has(i) || i == refID {
				continue
			}
			if ft.Geometry.typ != fts[refID].Geometry.typ {
				continue
			}
			switch ft.Geometry.typ {
			case GeomTypeLineString:
				l, merged := mergeLines(fts[refID].Geometry.g.(Line), ft.Geometry.g.(Line))
				if merged {
					fts[refID].Geometry.set(l)
					ignore.Add(i)
				}
			}
		}
	}

	var out = make([]Feature, 0, len(fts)-len(ignore))
	for pos, ft := range fts {
		if ignore.Has(pos) {
			continue
		}
		out = append(out, ft)
	}
	return out
}

func mergeLines(l1, l2 Line) (Line, bool) {
	if l1[len(l1)-1] == l2[0] {
		return append(l1, l2[1:]...), true
	}
	if l2[len(l2)-1] == l1[0] {
		return append(l2, l1[1:]...), true
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
