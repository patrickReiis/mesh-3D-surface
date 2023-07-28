package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

const (
	defaultWidth, defaultHeight = 600, 320                          // canvas size in pixels
	defaultCells                = 100                               // number of grid cells
	defaultXyrange              = 30.0                              // axis ranges (-xyrange..+xyrange)
	defaultXyscale              = defaultWidth / 2 / defaultXyrange // pixels per x or y unit
	defaultZscale               = defaultHeight * 0.4               // pixels per z unit
	defaultAngle                = math.Pi / 6                       // angle of x, y axes (=30°)
	defaultPeaksColor           = "#ff0000"                         // Red color is the default color for peaks
	defaultValleysColor         = "#0000ff"                         // Blue color is the default color for valleys
)

var sin30, cos30 = math.Sin(defaultAngle), math.Cos(defaultAngle) // sin(30°), cos(30°)

const port = "8080"

func main() {

	http.HandleFunc("/", handleDisplaySurface)

	log.Printf("Running on http://localhost:%s", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleDisplaySurface(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Display Surface")
}
