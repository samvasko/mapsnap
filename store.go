package main

import (
	"errors"
	// "fmt"
	"strconv"
	"strings"
)

type coord uint

type point struct {
	x, y coord
}

type matrix struct {
	/*
	   TL -- TR
	   |      |
	   BL -- BR
	*/
	TL, TR, BL, BR point
}

func (m *matrix) size() coord {
	return (m.TR.x - m.BL.x) * (m.TR.y - m.BL.y)
}

func (m *matrix) width() coord {
	return m.TR.x - m.TL.x
}

func (m *matrix) height() coord {
	return m.TL.y - m.BL.y
}

func digestPoint(strpoint string) point {
	splitstrpoint := strings.Split(strpoint, ",")
	var coords []coord
	for _, c := range splitstrpoint {
		nc, err := strconv.ParseUint(c, 10, 64)
		if err != nil {
			bailout(errors.New("Incorrect point definition: Define points like this '10,10 20,34'"))
		}
		coords = append(coords, coord(nc))
	}

	return point{coords[0], coords[1]}
}

func CreateMatrix(points [2]string) matrix {

	// Create ambigous representation of two points
	a := digestPoint(points[0])
	b := digestPoint(points[1])

	// Now find their real position
	TL := point{smaller(a.x, b.x), bigger(a.y, b.y)}
	BR := point{bigger(a.x, b.x), smaller(a.y, b.y)}

	TR := point{BR.x, TL.y}
	BL := point{TL.x, BR.y}

	return matrix{TL, TR, BL, BR}
}

func bigger(x1, x2 coord) coord {
	if x1 > x2 {
		return x1
	} else if x1 < x2 {
		return x2
	} else {
		bailout(errors.New("Points do not form rectangle"))
		return 0
	}
}

func smaller(x1, x2 coord) coord {
	if x1 < x2 {
		return x1
	} else if x1 > x2 {
		return x2
	} else {
		bailout(errors.New("Points do not form rectangle"))
		return 0
	}
}
