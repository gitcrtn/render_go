// render_go - main
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "fmt"
import "flag"
import "os"

func main(){
	var x,y,sample,subpixel,max_depth int
	var image_prefix string

	flag.IntVar(&x, "x", 640, "resolution width")
	flag.IntVar(&y, "y", 480, "resolution height")
	flag.IntVar(&sample, "s", 4, "sample")
	flag.IntVar(&subpixel, "p", 2, "subpixel")
	flag.IntVar(&max_depth, "d", 5, "max depth")
	flag.StringVar(&image_prefix, "i", "out", "prefix of image name")
	flag.Parse()
	
	dir, _ := os.Getwd()
	image_name := fmt.Sprintf("%s\\%s.%%03d.png", dir, image_prefix)
	complete_image_name := fmt.Sprintf("%s\\%s.complete.png", dir, image_prefix)
	
	fmt.Println("width:",x)
	fmt.Println("height:",y)
	fmt.Println("sample:",sample)
	fmt.Println("subpixel:",subpixel)
	fmt.Println("max depth:",max_depth)
	fmt.Println("image name:",image_name)
	
	setting := Setting{x,y,sample,subpixel,max_depth,image_name}	
	camera := Camera{Vector{50.0, 52.0, 220.0}, *Normalize(&Vector{0.0, -0.04, -30.0}), Vector{0.0, 1.0, 0.0}}
	screen := NewScreen(&camera,&setting)
	
	scene := NewScene()
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{1e5+1.0, 40.8, 81.6}, 1e5}, Material{Vector{0.75, 0.25, 0.25},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{-1e5+99.0, 40.8, 81.6}, 1e5}, Material{Vector{0.25, 0.25, 0.75},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{50.0, 40.8, 1e5}, 1e5}, Material{Vector{0.75, 0.75, 0.75},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{50.0, 40.8, -1e5+250}, 1e5}, Material{Vector{0.0, 0.0, 0.0},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{50.0, 1e5, 81.6}, 1e5}, Material{Vector{0.75, 0.75, 0.75},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{50.0, -1e5+81.6, 81.6}, 1e5}, Material{Vector{0.75, 0.75, 0.75},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{65.0, 20.0, 20.0}, 20.0}, Material{Vector{0.25, 0.75, 0.25},Vector{0, 0, 0},Reflection_Diffuse}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{27.0, 16.5, 47.0}, 16.5}, Material{Vector{0.99, 0.99, 0.99},Vector{0, 0, 0},Reflection_Specular}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{77.0, 16.5, 78.0}, 16.5}, Material{Vector{0.99, 0.99, 0.99},Vector{0, 0, 0},Reflection_Refraction}})
	scene.geometries = append(scene.geometries, Geometry{Sphere{Vector{50.0, 90.0, 81.6}, 15.0}, Material{Vector{0.0, 0.0, 0.0},Vector{36, 36, 36},Reflection_Diffuse}})
	
	im := NewImageBuffer(x,y)
	
	Render(&setting,&camera,screen,scene,im)
	im.Out(complete_image_name)	
}
