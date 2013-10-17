package main

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

func CreateMatrix(a, b point) matrix {
	TL := point{smaller(a.x, b.x), bigger(a.y, b.y)}
	BR := point{bigger(a.x, b.x), smaller(a.y, b.y)}

	TR := point{BR.x, TL.y}
	BL := point{TL.x, BR.y}

	return matrix{TL, TR, BL, BR}
}

func digestPoint(strpoint string) point {
	splitstrpoint := strings.Split(strpoint, ",")
	var coords []coord
	for _, c := range splitstrpoint {
		nc, err := strconv.ParseUint(c, 10, 64)
		if err != nil {
			panic("Failied to parse point coordinates")
		}
		coords = append(coords, coord(nc))
	}

	return point{coords[0], coords[1]}
}

func bigger(x1, x2 coord) coord {
	if x1 > x2 {
		return x1
	} else if x1 < x2 {
		return x2
	} else {
		panic("Points do not form rectangle")
	}
}

func smaller(x1, x2 coord) coord {
	if x1 < x2 {
		return x1
	} else if x1 > x2 {
		return x2
	} else {
		panic("Points do not form rectangle")
	}
}
