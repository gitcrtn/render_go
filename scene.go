// render_go - scene
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "math"

type Geometry struct{
	//mesh     BaseMesh
	mesh Sphere
	material Material
}

type Scene struct{
	geometries []Geometry
	bgColor Vector
}

func NewScene(bgColor ...Vector) *Scene{
	s := new(Scene)
	switch{
		case len(bgColor) > 0:
			s.bgColor = bgColor[0]
		default:
			s.bgColor = Vector{0,0,0}
	}
	return s
}

func (s *Scene)Intersect(ray *Ray)(*Vertex,int,bool){
	dist := math.Inf(1)
	current_dist := dist
	result := false
	current_result := false
	intersection := new(Vertex)
	hit_point := new(Vertex)
	obj_id := *new(int)
	for i:=0; i < len(s.geometries); i++{
		hit_point, current_result = s.geometries[i].mesh.Intersect(ray)
		if current_result{
			dist = *Distance(&hit_point.pos,&ray.pos)
			if dist < current_dist{
				intersection = hit_point
				current_dist = dist
				obj_id = i
				result = true
			}
		}
	}
	return intersection, obj_id, result
}