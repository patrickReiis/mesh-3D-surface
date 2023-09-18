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
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	svgData := getSvgData(r.Form)

	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, svgData)
}

func getSvgData(options map[string][]string) string {
	return "Hi"
}

func corner(i, j int) (svgX float64, svgY float64, ok bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := defaultXyrange * (float64(i)/defaultCells - 0.5)
	y := defaultXyrange * (float64(j)/defaultCells - 0.5)
	// Compute surface height z.
	z, ok := f(x, y)
	if ok == false {
		return 0, 0, false
	}
	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := defaultWidth/2 + (x-y)*cos30*defaultXyscale
	sy := defaultHeight/2 + (x+y)*sin30*defaultXyscale - z*defaultZscale
	return sx, sy, true
}

func f(x, y float64) (value float64, ok bool) {
	r := math.Hypot(x, y) // distance from (0,0)
	if math.IsNaN(r) == true {
		return 0, false
	}

	return math.Sin(r) / r, true
}
