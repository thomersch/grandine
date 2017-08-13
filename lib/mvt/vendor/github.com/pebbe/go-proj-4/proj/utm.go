package proj

import (
	"errors"
)

// Get UTM zone for longitude and latitude in degrees.
//
// Reference:
//   UTM Grid Zones of the World compiled by Alan Morton
//   http://www.dmap.co.uk/utmworld.htm
func UTMzone(lng, lat float64) (xzone int, yzone string, err error) {

	if lat < -80 || lat > 84 {
		err = errors.New("Arctic and antarctic region are not in UTM")
		return
	}

	for lng < -180 {
		lng += 360
	}
	for lng > 180 {
		lng -= 360
	}

	xzone = 1 + int((lng + 180) / 6)
	if lat > 72 && lng > 0 && lng < 42 {
		if lng < 9 {
			xzone = 31
		} else if lng < 21 {
			xzone = 33
		} else if lng < 33 {
			xzone = 35
		} else {
			xzone = 37
		}
	}
	if lat > 56 && lat < 64 && lng > 3 && lng < 12 {
		xzone = 32
	}

	yzone = string("CDEFGHJKLMNPQRSTUVWXX"[int((lat + 80) / 8)])

    return
}
