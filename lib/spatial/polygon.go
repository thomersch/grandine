package spatial

import (
	"container/list"
	"log"
)

// Polygon is a data type for storing simple polygons: One outer ring and an arbitrary number of inner rings.
type Polygon []Line

func (p Polygon) ClipToBBox(nw, se Point) []Geom {
	// Is outer ring fully inside?
	oNW, oSE := p[0].BBox()
	if oNW.X() >= nw.X() && oNW.Y() >= nw.Y() && oSE.X() <= se.X() && oSE.Y() <= se.Y() {
		geom, _ := NewGeom(p)
		return []Geom{geom}
	}

	// TODO: inner ring handling
	clipLn := NewLinesFromSegments(BBoxBorders(nw, se))[0]

	// clockwise ordering first
	// subjLns := orderableLine{ln: p[0]}
	// sort.Sort(subjLns)

	var (
		subjLL = list.New()
		clipLL = list.New()
	)

	type refPoint struct {
		pt      Point
		clipRef *list.Element
	}

	// convert subj and clip slices into linked lists
	for _, subjPt := range p[0] {
		log.Printf("subj: %v", subjPt)
		subjLL.PushBack(refPoint{pt: subjPt})
	}
	for _, clipPt := range clipLn {
		clipLL.PushBack(clipPt)
	}

	for subjPt := subjLL.Front(); subjPt != nil; subjPt = subjPt.Next() {
		if subjPt.Next() == nil {
			// no further segment can be constructed
			break
		}
		subjSeg := Segment{subjPt.Value.(refPoint).pt, subjPt.Next().Value.(refPoint).pt}
		for clipPt := clipLL.Front(); clipPt != nil; clipPt = clipPt.Next() {
			if clipPt.Next() == nil {
				// no further clip segment can be constructed
				break
			}
			clipSeg := Segment{clipPt.Value.(Point), clipPt.Next().Value.(Point)}
			if intsct, isIntsct := subjSeg.Intersection(clipSeg); isIntsct {
				log.Printf("intersection: %v (%v %v)", intsct, subjSeg, clipSeg)
				clipRef := clipLL.InsertAfter(intsct, clipPt)
				clipPt = clipPt.Next()

				subjLL.InsertAfter(refPoint{
					pt:      intsct,
					clipRef: clipRef,
				}, subjPt)
				subjPt = subjPt.Next()
			}
		}
	}

	var (
		lines       []Line
		startIntsct *list.Element
	)
	for subjPt := subjLL.Front(); ; subjPt = nextElemOrWrap(subjLL, subjPt) {
		log.Printf("Processing %v", subjPt.Value.(refPoint))
		if subjPt == nil || (startIntsct != nil && subjPt == startIntsct) {
			break
		}
		// entering intersection
		if !subjPt.Value.(refPoint).pt.InPolygon(Polygon{clipLn}) && nextElemOrWrap(subjLL, subjPt).Value.(refPoint).pt.InPolygon(Polygon{clipLn}) {
			startIntsct = subjPt
			log.Printf("entering")
			var (
				nln        Line
				startingPt = subjPt.Value.(refPoint)
			)
			log.Printf("stop at: %v", startingPt)
			nln = append(nln, startingPt.pt)

			// walk the subject line until there is a leaving intersection
			for subjPt = subjPt.Next(); subjPt != nil; subjPt = subjPt.Next() {
				log.Printf("walking subject: %v", subjPt)
				// don't duplicate points
				if len(nln) == 0 || nln[len(nln)-1] != subjPt.Value.(refPoint).pt {
					log.Printf("app.")
					nln = append(nln, subjPt.Value.(refPoint).pt)
				}

				if len(nln) > 1 && subjPt.Next() != nil && subjPt.Value.(refPoint).pt.InPolygon(Polygon{clipLn}) && !subjPt.Next().Value.(refPoint).pt.InPolygon(Polygon{clipLn}) {
					log.Printf("breche")
					break
				}
				if subjPt.Next() == nil {
					// prevent applying Next()
					break
				}
			}
			log.Printf("ref: %p", subjPt)
			// walk the clip line until starting intersection is reached
			for clipPt := subjPt.Value.(refPoint).clipRef; clipPt != nil; clipPt = clipPt.Next() {
				if clipPt == startingPt.clipRef {
					log.Println("ref reached")
					break
				}
				log.Printf("clip ref %p", clipPt)

				// don't duplicate points
				if len(nln) == 0 || nln[len(nln)-1] != clipPt.Value.(Point) {
					nln = append(nln, clipPt.Value.(Point))
				}
			}
			if len(nln) != 0 {
				lines = append(lines, nln)
			}
		}
	}

	var geoms []Geom
	for _, ln := range lines {
		ng, _ := NewGeom(Polygon{ln})
		geoms = append(geoms, ng)
	}
	return geoms
}

func PolygonsFromLines(ls []Line) []Polygon {
	var mlines []Line
	for lni := range ls {
		var merged bool
		for mi := range mlines {
			if mlines[mi].IsExtendedBy(ls[lni]) {
				mlines[mi] = MergeLines(mlines[mi], ls[lni])
				merged = true
			}
		}
		if !merged {
			mlines = append(mlines, ls[lni])
		}
	}
	var polys []Polygon
	for i := range mlines {
		if !mlines[i].Closed() {
			mlines[i] = append(mlines[i], mlines[i][0])
		}
		polys = append(polys, Polygon{mlines[i]})
	}
	return polys
}

func nextElemOrWrap(l *list.List, elem *list.Element) *list.Element {
	if elem.Next() == nil {
		return l.Front()
	}
	return elem.Next()
}
