package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/thomersch/grandine/lib/cugdf"
	"github.com/thomersch/grandine/lib/spatial"

	"github.com/thomersch/gosmparse"
)

type dataHandler struct {
	cond condition

	// dependent objects, which will be collected in the second pass
	depNodes    []int64
	depNodesMtx sync.Mutex

	nodes    []gosmparse.Node
	nodesMtx sync.Mutex
	ways     []gosmparse.Way
	waysMtx  sync.Mutex
	rels     []gosmparse.Relation
	relsMtx  sync.Mutex
}

func (d *dataHandler) ReadNode(n gosmparse.Node) {
	// TODO: make it possible to specify condition type (node/way/rel)
	if v, ok := n.Tags[d.cond.key]; ok {
		if len(d.cond.value) != 0 && d.cond.value != v {
			return
		}
		d.nodesMtx.Lock()
		d.nodes = append(d.nodes, n)
		d.nodesMtx.Unlock()
	}
}

func (d *dataHandler) ReadWay(w gosmparse.Way) {
	// TODO: make it possible to specify condition type (node/way/rel)
	if v, ok := w.Tags[d.cond.key]; ok {
		if len(d.cond.value) != 0 && d.cond.value != v {
			return
		}

		d.depNodesMtx.Lock()
		d.depNodes = append(d.depNodes, w.NodeIDs...)
		d.depNodesMtx.Unlock()

		d.waysMtx.Lock()
		d.ways = append(d.ways, w)
		d.waysMtx.Unlock()
	}
}

func (d *dataHandler) ReadRelation(r gosmparse.Relation) {
	// TODO: make it possible to specify condition type (node/way/rel)
	if v, ok := r.Tags[d.cond.key]; ok {
		if len(d.cond.value) != 0 && d.cond.value != v {
			return
		}

		d.relsMtx.Lock()
		d.rels = append(d.rels, r)
		d.relsMtx.Unlock()

		for _, memb := range r.Members {
			switch memb.Type {
			case gosmparse.NodeType:
				d.depNodesMtx.Lock()
				d.depNodes = append(d.depNodes, memb.ID)
				d.depNodesMtx.Unlock()
			}
		}
	}
}

type relCollector struct {
	nds    map[int64]spatial.Point
	ndsMtx sync.Mutex
}

func (d *relCollector) ReadNode(n gosmparse.Node) {
	d.ndsMtx.Lock()
	defer d.ndsMtx.Unlock()
	if _, ok := d.nds[n.ID]; ok {
		d.nds[n.ID] = spatial.Point{float64(n.Lon), float64(n.Lat)}
	}
}
func (d *relCollector) ReadWay(w gosmparse.Way)           {}
func (d *relCollector) ReadRelation(r gosmparse.Relation) {}

// const (
// 	typAny      = 0
// 	typNode     = 1
// 	typWay      = 2
// 	typRelation = 3
// )

type condition struct {
	key   string
	value string
}

func main() {
	cond := condition{"highway", "primary"}

	source := flag.String("src", "osm.pbf", "")
	outfile := flag.String("out", "osm.cudgf", "")
	flag.Parse()

	f, err := os.Open(*source)
	if err != nil {
		log.Fatal(err)
	}
	dec := gosmparse.NewDecoder(f)
	dh := dataHandler{
		cond: cond,
	}
	log.Println("Collecting data...")
	err = dec.Parse(&dh)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Seek(0, 0) // jumps to beginning of file
	if err != nil {
		log.Fatal(err)
	}

	ndmap := map[int64]spatial.Point{}
	for _, ndID := range dh.depNodes {
		ndmap[ndID] = spatial.Point{}
	}

	log.Println("Resolving dependent objects")
	rc := relCollector{
		nds: ndmap,
	}
	err = dec.Parse(&rc)
	if err != nil {
		log.Fatal(err)
	}

	var fc []spatial.Feature

	log.Println("Assembling ways")
	for _, wy := range dh.ways {
		props := map[string]interface{}{}
		for k, v := range wy.Tags {
			props[k] = v
		}

		var ls spatial.Line
		for _, nID := range wy.NodeIDs {
			ls = append(ls, rc.nds[nID])
		}

		fc = append(fc, spatial.Feature{
			Props:    props,
			Geometry: spatial.MustNewGeom(ls),
		})
	}

	log.Println("Writing out")
	of, err := os.Create(*outfile)
	if err != nil {
		log.Fatal(err)
	}
	err = cugdf.Marshal(fc, of)
	if err != nil {
		log.Fatal(err)
	}
}
