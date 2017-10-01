package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/thomersch/grandine/cmd/spatialize/mapping"
	"github.com/thomersch/grandine/lib/spaten"
	"github.com/thomersch/grandine/lib/spatial"

	"github.com/thomersch/gosmparse"
)

type nd struct {
	Lat, Lon float64
	Tags     map[string]interface{}
}
type wy struct {
	ID      int64
	NodeIDs []int64
	Tags    map[string]interface{}
}
type rl struct {
	Members []gosmparse.RelationMember
	Tags    map[string]interface{}
}

type dataHandler struct {
	conds []mapping.Condition

	ec *elemCache

	nodes    []nd
	nodesMtx sync.Mutex
	ways     []wy
	waysMtx  sync.Mutex
	rels     []rl
	relsMtx  sync.Mutex
}

func (d *dataHandler) ReadNode(n gosmparse.Node) {
	for _, cond := range d.conds {
		if cond.Matches(n.Tags) {
			d.nodesMtx.Lock()
			d.nodes = append(d.nodes, nd{
				Lat:  n.Lat,
				Lon:  n.Lon,
				Tags: cond.Map(n.Tags),
			})
			d.nodesMtx.Unlock()
		}
	}
}

func (d *dataHandler) ReadWay(w gosmparse.Way) {
	for _, cond := range d.conds {
		if cond.Matches(w.Tags) {
			d.ec.AddNodes(w.NodeIDs...)
			d.ec.setMembers(w.ID, w.NodeIDs)

			d.waysMtx.Lock()
			d.ways = append(d.ways, wy{
				ID:      w.ID,
				NodeIDs: w.NodeIDs,
				Tags:    cond.Map(w.Tags),
			})
			d.waysMtx.Unlock()
		}
	}
}

func (d *dataHandler) ReadRelation(r gosmparse.Relation) {
	for _, cond := range d.conds {
		if cond.Matches(r.Tags) {
			d.relsMtx.Lock()
			d.rels = append(d.rels, rl{
				Members: r.Members,
				Tags:    cond.Map(r.Tags),
			})
			d.relsMtx.Unlock()

			for _, memb := range r.Members {
				switch memb.Type {
				case gosmparse.WayType:
					d.ec.AddWay(memb.ID)
				} // TODO: check if relations of nodes/relations are necessary
			}
		}
	}
}

type elemCache struct {
	nodes    map[int64]spatial.Point
	nodesMtx sync.Mutex
	ways     map[int64][]int64
	waysMtx  sync.Mutex
}

func NewElemCache() *elemCache {
	return &elemCache{
		nodes: map[int64]spatial.Point{},
		ways:  map[int64][]int64{},
	}
}

func (d *elemCache) AddNodes(nIDs ...int64) {
	d.nodesMtx.Lock()
	for _, nID := range nIDs {
		d.nodes[nID] = spatial.Point{}
	}
	d.nodesMtx.Unlock()
}

func (d *elemCache) AddWay(wID int64) {
	d.waysMtx.Lock()
	d.ways[wID] = []int64{}
	d.waysMtx.Unlock()
}

func (d *elemCache) SetCoord(nID int64, coord spatial.Point) {
	d.nodesMtx.Lock()
	d.nodes[nID] = coord
	d.nodesMtx.Unlock()
}

func (d *elemCache) setMembers(wID int64, members []int64) {
	d.waysMtx.Lock()
	d.ways[wID] = members
	d.waysMtx.Unlock()
}

func (d *elemCache) ReadWay(w gosmparse.Way) {
	d.waysMtx.Lock()
	_, ok := d.ways[w.ID]
	d.waysMtx.Unlock()
	if ok {
		d.setMembers(w.ID, w.NodeIDs)
		d.AddNodes(w.NodeIDs...)
	}
}

func (d *elemCache) Line(wID int64) spatial.Line {
	// check if mutex is needed
	membs, ok := d.ways[wID]
	if !ok {
		log.Fatalf("missing referenced way: %v", wID)
	}

	var l spatial.Line
	for _, memb := range membs {
		l = append(l, d.nodes[memb])
	}
	return l
}

