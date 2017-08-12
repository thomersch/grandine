package spatial

import (
	"container/list"
	"fmt"
)

type pointList struct {
	l *list.List
}

func newPointList() *pointList {
	return &pointList{
		l: list.New(),
	}
}

func (p *pointList) Front() *list.Element {
	return p.l.Front()
}

func (p *pointList) InsertAfter(v interface{}, mark *list.Element) *list.Element {
	return p.l.InsertAfter(v, mark)
}

func (p *pointList) PushBack(v refPoint) *list.Element {
	return p.l.PushBack(v)
}

func (p *pointList) String() string {
	const sep = ", "

	var buf string
	for elem := p.l.Front(); elem != nil; elem = elem.Next() {
		buf += elem.Value.(refPoint).pt.String() + sep
	}
	return buf[:len(buf)-len(sep)]
}

// Polygon is a data type for storing simple polygons.
type Polygon []Line

func (p Polygon) String() string {
	s := "Polygon{"
	for _, line := range p {
		s += fmt.Sprintf("%v", line)
	}
	return s + "}"
}

func (p Polygon) ClipToBBox(bbox BBox) []Geom {
	// Is outer ring fully inside?
	orBBox := p[0].BBox()
	if orBBox.SW.X() >= bbox.SW.X() && orBBox.SW.Y() >= bbox.SW.Y() && orBBox.NE.X() <= bbox.NE.X() && orBBox.NE.Y() <= bbox.NE.Y() {
		geom, _ := NewGeom(p)
		return []Geom{geom}
	}

	// This clipping method uses Weiler-Atherton Polygon Clipping
	// It is implemented with in a two-pass manner.
	// Step 1: All edges (subject and clipping) are traversed and intersections are added.
	// Step 2: Starting with an incoming intersection, all edges of the subject polygon are
	//         iterated and until the clipping region is exited, those edges are collected.
	//         After this the clipping polygon is traversed until the original entering
	//         intersection is reached.
	//         This is repeated until the starting intersection is reached.

	// TODO: inner ring handling
	var (
		clipLn = NewLinesFromSegments(BBoxBorders(bbox.SW, bbox.NE))[0]
		subjLn = p[0]
	)
	// weiler-atherton needs clockwise ordering
	if !subjLn.Clockwise() {
		subjLn.Reverse()
	}
	clipLn.Reverse()

	var (
		subjLL = list.New()
		clipLL = list.New()
	)
	// convert subj and clip slices into linked lists
	for _, subjPt := range subjLn {
		subjLL.PushBack(refPoint{pt: subjPt.RoundedCoords()})
		// log.Printf("subjpt: %v", subjPt.RoundedCoords())

	}
	for _, clipPt := range clipLn {
		clipLL.PushBack(clipPt.RoundedCoords())
		// log.Printf("clippt: %v", clipPt.RoundedCoords())
	}

	// build intersections
	nIntsct := insertIntersections(subjLL, clipLL)
	if nIntsct == 0 {
		// Polygons don't intersect at all and polygon is not fully inside.
		return nil
	}

	var (
		lines       []Line
		startIntsct *list.Element
	)
	for subjPt := subjLL.Front(); ; subjPt = nextElemOrWrap(subjLL, subjPt) {
		if subjPt == nil || (startIntsct != nil && subjPt == startIntsct) {
			break
		}
		// entering intersection
		if !subjPt.Value.(refPoint).pt.InPolygon(Polygon{clipLn}) && nextElemOrWrap(subjLL, subjPt).Value.(refPoint).pt.InPolygon(Polygon{clipLn}) {
			if startIntsct == nil {
				startIntsct = subjPt
			}

			var (
				nln        Line
				startingPt = nextElemOrWrap(subjLL, subjPt).Value.(refPoint)
			)
			nln = append(nln, startingPt.pt)

			// walk the subject line until there is a leaving intersection
			for subjPt = nextElemOrWrap(subjLL, subjPt); ; subjPt = nextElemOrWrap(subjLL, subjPt) {
				// don't duplicate points
				if len(nln) == 0 || nln[len(nln)-1] != subjPt.Value.(refPoint).pt {
					nln = append(nln, subjPt.Value.(refPoint).pt)
				}

				if len(nln) > 1 && !nextElemOrWrap(subjLL, subjPt).Value.(refPoint).pt.InPolygon(Polygon{clipLn}) {
					break
				}
			}
			// walk the clip line until starting intersection is reached
			for clipPt := subjPt.Value.(refPoint).clipRef; ; clipPt = nextElemOrWrap(clipLL, clipPt) {
				if clipPt == nil || clipPt == startingPt.clipRef || startingPt.clipRef == nil {
					break
				}

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
		if len(ln) < 3 {
			continue
		}
		ln.Reverse() // invert back to counter-clockwise order
		ng, _ := NewGeom(Polygon{ln})
		geoms = append(geoms, ng)
	}
	return geoms
}

func insertIntersections(subjLL, clipLL *list.List) int {
	var nIntsct = 0
	for subjPt := subjLL.Front(); subjPt != nil; subjPt = subjPt.Next() {
		subjSeg := Segment{subjPt.Value.(refPoint).pt, nextElemOrWrap(subjLL, subjPt).Value.(refPoint).pt}
		for clipPt := clipLL.Front(); clipPt != nil; clipPt = clipPt.Next() {
			clipSeg := Segment{clipPt.Value.(Point), nextElemOrWrap(clipLL, clipPt).Value.(Point)}
			if intsct, isIntsct := subjSeg.Intersection(clipSeg); isIntsct {
				nIntsct += 1
				// log.Println(">>>>>>>>", intsct)
				clipRef := clipLL.InsertAfter(intsct, clipPt)
				clipPt = clipPt.Next()

				if existingElem := hasPointElement(subjLL, refPoint{pt: intsct}); existingElem == nil {
					subjLL.InsertAfter(refPoint{
						pt:      intsct,
						clipRef: clipRef,
					}, subjPt)
				} else {
					existingElem.Value = refPoint{
						pt:      intsct,
						clipRef: clipRef,
					}
				}

				if subjPt.Next() != nil {
					subjPt = subjPt.Next()
				}
			}
		}
	}
	return nIntsct
}

func nextElemOrWrap(l *list.List, elem *list.Element) *list.Element {
	if elem.Next() == nil {
		return l.Front()
	}
	return elem.Next()
}

type refPoint struct {
	pt      Point
	clipRef *list.Element
}

func hasPointElement(l *list.List, ref1 refPoint) *list.Element {
	for ref2 := l.Front(); ref2 != nil; ref2 = ref2.Next() {
		if ref2.Value.(refPoint).pt == ref1.pt {
			return ref2
		}
	}
	return nil
}
