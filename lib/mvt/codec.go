package mvt

import "github.com/thomersch/grandine/lib/spatial"

type cmd uint32

const (
	cmdMoveTo    cmd = 1
	cmdLineTo    cmd = 2
	cmdClosePath cmd = 7

	extent = 4096
)

func encodeCommandInt(c cmd, count uint32) uint32 {
	return (uint32(c) & 0x7) | (count << 3)
}

func decodeCommandInt(cmdInt uint32) (cmd, uint32) {
	return cmd(cmdInt & 0x7), cmdInt >> 3
}

func encodeZigZag(i int) uint32 {
	return uint32((i << 1) ^ (i >> 31))
}

func encodeGeometry(fs []spatial.Feature) (commands []uint32) {
	var (
		cur    [2]uint32
		dx, dy int32
	)
	for _, feat := range fs {
		switch feat.Geometry.Typ() {
		case spatial.GeomTypePoint:
			pt, _ := feat.Geometry.Point()
			dx = encodeZigZag(pt[0] - cur[0])
			dy = encodeZigZag(pt[1] - cur[1])
			commands = append(commands, encodeCommandInt(cmdMoveTo, 1), dx, dy)
		}
	}
	return
}