// Interface enforces this. Probably I should change the behavior.
func (d *elemCache) ReadNode(n gosmparse.Node)         {}
func (d *elemCache) ReadRelation(r gosmparse.Relation) {}

type nodeCollector struct {
	ec *elemCache
}

func (d *nodeCollector) ReadNode(n gosmparse.Node) {
	d.ec.SetCoord(n.ID, spatial.Point{float64(n.Lon), float64(n.Lat)})
}
func (d *nodeCollector) ReadWay(w gosmparse.Way)           {}
func (d *nodeCollector) ReadRelation(r gosmparse.Relation) {}

func main() {
	source := flag.String("in", "osm.pbf", "")
	outfile := flag.String("out", "osm.spaten", "")
	mappingPath := flag.String("mapping", "", "path to mapping file. default mapping will be applied if none is specified")
	flag.Parse()

	var conds []mapping.Condition
	if len(*mappingPath) == 0 {
		log.Println("No mapping specified. Using default tag mapping.")
		conds = mapping.Default
	}

	f, err := os.Open(*source)
	if err != nil {
		log.Fatal(err)
	}
	dec := gosmparse.NewDecoder(f)

	// First pass
	ec := NewElemCache()
	dh := dataHandler{
		conds: conds,
		ec:    ec,
	}
	log.Println("Starting 3 step parsing")
	log.Println("Reading data (1/3)...")
	err = dec.Parse(&dh)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(0, 0) // jumps to beginning of file
	if err != nil {
		log.Fatal(err)
	}

	// Second pass
	log.Println("Collecting nodes (2/3)...")
	err = dec.Parse(ec)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Third pass
	log.Println("Resolving dependent objects (3/3)...")
	rc := nodeCollector{
		ec: ec,
	}
	err = dec.Parse(&rc)
	if err != nil {
		log.Fatal(err)
	}

	var fc []spatial.Feature

	log.Println("Parsing completed.")

	log.Println("Collecting points...")
	for _, pt := range dh.nodes {
		props := map[string]interface{}{}
		for k, v := range pt.Tags {
			props[k] = v
		}
		fc = append(fc, spatial.Feature{
			Props:    props,
			Geometry: spatial.MustNewGeom(spatial.Point{float64(pt.Lon), float64(pt.Lat)}),
		})
	}

	log.Println("Assembling ways...")
	// TODO: auto-detect if linestring or polygon, based on tags
	for _, wy := range dh.ways {
		var (
			area  bool
			props = map[string]interface{}{}
			geom  interface{}
		)
		for k, v := range wy.Tags {
			if k == "__type" && v == "area" {
				area = true
				continue
			}
			props[k] = v
		}
		ln := ec.Line(wy.ID)
		if ln.Clockwise() {
			ln.Reverse()
		}
		if area {
			geom = spatial.Polygon{ln}
		} else {
			geom = ln
		}

		fc = append(fc, spatial.Feature{
			Props:    props,
			Geometry: spatial.MustNewGeom(geom),
		})
	}

	log.Println("Assembling relations...")
	for _, rl := range dh.rels {
		if v, ok := rl.Tags["type"]; !ok || v != "multipolygon" {
			continue
		}
		var poly spatial.Polygon

		for _, memb := range rl.Members {
			if memb.Role == "outer" || memb.Role == "inner" {
				ring := ec.Line(memb.ID)
				if (memb.Role == "outer" && ring.Clockwise()) || (memb.Role == "inner" && !ring.Clockwise()) {
					ring.Reverse()
				}
				poly = append(poly, ring)
			}
		}
	}

	log.Println("Writing out")
	of, err := os.Create(*outfile)
	if err != nil {
		log.Fatal(err)
	}
	var outCodec spaten.Codec
	err = outCodec.Encode(of, &spatial.FeatureCollection{fc, "4326"})
	if err != nil {
		log.Fatal(err)
	}
}
