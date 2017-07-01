# Grandine

[![GoDoc](https://godoc.org/github.com/thomersch/grandine?status.svg)](https://godoc.org/github.com/thomersch/grandine) [![Build Status](https://travis-ci.org/thomersch/grandine.svg?branch=master)](https://travis-ci.org/thomersch/grandine) 

This repository contains libraries and command line tools for working with geospatial data. It aims to streamline vector tile generation and provides tooling for standardized geo data serialization.

The work is partly funded by the [Prototype Fund](https://prototypefund.de), powered by Open Knowledge Foundation Germany.

![](https://files.skowron.eu/grandine/logo-prototype.svg) ![](https://files.skowron.eu/grandine/logo-bmbf.svg) ![](https://files.skowron.eu/grandine/logo-okfn.svg)

## Structure

* `fileformat` contains a draft spec for a new geo data format that aims to be flexible, with a big focus on being very fast to serialize/deserialize.
* In `lib` you'll find a few Go libraries that provide a few primitives for handling spatial data:
	* `lib/spatial` contains functionality for handling points/lines/polygons and basic transformation operations. If you miss functionality, feel free to send a Pull Request, it would be greatly appreciated.
	* `lib/mvt` contains code for serializing Mapbox Vector Tiles.
* There are a few command line tools in `cmd`:
	* `spatialize` converts OpenStreetMap data into a spatial data format as defined in `fileformat`
	* `tiler` generates vector tiles from spatial data
