// render_go - vector
// Copyright (c) 2014 Keita Yamada
// 2014.09.20

package main

import "math"

type Vector struct{
	x float64
	y float64
	z float64
}

func (v *Vector) Set(vin Vector){
	v.x = vin.x
	v.y = vin.y
	v.z = vin.z	
}

func (v *Vector) Add(vin Vector){
	v.x += vin.x
	v.y += vin.y
	v.z += vin.z
}

func (v *Vector) Sub(vin Vector){
	v.x -= vin.x
	v.y -= vin.y
	v.z -= vin.z
}

func (v *Vector) Div(vin Vector){
	v.x /= vin.x
	v.y /= vin.y
	v.z /= vin.z
}

func (v *Vector) Mul(vin Vector){
	v.x *= vin.x
	v.y *= vin.y
	v.z *= vin.z
}

func (v *Vector) AddFloat(fin float64){
	v.x += fin
	v.y += fin
	v.z += fin
}

func (v *Vector) SubFloat(fin float64){
	v.x -= fin
	v.y -= fin
	v.z -= fin
}

func (v *Vector) DivFloat(fin float64){
	v.x /= fin
	v.y /= fin
	v.z /= fin
}

func (v *Vector) MulFloat(fin float64){
	v.x *= fin
	v.y *= fin
	v.z *= fin
}

func Cross(v1, v2 *Vector) *Vector{
	vout := new(Vector)
	vout.x = v1.y * v2.z - v1.z * v2.y
	vout.y = v1.z * v2.x - v1.x * v2.z
	vout.z = v1.x * v2.y - v1.y * v2.x
	return vout 
}

func Dot(v1, v2 *Vector) *float64{
	f := new(float64)
	*f = v1.x * v2.x + v1.y * v2.y + v1.z * v2.z
	return f
}

func Distance(v1, v2 *Vector) *float64{
	f := new(float64)
	*f = math.Sqrt(math.Pow((v1.x - v2.x),2) + math.Pow((v1.y - v2.y),2) + math.Pow((v1.z - v2.z),2))
	return f
}

func Normalize(v *Vector) *Vector{
	vout := new(Vector)
	*vout = *v
	d := Distance(vout,&Vector{0,0,0})
	if *d==0.0{		
		return vout
	}
	vout.DivFloat(*d)
	return vout
}
