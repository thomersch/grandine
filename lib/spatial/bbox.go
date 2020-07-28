package spatial

import "math"

type BBox struct {
	SW, NE Point
}

func (b1 *BBox) ExtendWith(b2 BBox) {
	b1.SW = Point{math.Min(b1.SW.X, b2.SW.X), math.Min(b1.SW.Y, b2.SW.Y)}
	b1.NE = Point{math.Max(b1.NE.X, b2.NE.X), math.Max(b1.NE.Y, b2.NE.Y)}
}

// In determines if the bboxes overlap.
func (b BBox) Overlaps(b2 BBox) bool {
	return b.SW.InBBox(b2) || b.NE.InBBox(b2) || b2.SW.InBBox(b) || b2.NE.InBBox(b)
}

func (b BBox) FullyIn(b2 BBox) bool {
	return b.SW.InBBox(b2) && b.NE.InBBox(b2)
}

func (b BBox) Segments() []Segment {
	return []Segment{
		{
			{b.SW.X, b.SW.Y},
			{b.SW.X, b.NE.Y},
		},
		{
			{b.SW.X, b.NE.Y},
			{b.NE.X, b.NE.Y},
		},
		{
			{b.NE.X, b.NE.Y},
			{b.NE.X, b.SW.Y},
		},
		{
			{b.NE.X, b.SW.Y},
			{b.SW.X, b.SW.Y},
		},
	}
}
