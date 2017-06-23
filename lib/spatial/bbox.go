package spatial

type BBox struct {
	SW, NE Point
}

// In determines if the bboxes overlap.
func (b BBox) In(b2 BBox) bool {
	return b.SW.InBBox(b2) || b.NE.InBBox(b2) || b2.SW.InBBox(b) || b2.NE.InBBox(b)
}
