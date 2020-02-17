package spatial

type BBox struct {
	SW, NE Point
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
