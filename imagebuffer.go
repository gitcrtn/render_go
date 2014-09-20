// render_go - imagebuffer
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "os"
import "image"
import "image/color"
import "image/png"
import "math"

import "fmt"

func clamp(f float64) float64{
	if f < 0.0{
		return 0.0
	}
	if f > 1.0 {
		return 1.0
	}
	return f
}

func f2c(f float64) uint8{
	return uint8(math.Pow(clamp(f), 1.0/2.2) * 255.0 + 0.5)
}

type ImageBuffer struct{
	x int
	y int
	instance []float64
	buf *image.NRGBA
}

func NewImageBuffer(x,y int) *ImageBuffer {
	i := new(ImageBuffer)
	i.x = x
	i.y = y
	i.instance = make([]float64,x*y*3)
	i.buf = image.NewNRGBA(image.Rect(0, 0, x, y))
	return i
}

func (i *ImageBuffer)SetColor(x,y int, c Vector){
	index := (( i.y - y - 1 ) * i.x + x ) * 3
	i.instance[index+0] = c.x
	i.instance[index+1] = c.y
	i.instance[index+2] = c.z
}

func (i *ImageBuffer)AddColor(x,y int, c Vector){
	index := (( i.y - y - 1 ) * i.x + x ) * 3
	i.instance[index+0] += c.x
	i.instance[index+1] += c.y
	i.instance[index+2] += c.z
}

func (i *ImageBuffer) test(){
	index := 0
	for y:=0; y < i.y; y++{
		for x:=0; x < i.x; x++{
			index = (( i.y - y - 1 ) * i.x + x ) * 3
			i.instance[index+0] = 1.0
		}
	}
}

func (i *ImageBuffer) Out(image_name string){
	index := 0
	for y:=0; y < i.y; y++{
		for x:=0; x < i.x; x++{
			index = (( i.y - y - 1 ) * i.x + x ) * 3
			i.buf.Set(x,y,color.NRGBA{f2c(i.instance[index+0]),
									f2c(i.instance[index+1]),
									f2c(i.instance[index+2]),255})
		}
	}
	imgfile,_:=os.Create(fmt.Sprintf(image_name)) 
	png.Encode(imgfile,i.buf)
}
