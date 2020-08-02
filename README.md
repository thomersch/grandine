# Grandine

[![GoDoc](https://godoc.org/github.com/thomersch/grandine?status.svg)](https://godoc.org/github.com/thomersch/grandine) [![Build Status](https://travis-ci.org/thomersch/grandine.svg?branch=master)](https://travis-ci.org/thomersch/grandine) 

This repository contains libraries and command line tools for working with geospatial data. It aims to streamline vector tile generation and provides tooling for standardized geo data serialization.

Initial work has been funded by the [Prototype Fund](https://prototypefund.de), powered by Open Knowledge Foundation Germany and the German Ministry for Research and Education.

![Prototype Fund](https://files.skowron.eu/grandine/logo-prototype.svg) ![Bundesministerium für Bildung und Forschung](https://files.skowron.eu/grandine/logo-bmbf.svg) ![Open Knowledge Foundation Deutschland](https://files.skowron.eu/grandine/logo-okfn.svg)

## Spaten File Format

If you are looking for documentation on the Spaten geo spatial file format, look [here](https://thomas.skowron.eu/spaten/). The Go implementation resides inside [lib/spaten](https://github.com/thomersch/grandine/tree/master/lib/spaten).

## Requirements

* Go ≥ 1.11, Go 1.14 recommended
* GEOS
* LevelDB

## Quickstart

If you have built a Go project before, you probably already know what to do. If not:

* Make sure you have Go installed. Preferably version 1.9 or higher. (If you can, use 1.10 or newer as it allows for some impactful performance improvements)
* Execute `go get -u github.com/thomersch/grandine`. This will checkout a current version of the code into `~/go/src/github.com/thomersch/grandine`
* Go to the checkout directory. Execute `make build`, this will put all executables into the `bin` directory.
* All the executables can be called with the `-help` flag which will print out basic usage info.

## Tips

### How to concatinate/merge multiple spaten files into one

	grandine-converter -in fileA,fileB,fileC -out outfile

Alternatively, if you want to use a pipe:

	grandine-converter -in fileA,fileB,fileC | your-app-here

### How to render a tile set from a spaten file

	grandine-tiler -in some_geodata.spaten -zoom 9,10,11 -out tiles/

By default, all data will be on the `default` layer.

## Structure

* `fileformat` contains a draft spec for a new geo data format that aims to be flexible, with a big focus on being very fast to serialize/deserialize.
* In `lib` you'll find a few Go libraries that provide a few primitives for handling spatial data:
	* `lib/spatial` contains functionality for handling points/lines/polygons and basic transformation operations. If you miss functionality, feel free to send a Pull Request, it would be greatly appreciated.
	* `lib/mvt` contains code for serializing Mapbox Vector Tiles.
* There are a few command line tools in `cmd`:
	* `converter` is a helper tool for converting and concatenating geo data files
	* `spatialize` converts OpenStreetMap data into a Spaten data file as defined in `fileformat`
	* `tiler` generates vector tiles from spatial data
