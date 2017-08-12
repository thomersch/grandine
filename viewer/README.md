# Debug Vector Tile Viewer

## Usage

	npm install # installs dependencies
	npm run server # starts an HTTP server on localhost:8080

Please note that the code assumes that tiles are placed in the `tiles/` directory. Consider placing a symlink.


### Generating Fonts

You can build the fonts using [fontnik](https://github.com/mapbox/node-fontnik), which can be installed via `npm install fontnik`. The executable will be located in `node_modules/fontnik/bin/build-glyphs`.
