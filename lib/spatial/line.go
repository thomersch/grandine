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

func (l Line) String() string {
	return l.string()
}

func (l Line) Center() Point {
	var (
		sum float64
		pt  Point
	)

	for i := range l {
		pt1 := l[i]
		pt2 := l[(i+1)%len(l)]
		cross := pt1.X*pt2.Y - pt1.Y*pt2.X
		sum += cross
		pt = Point{((pt1.X + pt2.X) * cross) + pt.X, ((pt1.Y + pt2.Y) * cross) + pt.Y}
	}
	z := 1 / (3 * sum)
	return Point{pt.X * z, pt.Y * z}
}

// Segments splits a Line into Segments (line with two points).
func (l Line) Segments() []Segment {
	var segs = make([]Segment, 0, len(l)-1)
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

func (l Line) BBox() BBox {
	var bb BBox
	bb.SW.X = l[0].X
	bb.SW.Y = l[0].Y
	bb.NE.X = bb.SW.X
	bb.NE.Y = bb.SW.Y

	for _, pt := range l {
		bb.SW.X = math.Min(bb.SW.X, pt.X)
		bb.SW.Y = math.Min(bb.SW.Y, pt.Y)
		bb.NE.X = math.Max(bb.NE.X, pt.X)
		bb.NE.Y = math.Max(bb.NE.Y, pt.Y)
	}
	return bb
}

func (l Line) ClipToBBox(bbox BBox) []Geom {
	lsBBox := l.BBox()
	// Is linestring completely inside bbox?
	if bbox.SW.X <= lsBBox.SW.X && bbox.NE.X >= lsBBox.NE.X &&
		bbox.SW.Y <= lsBBox.SW.Y && bbox.NE.Y >= lsBBox.NE.Y {
		// no clipping necessary
		g, _ := NewGeom(l)
		return []Geom{g}
	}

	// Is linestring fully outside the bbox?
	if lsBBox.NE.X < bbox.SW.X || lsBBox.NE.Y < bbox.SW.Y || lsBBox.SW.X > bbox.NE.X || lsBBox.SW.Y > bbox.NE.Y {
		return []Geom{}
	}

	var cutsegs []Segment
	for _, seg := range l.Segments() {
		if seg.FullyInBBox(bbox.SW, bbox.NE) {
			cutsegs = append(cutsegs, seg)
			continue
		}
		ns := seg.ClipToBBox(bbox.SW, bbox.NE)
		if len(ns) != 0 {
			cutsegs = append(cutsegs, ns...)
		}
	}
	if len(cutsegs) == 0 {
		return nil
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

func (l Line) Clockwise() bool {
	var (
		area float64
		n    = len(l)
	)
	for i := range l {
		area += (l[i].Y + l[(i+1)%n].Y) * (l[i].X - l[(i+1)%n].X)
	}
	return area >= 0
}

// Simplify returns a simplified copy of the Line using the Ramer-Douglas-Peucker algorithm.
func (l Line) Simplify(e float64) Line {
	if len(l) < 3 {
		return l
	}

	var (
		seg     = Segment{l[0], l[len(l)-1]}
		maxDist float64
		index   int
	)
	for i, pt := range l[1 : len(l)-1] {
		dist := seg.DistanceToPt(pt)
		if dist > maxDist {
			maxDist = dist
			index = (i + 1) // starting with the second point
		}
	}

	if maxDist > e {
		// divide line in two sublines
		l1 := l[:index+1].Simplify(e)
		l2 := l[index:].Simplify(e)
		return append(l1[:len(l1)-1], l2...)
	}
	return Line{seg[0], seg[1]}
}

func (l Line) Copy() Line {
	return append(l[:0:0], l...) // https://github.com/go101/go101/wiki/How-to-efficiently-clone-a-slice%3F
}

// MergeLines is deprecated.
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

type Segment [2]Point

func (s *Segment) HasPoint(pt Point) bool {
	var (
		s1 = s[0].RoundedCoords()
		s2 = s[1].RoundedCoords()
	)

	if math.Min(s1.X, s2.X) <= pt.X &&
		pt.X <= math.Max(s1.X, s2.X) &&
		math.Min(s1.Y, s2.Y) <= pt.Y &&
		pt.Y <= math.Max(s1.Y, s2.Y) {
		return true
	}
	return false
}

func (s *Segment) SplitAt(p Point) (Segment, Segment) {
	return Segment{s[0], p}, Segment{p, s[1]}
}

// ClipToBBox returns 0 or 1 Segment which is inside bbox.
func (s *Segment) ClipToBBox(sw, ne Point) []Segment {
	var intersections []Point
	for _, bbrd := range BBoxBorders(sw, ne) {
		if ipt, ok := s.Intersection(bbrd); ok {
			intersections = append(intersections, ipt)
		}
	}
	for i, is := range intersections {
		s1, s2 := s.SplitAt(is)
		if s1.Length() != 0 && s1.FullyInBBox(sw, ne) {
			return []Segment{s1}
		}
		if s2.Length() != 0 && s2.FullyInBBox(sw, ne) {
			return []Segment{s2}
		}
		// segment starts and ends outside bbox
		// TODO: this could probably be solved cleaner
		for ii, iis := range intersections {
			if i == ii {
				continue
			}
			is1, is2 := s1.SplitAt(iis)
			if is1.Length() != 0 && is1.FullyInBBox(sw, ne) {
				return []Segment{is1}
			}
			if is2.Length() != 0 && is2.FullyInBBox(sw, ne) {
				return []Segment{is2}
			}
			is1, is2 = s2.SplitAt(iis)
			if is1.Length() != 0 && is1.FullyInBBox(sw, ne) {
				return []Segment{is1}
			}
			if is2.Length() != 0 && is2.FullyInBBox(sw, ne) {
				return []Segment{is2}
			}
		}
	}
	// no intersection
	return nil
}

func (s *Segment) FullyInBBox(sw, ne Point) bool {
	// TODO: check rounding
	sw = sw.RoundedCoords()
	ne = ne.RoundedCoords()
	return s[0].X >= sw.X && s[0].Y >= sw.Y &&
		s[1].X >= sw.X && s[1].Y >= sw.Y &&
		s[0].X <= ne.X && s[0].Y <= ne.Y &&
		s[1].X <= ne.X && s[1].Y <= ne.Y
}

func (s *Segment) Length() float64 {
	if s[0].X == s[1].X && s[0].Y == s[1].Y {
		return 0
	}
	return math.Sqrt(
		math.Pow(s[0].X-s[1].X, 2) +
			math.Pow(s[0].Y-s[1].Y, 2),
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

// DistanceToPt determines the Segment's perpendicular distance to the point.
func (s Segment) DistanceToPt(p Point) float64 {
	if s[0].X == s[1].X {
		return math.Abs(p.X - s[0].X)
	}
	slope := s[1].Y - s[0].Y/s[1].X - s[0].X
	ict := s[0].Y - slope*s[0].X
	return math.Abs(slope*p.X-p.Y+ict) / math.Sqrt(math.Pow(slope, 2)+1)
}

// BBoxBorders returns the lines which describe the rectangle of the BBox.
// Segments are returned in counter-clockwise order.
func BBoxBorders(sw, ne Point) []Segment {
	return []Segment{
		{
			{sw.X, sw.Y},
			{sw.X, ne.Y},
		},
		{
			{sw.X, ne.Y},
			{ne.X, ne.Y},
		},
		{
			{ne.X, ne.Y},
			{ne.X, sw.Y},
		},
		{
			{ne.X, sw.Y},
			{sw.X, sw.Y},
		},
	}
}
