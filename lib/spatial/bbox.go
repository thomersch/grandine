package spatial

type BBox struct {
	SW, NE Point
}

// In determines if the bboxes overlap.
func (b BBox) In(b2 BBox) bool {
	// b.SW.InBBox(b.SW, ne)
	return false

}
