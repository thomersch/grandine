# lib/mapping

[![GoDoc](https://godoc.org/github.com/thomersch/grandine?status.svg)](https://godoc.org/github.com/thomersch/grandine/lib/mapping)

This package provides facilities for parsing and applying mapping files. Mapping files allow to convert existing, freely tagged geodata into a stricter subset, as defined by a mapping file. This approach can be very useful if you want to rename keys, want to convert data types or just reduce the dataset without having to write code.

A mapping file is a yaml file (thus JSON is fully compatible, if you prefer to use that), which defines a list of projections, consisting of an input condition and one or more output statements.

## Condition

The key and value specified in `src` is the input condition. Any element that matches the given combination will be transformed using the `dest` parameters.

### Examples

* `{key: "highway", value: "primary"}` captures all elements with `highway=primary`
* `{key: "highway", value: "*"}` captures all elements that have any tag with the key `highway`.

## Output

In `dest` one can specify one or more output tags. Those can be either static for any matched element or have dynamic, typed values, which are inherited from the original element.

### Examples

* `- {key: "class", value: "railway"}` inserts a `class=highway` into all matched elements, as defined in `src`.
* `- {key: "v-max", value: "$maxspeed", type: int}` retrieves the maxspeed value from the source element, converts it to an integer and inserts it into `v-max`.

### Types

Any output element can have a `type` element, which defines the data type. Currently supported:

* `int`, casting to integer
* `string`, a series of bytes, no conversion, equivalent to not specifying any type
* no type, just interpreting as string

## Full Example

	- src:
	    key: highway
	    value: primary
	  dest:
	    - {key: "@layer", value: "transportation"}
	    - {key: "class", value: "$highway"}

	- src:
	    key: building
	    value: "*"
	  dest:
	    - {key: "@layer", value: "building"}
	    - {key: "@zoom:min", value: 14}

	- src:
	    key: railway
	    value: "*"
	  dest:
	    - {key: "@layer", value: "transportation"}
	    - {key: "class", value: "railway"}
	    - {key: "maxspeed", value: "$maxspeed", type: int}
