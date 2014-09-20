// render_go - render
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "time"
import "math/rand"
import "fmt"
import "math"

type Camera struct{
	pos Vector
	dir Vector
	up  Vector
}

type Screen struct{
	dist   float64
	width  float64
	height float64
	x	   Vector
	y	   Vector
	center Vector
}

// default-args: scale=30.0, dist=40.0
func NewScreen(camera *Camera, setting *Setting, args ...float64) *Screen{
	s := new(Screen)
	s.dist = 40.0
	scale := 30.0
	if len(args) > 0{
		scale = args[0]
	}
	if len(args) > 1{
		s.dist = args[1]
	}
	s.width = scale * float64(setting.width) / float64(setting.height)
	s.height = scale
	s.x.Set(*Normalize(Cross(&camera.dir,&camera.up)))
	s.x.MulFloat(s.width)
	s.y.Set(*Normalize(Cross(&camera.dir,&s.x)))
	s.y.MulFloat(s.height)		
	s.center.Set(camera.dir)
	s.center.MulFloat(s.dist)
	s.center.Add(camera.pos)
	return s
}

type Setting struct{
	width   int
	height  int
	samples int
	supersamples int
	max_depth int
	imagePrefix string
}

type Ray struct{
	pos Vector
	dir Vector
}

func NewRay(camera *Camera, screen *Screen, setting *Setting, x, y, supersample_x, supersample_y int) *Ray {
	ray := new(Ray)
	rate := 1.0 / float64(setting.supersamples)
	rx := float64(supersample_x) * rate + rate / 2.0
	ry := float64(supersample_y) * rate + rate / 2.0
	var end Vector
	var buf Vector
	end.Set(screen.center)
	buf.Set(screen.x)
	buf.MulFloat((rx + float64(x)) / float64(setting.width) - 0.5)
	end.Add(buf)
	buf.Set(screen.y)
	buf.MulFloat((ry + float64(y)) / float64(setting.height) - 0.5)
	end.Add(buf)
	end.Sub(camera.pos)
	ray.dir.Set(*Normalize(&end))
	ray.pos.Set(camera.pos)
	return ray
}

func Radiance(scene *Scene, setting *Setting, ray *Ray, rnd *rand.Rand, depth int) Vector{
	intersection, obj_id, result := scene.Intersect(ray)
	
	if !result{
		return scene.bgColor
	}
		
	geom := scene.geometries[obj_id]

	orienting_normal := intersection.normal
	if *Dot(&intersection.normal,&ray.dir) > 0.0{
		orienting_normal.MulFloat(-1.0)	
	}
	
	if depth > setting.max_depth{
		return geom.material.emission
	}
	
	reflect_ratio := float64(setting.max_depth - depth) / float64(setting.max_depth)

	var inc_rad Vector
	switch{
		case geom.material.reflection == Reflection_Diffuse:
			inc_rad.Diffuse(scene,setting,ray,rnd,depth,intersection,&orienting_normal)
		case geom.material.reflection == Reflection_Specular:
			inc_rad.Specular(scene,setting,ray,rnd,depth,intersection,&orienting_normal)		
		case geom.material.reflection == Reflection_Refraction:
			inc_rad.Refraction(scene,setting,ray,rnd,depth,intersection,&orienting_normal)
	}	
	weight := geom.material.color
	weight.MulFloat(reflect_ratio)
	
	inc_rad.Mul(weight)
	inc_rad.Add(geom.material.emission)
	
	return inc_rad
}

func Render(setting *Setting, camera *Camera, screen *Screen, scene *Scene, image *ImageBuffer){
	
	start_time := time.Now()
	current_time := time.Now()
	var now_duration time.Duration
	time_count := 1.0
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	
	for y:=0; y < setting.height; y++{
		rnd.Seed(time.Now().Unix())
		fmt.Println(fmt.Sprintf("Rendering (y = \"%d\") %.02f%%", y, 100.0 * float64(y) / float64(setting.height - 1)))
		
		for x:=0; x < setting.width; x++{
			image.SetColor(x,y,Vector{0,0,0})
			
			for sy:=0; sy < setting.supersamples; sy++{
			
				for sx:=0; sx < setting.supersamples; sx++{
					acm_rad := Vector{0,0,0}
					buf_rad := Vector{0,0,0}
					
					for s:=0; s < setting.samples; s++{
						ray := NewRay(camera,screen,setting,x,y,sx,sy)			
						
						buf_rad = Radiance(scene,setting,ray,rnd,0)
						buf_rad.DivFloat(float64(setting.samples) * math.Pow(float64(setting.supersamples),2))						
						acm_rad.Add(buf_rad)
						image.AddColor(x,y,acm_rad)
						
						current_time = time.Now()
						now_duration = current_time.Sub(start_time)
						if now_duration.Minutes() > time_count  {
							fmt.Println(time_count, "minute(s)")
							fmt.Println("image output...")
							image.Out( fmt.Sprintf(setting.imagePrefix,int(time_count)) )
							time_count += 1.0
						}
					}				
				}			
			}
		}
	}	
}
