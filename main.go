package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func snowFlakeHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	p1 := Point{0, 100}
	p2 := Point{-100, -50}
	p3 := Point{100, -50}

	zoom,err := strconv.ParseFloat(r.URL.Query().Get("zoom"),64)
	if err!=nil {
		zoom=1
	}
	depth := depthForZoonm(zoom)
	if dq := r.URL.Query().Get("depth"); dq != "" {
    if d, err := strconv.Atoi(dq); err == nil && d > 0 && d <= 50 {
        depth = d
    }
	}
	centerX,err := strconv.ParseFloat(r.URL.Query().Get("centerX"),64)
	if err!=nil {
		centerX=0
	}
	centerY,err := strconv.ParseFloat(r.URL.Query().Get("centerY"),64)
	if err!=nil {
		centerY=0
	}

	// Use actual screen dimensions if provided for accurate culling
	screenW := 1920.0
	screenH := 1080.0
	if sw, err := strconv.ParseFloat(r.URL.Query().Get("screenW"), 64); err == nil && sw > 0 {
		screenW = sw
	}
	if sh, err := strconv.ParseFloat(r.URL.Query().Get("screenH"), 64); err == nil && sh > 0 {
		screenH = sh
	}

	// Match the frontend's scale formula: scale = (H/400) * zoom
	// Visible world width  = screenW / scale = screenW * 400 / (screenH * zoom)
	// Visible world height = screenH / scale = 400 / zoom
	// Add 50% padding on each side to avoid edge-case culling during panning
	visibleW := (screenW * 400.0 / (screenH * zoom)) * 2.0
	visibleH := (400.0 / zoom) * 2.0

	box:= Box{
		centerX-visibleW/2,
		centerX+visibleW/2,
		centerY-visibleH/2,
		centerY+visibleH/2,
	}
	

	result := []Segment{}

	generateSegment(p1, p2, depth, box, &result)
	generateSegment(p2, p3, depth, box, &result)
	generateSegment(p3, p1, depth, box, &result)

	fmt.Println("result",len(result))

	json.NewEncoder(w).Encode(result)	

}

func main(){
	fmt.Println("fractals are cool")
 	http.HandleFunc("/snowflake",snowFlakeHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
    http.ListenAndServe(":8083",nil)
}
