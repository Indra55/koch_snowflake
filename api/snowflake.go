package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

type Point struct {
	X float64
	Y float64
}

type Box struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

type Segment struct {
	A Point
	B Point
}

func segmentOutsideBox(a Point, b Point, box Box) bool {
	dx := b.X - a.X
	dy := b.Y - a.Y
	pad := math.Sqrt(dx*dx+dy*dy) * 0.5

	segMinX := math.Min(a.X, b.X) - pad
	segMaxX := math.Max(a.X, b.X) + pad
	segMinY := math.Min(a.Y, b.Y) - pad
	segMaxY := math.Max(a.Y, b.Y) + pad

	if segMaxX < box.MinX || segMinX > box.MaxX || segMaxY < box.MinY || segMinY > box.MaxY {
		return true
	}
	return false
}

func generateSegment(a Point, b Point, depth int, box Box, result *[]Segment) {
	if depth == 0 {
		*result = append(*result, Segment{A: a, B: b})
		return
	} else {
		if segmentOutsideBox(a, b, box) {
			return
		}

		p1 := Point{
			X: a.X + (1.0/3.0)*(b.X-a.X),
			Y: a.Y + (1.0/3.0)*(b.Y-a.Y),
		}
		p2 := Point{
			X: b.X + (1.0/3.0)*(a.X-b.X),
			Y: b.Y + (1.0/3.0)*(a.Y-b.Y),
		}

		dx := p2.X - p1.X
		dy := p2.Y - p1.Y
		angle := 60 * math.Pi / 180
		dx_rot := dx*math.Cos(angle) - dy*math.Sin(angle)
		dy_rot := dx*math.Sin(angle) + dy*math.Cos(angle)

		p3 := Point{
			X: p1.X + dx_rot,
			Y: p1.Y + dy_rot,
		}

		generateSegment(a, p1, depth-1, box, result)
		generateSegment(p1, p3, depth-1, box, result)
		generateSegment(p3, p2, depth-1, box, result)
		generateSegment(p2, b, depth-1, box, result)
	}
}

func depthForZoonm(zoom float64) int {
	depth := int(math.Log(zoom) / math.Log(3))
	depth += 2
	if depth > 50 {
		depth = 50
	} else if depth < 1 {
		depth = 1
	}
	return depth
}

func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	p1 := Point{X: 0, Y: 100}
	p2 := Point{X: -100, Y: -50}
	p3 := Point{X: 100, Y: -50}

	zoom, err := strconv.ParseFloat(r.URL.Query().Get("zoom"), 64)
	if err != nil {
		zoom = 1
	}
	depth := depthForZoonm(zoom)
	if dq := r.URL.Query().Get("depth"); dq != "" {
		if d, err := strconv.Atoi(dq); err == nil && d > 0 && d <= 50 {
			depth = d
		}
	}
	centerX, err := strconv.ParseFloat(r.URL.Query().Get("centerX"), 64)
	if err != nil {
		centerX = 0
	}
	centerY, err := strconv.ParseFloat(r.URL.Query().Get("centerY"), 64)
	if err != nil {
		centerY = 0
	}

	screenW := 1920.0
	screenH := 1080.0
	if sw, err := strconv.ParseFloat(r.URL.Query().Get("screenW"), 64); err == nil && sw > 0 {
		screenW = sw
	}
	if sh, err := strconv.ParseFloat(r.URL.Query().Get("screenH"), 64); err == nil && sh > 0 {
		screenH = sh
	}

	visibleW := (screenW * 400.0 / (screenH * zoom)) * 2.0
	visibleH := (400.0 / zoom) * 2.0

	box := Box{
		MinX: centerX - visibleW/2,
		MaxX: centerX + visibleW/2,
		MinY: centerY - visibleH/2,
		MaxY: centerY + visibleH/2,
	}

	result := []Segment{}

	generateSegment(p1, p2, depth, box, &result)
	generateSegment(p2, p3, depth, box, &result)
	generateSegment(p3, p1, depth, box, &result)

	json.NewEncoder(w).Encode(result)
}
