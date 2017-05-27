package spatial

import "math"

type Line []Point

// NewLinesFromSegments creates a line from continous segments.
func NewLinesFromSegments(segs []Segment) []Line {
	ls := []Line{Line{}}
	for i, seg := range segs {
		// if start of current segment isn't the end of the last one, create a new line
		if i != 0 && seg[0] != segs[i-1][1] {
			ls[len(ls)-1] = append(ls[len(ls)-1], segs[i-1][1])
			ls = append(ls, Line{seg[0]})
			continue
		}
		ls[len(ls)-1] = append(ls[len(ls)-1], seg[0])
	}
	ls[len(ls)-1] = append(ls[len(ls)-1], segs[len(segs)-1][1]) // last point of last segment is added to the last linestring
	return ls
}

func (l Line) Center() Point {
	var (
		sum float64
		pt  Point
	)

	for i := range l {
		pt1 := l[i]
		pt2 := l[(i+1)%len(l)]
		cross := pt1.X()*pt2.Y() - pt1.Y()*pt2.X()
		sum += cross
		pt = Point{((pt1.X() + pt2.X()) * cross) + pt.X(), ((pt1.Y() + pt2.Y()) * cross) + pt.Y()}
	}
	z := 1 / (3 * sum)
	return Point{pt.X() * z, pt.Y() * z}
}

// Segments splits a Line into Segments (line with two points).
func (l Line) Segments() []Segment {
	var segs []Segment
	for i := range l {
		if i+1 < len(l) {
			segs = append(segs, Segment{l[i], l[i+1]})
		}
	}
	return segs
}

func (l Line) Intersections(segments []Segment) []Point {
	var intersectionSet = map[Point]struct{}{} // set
	for _, seg := range l.Segments() {
		for _, seg2 := range segments {
			if ipt, intersects := seg.Intersection(seg2); intersects {
				intersectionSet[ipt] = struct{}{}
			}
		}
	}

	var intersections []Point
	for inter := range intersectionSet {
		intersections = append(intersections, inter)
	}
	return intersections
}

func (l Line) BBox() (nw, se Point) {
	nw[0] = l[0][0]
	nw[1] = l[0][1]
	se[0] = nw[0]
	se[1] = nw[1]
	for _, pt := range l {
		nw[0] = math.Min(nw[0], pt[0])
		nw[1] = math.Min(nw[1], pt[1])
		se[0] = math.Max(se[0], pt[0])
		se[1] = math.Max(se[1], pt[1])
	}
	return
}

func (l Line) ClipToBBox(nw, se Point) []Geom {
	lsNW, lsSE := l.BBox()
	// Is linestring completely inside bbox?
	if nw[0] <= lsNW[0] && se[0] >= lsSE[0] &&
		nw[1] <= lsNW[1] && se[1] >= lsSE[1] {
		// no clipping necessary
		g, _ := NewGeom(l)
		return []Geom{g}
	}

	// Is linestring fully outside the bbox?
	if lsSE[0] < nw[0] || lsSE[1] < nw[1] || lsNW[0] > se[0] || lsNW[1] > se[1] {
		return []Geom{}
	}

	var cutsegs []Segment
	for _, seg := range l.Segments() {
		if seg.FullyInBBox(nw, se) {
			cutsegs = append(cutsegs, seg)
			continue
		}
		ns := seg.ClipToBBox(nw, se)
		if len(ns) != 0 {
			cutsegs = append(cutsegs, ns...)
		}
	}
	var gms []Geom
	for _, ln := range NewLinesFromSegments(cutsegs) {
		gms = append(gms, Geom{typ: GeomTypeLineString, g: ln})
	}
	return gms
}

func (l Line) Closed() bool {
	if len(l) < 2 {
		return false
	}
	return l[0] == l[len(l)-1]
}

func (l1 Line) IsExtendedBy(l2 Line) bool {
	return l1[0] == l2[0] || l1[1] == l2[1] || l1[len(l1)-1] == l2[0] || l1[0] == l2[len(l2)-1]
}

func (l Line) Reverse() {
	for i := len(l)/2 - 1; i >= 0; i-- {
		opp := len(l) - 1 - i
		l[i], l[opp] = l[opp], l[i]
	}
}

func MergeLines(l1, l2 Line) Line {
	// head to head
	if l1[0] == l2[0] {
		l1.Reverse()
		return append(l1, l2[1:]...)
	}
	// tail to tail
	if l1[len(l1)-1] == l2[len(l2)-1] {
		l2.Reverse()
		return append(l1, l2[1:]...)
	}
	// head to tail
	if l1[0] == l2[len(l2)-1] {
		return append(l2, l1[1:]...)
	}
	// tail to head
	if l1[len(l1)-1] == l2[0] {
		return append(l1, l2[1:]...)
	}
	return Line{}
}

type orderableLine struct {
	ln     Line
	center Point
}

func newOrderableLine(l Line) orderableLine {
	return orderableLine{ln: l, center: l.Center()}
}

