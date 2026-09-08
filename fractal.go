package main 
import (
	"math"
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

func generateSegment(a Point, b Point, depth int, box Box, result *[]Segment)  {
	if depth==0 {
		*result=append(*result, Segment{A:a,B:b})
		return
	}else{
		if segmentOutsideBox(a,b,box) {
			return
		}
		
		p1:= Point{
			a.X+(1.0/3.0)*(b.X-a.X), 
			a.Y+(1.0/3.0)*(b.Y-a.Y)}
		p2:= Point{
			b.X+(1.0/3.0)*(a.X-b.X),b.Y+(1.0/3.0)*(a.Y-b.Y)}
		
		dx:= p2.X-p1.X
		dy:= p2.Y-p1.Y
		angle:= 60*math.Pi/180
		dx_rot:=dx*math.Cos(angle)-dy*math.Sin(angle)
		dy_rot:=dx*math.Sin(angle)+dy*math.Cos(angle)

		p3:= Point{
			p1.X+dx_rot, p1.Y+dy_rot}
		
		generateSegment(a,p1,depth-1,box,result)
		generateSegment(p1,p3,depth-1,box,result)
		generateSegment(p3,p2,depth-1,box,result)
		generateSegment(p2,b,depth-1,box,result)
		
	}
}

func depthForZoonm(zoom float64) int {
	depth:=int(math.Log(zoom)/math.Log(3))
	depth+=2
	if depth>50 {
		depth=50
	} else if depth<1 {
		depth=1
	}
	return depth
}