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
	depth:=depthForZoonm(zoom)

	centerX,err := strconv.ParseFloat(r.URL.Query().Get("centerX"),64)
	if err!=nil {
		centerX=0
	}
	centerY,err := strconv.ParseFloat(r.URL.Query().Get("centerY"),64)
	if err!=nil {
		centerY=0
	}

	baseWidth:=300.0

	visibleWidth:= baseWidth/zoom

	box:= Box{
		centerX-visibleWidth/2,
		centerX+visibleWidth/2,
		centerY-visibleWidth/2,
		centerY+visibleWidth/2,
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
    http.ListenAndServe(":8083",nil)
}