// Methods for sorting in a clockwise order
func (l orderableLine) Len() int      { return len(l.ln) }
func (l orderableLine) Swap(i, j int) { l.ln[i], l.ln[j] = l.ln[j], l.ln[i] }
func (l orderableLine) Less(i, j int) bool {
	// inspired by https://stackoverflow.com/a/6989383/552651
	var (
		center = l.center
		b      = l.ln[i]
		a      = l.ln[j]
	)
	if a.X()-center.X() >= 0 && b.X()-center.X() < 0 {
		return true
	}
	if a.X()-center.X() < 0 && b.X()-center.X() >= 0 {
		return false
	}
	if a.X()-center.X() == 0 && b.X()-center.X() == 0 {
		if a.Y()-center.Y() >= 0 || b.Y()-center.Y() >= 0 {
			return a.Y() > b.Y()
		}
		return b.Y() > a.Y()
	}
	det := (a.X()-center.X())*(b.Y()-center.Y()) - (b.X()-center.X())*(a.Y()-center.Y())
	if det < 0 {
		return true
	}
	if det > 0 {
		return false
	}
	d1 := (a.X()-center.X())*(a.X()-center.X()) + (a.Y()-center.Y())*(a.Y()-center.Y())
	d2 := (b.X()-center.X())*(b.X()-center.X()) + (b.Y()-center.Y())*(b.Y()-center.Y())
	return d1 > d2
}

type Segment [2]Point

func (s *Segment) HasPoint(pt Point) bool {
	if math.Min(s[0].X(), s[1].X()) <= pt.X() &&
		pt.X() <= math.Max(s[0].X(), s[1].X()) &&
		math.Min(s[0].Y(), s[1].Y()) <= pt.Y() &&
		pt.Y() <= math.Max(s[0].Y(), s[1].Y()) {
		return true
	}
	return false
}

func (s *Segment) SplitAt(p Point) (Segment, Segment) {
	return Segment{s[0], p}, Segment{p, s[1]}
}

// ClipToBBox returns 0 or 1 Segment which is inside bbox.
func (s *Segment) ClipToBBox(nw, se Point) []Segment {
	var intersections []Point
	for _, bbrd := range BBoxBorders(nw, se) {
		if ipt, ok := s.Intersection(bbrd); ok {
			intersections = append(intersections, ipt)
		}
	}
	for i, is := range intersections {
		s1, s2 := s.SplitAt(is)
		if s1.Length() != 0 && s1.FullyInBBox(nw, se) {
			return []Segment{s1}
		}
		if s2.Length() != 0 && s2.FullyInBBox(nw, se) {
			return []Segment{s2}
		}
		// segment starts and ends outside bbox
		// TODO: this could probably be solved cleaner
		for ii, iis := range intersections {
			if i == ii {
				continue
			}
			is1, is2 := s1.SplitAt(iis)
			if is1.Length() != 0 && is1.FullyInBBox(nw, se) {
				return []Segment{is1}
			}
			if is2.Length() != 0 && is2.FullyInBBox(nw, se) {
				return []Segment{is2}
			}
			is1, is2 = s2.SplitAt(iis)
			if is1.Length() != 0 && is1.FullyInBBox(nw, se) {
				return []Segment{is1}
			}
			if is2.Length() != 0 && is2.FullyInBBox(nw, se) {
				return []Segment{is2}
			}
		}
	}
	// no intersection
	return nil
}

func (s *Segment) FullyInBBox(nw, se Point) bool {
	return s[0].X() >= nw.X() && s[0].Y() >= nw.Y() &&
		s[1].X() >= nw.X() && s[1].Y() >= nw.Y() &&
		s[0].X() <= se.X() && s[0].Y() <= se.Y() &&
		s[1].X() <= se.X() && s[1].Y() <= se.Y()
}

func (s *Segment) Length() float64 {
	if s[0].X() == s[1].X() && s[0].Y() == s[1].Y() {
		return 0
	}
	return math.Sqrt(
		math.Pow(s[0].X()-s[1].X(), 2) +
			math.Pow(s[0].Y()-s[1].Y(), 2),
	)
}

// LineIntersect returns a point where the two lines intersect and whether the point is on both segments.
func (s *Segment) Intersection(s2 Segment) (Point, bool) {
	a1, b1, c1 := LineSegmentToCarthesian(s[0], s[1])
	a2, b2, c2 := LineSegmentToCarthesian(s2[0], s2[1])

	det := a1*b2 - a2*b1
	if det == 0 {
		// parallel lines
		return Point{}, false
	}

	intersection := Point{
		(b2*c1 - b1*c2) / det,
		(a1*c2 - a2*c1) / det,
	}.RoundedCoords()
	return intersection, s.HasPoint(intersection) && s2.HasPoint(intersection)
}

// BBoxBorders returns the lines which describe the rectangle of the BBox.
// Segments are returned in counter-clockwise order.
func BBoxBorders(nw, se Point) [4]Segment {
	return [4]Segment{
		{
			{nw.X(), nw.Y()},
			{nw.X(), se.Y()},
		},
		{
			{nw.X(), se.Y()},
			{se.X(), se.Y()},
		},
		{
			{se.X(), se.Y()},
			{se.X(), nw.Y()},
		},
		{
			{se.X(), nw.Y()},
			{nw.X(), nw.Y()},
		},
	}
}
