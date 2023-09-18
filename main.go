package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
)

const (
	defaultWidth, defaultHeight = 600, 320                          // canvas size in pixels
	defaultCells                = 100                               // number of grid cells
	defaultXYrange              = 30.0                              // axis ranges (-xyrange..+xyrange)
	defaultXYscale              = defaultWidth / 2 / defaultXYrange // pixels per x or y unit
	defaultZscale               = defaultHeight * 0.4               // pixels per z unit
	defaultAngle                = math.Pi / 6                       // angle of x, y axes (=30°)
	defaultColor                = "ff0000"                          // Red is the default color
)

var sin30, cos30 = math.Sin(defaultAngle), math.Cos(defaultAngle) // sin(30°), cos(30°)
var clientWidth, clientHeight = defaultWidth, defaultHeight
var clientColor = defaultColor
var clientXYscale = float64(clientWidth) / 2 / defaultXYrange
var clientZscale = float64(clientHeight) * 0.4

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

func getSvgData(options map[string][]string) (svgFile string) {

	if optValues, ok := options["height"]; ok && len(optValues) == 1 {
		if height, err := strconv.Atoi(optValues[0]); err == nil {
			clientHeight = height
			clientZscale = float64(clientHeight) * 0.4
		}
	}

	if optValues, ok := options["width"]; ok && len(optValues) == 1 {
		if width, err := strconv.Atoi(optValues[0]); err == nil {
			clientWidth = width
			clientXYscale = float64(clientWidth) / 2 / defaultXYrange
		}
	}

	if optValues, ok := options["color"]; ok && len(optValues) == 1 {
		color := []byte(optValues[0])
		if ok, _ := regexp.Match(`([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6})$`, color); ok == true {
			clientColor = string(color)
		}
	}

	svgFile += fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>", clientWidth, clientHeight)

	for i := 0; i < defaultCells; i++ {
		for j := 0; j < defaultCells; j++ {
			ax, ay, ok := corner(clientWidth, clientHeight, clientXYscale, clientZscale, i+1, j)
			if ok == false {
				continue
			}
			bx, by, ok := corner(clientWidth, clientHeight, clientXYscale, clientZscale, i, j)
			if ok == false {
				continue
			}
			cx, cy, ok := corner(clientWidth, clientHeight, clientXYscale, clientZscale, i, j+1)
			if ok == false {
				continue
			}
			dx, dy, ok := corner(clientWidth, clientHeight, clientXYscale, clientZscale, i+1, j+1)
			if ok == false {
				continue
			}

			svgFile += fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, clientColor)
		}
	}
	svgFile += fmt.Sprint("</svg>\n")

	return svgFile
}

func corner(width, height int, XYscale, Zscale float64, i, j int) (svgX float64, svgY float64, ok bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := defaultXYrange * (float64(i)/defaultCells - 0.5)
	y := defaultXYrange * (float64(j)/defaultCells - 0.5)
	// Compute surface height z.
	z, ok := f(x, y)
	if ok == false {
		return 0, 0, false
	}
	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := float64(width)/2 + (x-y)*cos30*float64(XYscale)
	sy := float64(height)/2 + (x+y)*sin30*float64(XYscale) - z*Zscale
	if math.IsNaN(sx) == true || math.IsNaN(sy) == true {
		return 0, 0, false
	}

	return sx, sy, true
}

func f(x, y float64) (value float64, ok bool) {
	r := math.Hypot(x, y) // distance from (0,0)
	if math.IsNaN(r) == true {
		return 0, false
	}

	return math.Sin(r) / r, true
}
