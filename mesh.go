// render_go - mesh
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "math"

type Vertex struct{
	pos Vector
	normal Vector
}

type BaseMesh interface{
	Intersect(*Ray)(*Vertex, bool)
}

type Sphere struct{
	pos Vector
	radius float64
}

func (s *Sphere) Intersect(ray *Ray)(*Vertex, bool){
	intersection := new(Vertex)
	result := false
	v := ray.pos
	v.Sub(s.pos)
	B := *Dot(&ray.dir,&v) * 2.0
	C := *Dot(&v,&v) - math.Pow(s.radius,2)
	discr := math.Pow(B,2) - 4.0 * C
	if discr >= 0.0{
		sqroot := math.Sqrt(discr)
		t := (-B - sqroot) * 0.5
		if t < 0.0{
			t = (-B + sqroot) * 0.5
		}
		if t >= 0.0{
			result = true
			intersection.pos = ray.dir
			intersection.pos.MulFloat(t)
			intersection.pos.Add(ray.pos)
			pos := intersection.pos
			pos.Sub(s.pos)
			intersection.normal = *Normalize(&pos)
		}
	}
	return intersection, result
}
