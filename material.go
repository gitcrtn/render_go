// render_go - material
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "math"
import "math/rand"

type Reflection int

type Reflection_Algorithm interface{
	Radiance(*Scene, *Setting, *Ray, *rand.Rand, int, *Vertex, *Vector)
}

func (color *Vector) Diffuse(scene *Scene, setting *Setting, ray *Ray, rnd *rand.Rand, depth int, intersection *Vertex, orienting_normal *Vector){
	w := *orienting_normal
	u := new(Vector)
	switch{
		case math.Abs(w.x) > 0.0000009:
			u = Normalize(Cross(&Vector{0,1,0},&w))
		default:
			u = Normalize(Cross(&Vector{1,0,0},&w))
	}
	v := Cross(&w,u)
	r1 := 2.0 * math.Pi * rnd.Float64()
	r2 := rnd.Float64()
	rr2 := math.Sqrt(r2)
	u.MulFloat(math.Cos(r1) * rr2)
	v.MulFloat(math.Sin(r1) * rr2)
	w.MulFloat(math.Sqrt(1.0 - r2))
	u.Add(*v)
	u.Add(w)
	ray.pos.Set(intersection.pos)
	ray.dir.Set(*Normalize(u))
	*color = Radiance(scene,setting,ray,rnd,depth+1)
}

func (color *Vector) Specular(scene *Scene, setting *Setting, ray *Ray, rnd *rand.Rand, depth int, intersection *Vertex, orienting_normal *Vector){
	buf := intersection.normal
	buf.MulFloat(*Dot(&intersection.normal,&ray.dir) * 2.0)
	ray.pos = intersection.pos
	ray.dir.Sub(buf)
	*color = Radiance(scene,setting,ray,rnd,depth+1)
}

func (color *Vector) Refraction(scene *Scene, setting *Setting, ray *Ray, rnd *rand.Rand, depth int, intersection *Vertex, orienting_normal *Vector){
	into := *Dot(orienting_normal,&intersection.normal) > 0.0
	default_refraction := 1.0
	object_refraction := 1.5
	ray_refraction := *new(float64)
	switch{
		case into:
			ray_refraction = default_refraction / object_refraction
		default:
			ray_refraction = object_refraction / default_refraction
	}
	incident_dot := *Dot(&ray.dir,orienting_normal)
	critical_factor := 1.0 - math.Pow(ray_refraction,2) * (1.0 - math.Pow(incident_dot,2))
	
	reflection_ray := new(Ray)
	reflection_ray.pos = intersection.pos
	reflection_ray.dir = intersection.normal
	reflection_ray.dir.MulFloat(*Dot(&intersection.normal,&ray.dir) * -2.0)
	reflection_ray.dir.Add(ray.dir)
	
	// total reflection
	if critical_factor < 0.0{
		*color = Radiance(scene,setting,reflection_ray,rnd,depth+1)
		return
	}
	
	refraction_ray := new(Ray)
	refraction_ray.pos = intersection.pos
	bufDir := intersection.normal
	switch{
		case into:
			bufDir.MulFloat(-incident_dot * ray_refraction)
		default:
			bufDir.MulFloat(incident_dot * ray_refraction)			
	}
	bufDir.Add(ray.dir)
	bufDir.AddFloat(ray_refraction + math.Sqrt(critical_factor))
	refraction_ray.dir = *Normalize(&bufDir)
	
	p := object_refraction - default_refraction
	q := object_refraction + default_refraction
	vertical_incidence_factor := math.Pow(p,2) / math.Pow(q,2)
	r := *new(float64)
	switch{
		case into:
			r = 1.0 + incident_dot
		default:
			tmp := *orienting_normal
			tmp.MulFloat(-1.0)
			r = 1.0 - *Dot(&refraction_ray.dir,&tmp)
	}
	fresnel_incidence_factor := vertical_incidence_factor + (1.0 - vertical_incidence_factor) * math.Pow(r,5)
	radiance_scale := math.Pow(ray_refraction,2)
	refraction_factor := (1.0 - fresnel_incidence_factor) * radiance_scale
	
	probability := 0.75 + fresnel_incidence_factor
	
	switch{
		case depth > 2:
			switch{
				case rnd.Float64() < probability:
					*color = Radiance(scene,setting,reflection_ray,rnd,depth+1)
					color.MulFloat(fresnel_incidence_factor)
					return
				default:
					*color = Radiance(scene,setting,reflection_ray,rnd,depth+1)
					color.MulFloat(refraction_factor)
					return
			}
		default:
			*color = Radiance(scene,setting,reflection_ray,rnd,depth+1)
			color.MulFloat(fresnel_incidence_factor)
			refr := Radiance(scene,setting,reflection_ray,rnd,depth+1)
			refr.MulFloat(refraction_factor)
			color.Add(refr)
			return
	}	
}

const (
	Reflection_Diffuse Reflection = iota
	Reflection_Specular
	Reflection_Refraction
)

type Material struct{
	color Vector
	emission Vector
	reflection Reflection
}
