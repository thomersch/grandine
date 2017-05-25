package spatial

import "math"

type Line []Point

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
	}
	return intersection, s.HasPoint(intersection) && s2.HasPoint(intersection)
}

// BBoxBorders returns the lines which describe the rectangle of the BBox.
func BBoxBorders(nw, se Point) [4]Segment {
	return [4]Segment{
		{
			{nw.X(), nw.Y()},
			{nw.Y(), se.X()},
		},
		{
			{nw.Y(), se.X()},
			{se.X(), se.Y()},
		},
		{
			{se.X(), se.Y()},
			{nw.X(), se.Y()},
		},
		{
			{nw.X(), se.Y()},
			{nw.X(), nw.Y()},
		},
	}
}
