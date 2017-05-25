package spatial

// LineSegmentToCarthesian converts a line segment into a carthesian representation.
// Possible improvement: normalization of values
func LineSegmentToCarthesian(p1, p2 Point) (a, b, c float64) {
	a = p1.Y() - p2.Y()
	b = p2.X() - p1.X()
	c = p2.X()*p1.Y() - p1.X()*p2.Y()
	return
}
