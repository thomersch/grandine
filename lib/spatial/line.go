package spatial

import "math"

// TODO: Consider having a line/segment struct with methods

func PointOnLine(pt Point, ln [2]Point) bool {
	if math.Min(ln[0].X(), ln[1].X()) <= pt.X() &&
		pt.X() <= math.Max(ln[0].X(), ln[1].X()) &&
		math.Min(ln[0].Y(), ln[1].Y()) <= pt.Y() &&
		pt.Y() <= math.Max(ln[0].Y(), ln[1].Y()) {
		return true
	}
	return false
}

// LineIntersect returns a point where the two lines intersect and whether the point is on both segments.
func LineIntersect(l1, l2 [2]Point) (Point, bool) {
	a1, b1, c1 := LineSegmentToCarthesian(l1[0], l1[1])
	a2, b2, c2 := LineSegmentToCarthesian(l2[0], l2[1])

	det := a1*b2 - a2*b1
	if det == 0 {
		// parallel lines
		return Point{}, false
	}

	intersection := Point{
		(b2*c1 - b1*c2) / det,
		(a1*c2 - a2*c1) / det,
	}
	return intersection, PointOnLine(intersection, l1) && PointOnLine(intersection, l2)
}

// BBoxBorders returns the lines which describe the rectangle of the BBox.
func BBoxBorders(nw, se Point) [4][2]Point {
	return [4][2]Point{
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
