package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"slices"
)

func SaveAsGray16(d [][]uint16, name string) error {
	w := len(d)
	img := image.NewGray16(image.Rectangle{Max: image.Point{w, w}})
	for i := range d {
		for j := range d {
			img.SetGray16(i, j, color.Gray16{d[i][j]})
		}
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	return png.Encode(f, img)
}

type earthColorMapPoint struct {
	v, r, g, b float32
}

// earth color map obtained from
// https://janert.me/blog/2022/the-diamond-square-algorithm-for-terrain-generation/
var earthColorMap = []earthColorMapPoint{
	{0.000000, 0.44314, 0.67059, 0.84706},
	{0.037037, 0.47451, 0.69804, 0.87059},
	{0.074074, 0.51765, 0.72549, 0.89020},
	{0.111111, 0.55294, 0.75686, 0.91765},
	{0.148148, 0.58824, 0.78824, 0.94118},
	{0.185185, 0.63137, 0.82353, 0.96863},
	{0.222222, 0.67451, 0.85882, 0.98431},
	{0.259259, 0.72549, 0.89020, 1.00000},
	{0.296296, 0.77647, 0.92549, 1.00000},
	{0.333333, 0.84706, 0.94902, 0.99608},
	{0.333333, 0.67451, 0.81569, 0.64706},
	{0.370370, 0.58039, 0.74902, 0.54510},
	{0.407407, 0.65882, 0.77647, 0.56078},
	{0.444444, 0.74118, 0.80000, 0.58824},
	{0.481481, 0.81961, 0.84314, 0.67059},
	{0.518519, 0.88235, 0.89412, 0.70980},
	{0.555556, 0.93725, 0.92157, 0.75294},
	{0.592593, 0.90980, 0.88235, 0.71373},
	{0.629630, 0.87059, 0.83922, 0.63922},
	{0.666667, 0.82745, 0.79216, 0.61569},
	{0.703704, 0.79216, 0.72549, 0.50980},
	{0.740741, 0.76471, 0.65490, 0.41961},
	{0.777778, 0.72549, 0.59608, 0.35294},
	{0.814815, 0.66667, 0.52941, 0.32549},
	{0.851852, 0.67451, 0.60392, 0.48627},
	{0.888889, 0.72941, 0.68235, 0.60392},
	{0.925926, 0.79216, 0.76471, 0.72157},
	{0.962963, 0.87843, 0.87059, 0.84706},
	{1.000000, 0.96078, 0.95686, 0.94902},
}

func getInterColor(v uint16) color.RGBA {
	// convert v into a float32 in the range [0..1]
	x := float32(v) / float32(0xFFFF)
	i, ok := slices.BinarySearchFunc(earthColorMap, x, func(c earthColorMapPoint, x float32) int {
		if x < c.v {
			return 1
		} else if x > c.v {
			return -1
		}
		return 0
	})
	if ok {
		return color.RGBA{
			R: uint8(255 * earthColorMap[i].r),
			G: uint8(255 * earthColorMap[i].g),
			B: uint8(255 * earthColorMap[i].b),
			A: 255,
		}
	}
	return color.RGBA{
		R: interpolate(earthColorMap[i-1].r, earthColorMap[i].r, x),
		G: interpolate(earthColorMap[i-1].g, earthColorMap[i].g, x),
		B: interpolate(earthColorMap[i-1].b, earthColorMap[i].b, x),
		A: 255,
	}
}

func interpolate(a, b, x float32) uint8 {
	return uint8(255 * (a + x*(b-a)))
}

func SaveAsColor(d [][]uint16, name string) error {
	w := len(d)
	img := image.NewRGBA(image.Rectangle{Max: image.Point{w, w}})
	for i := range d {
		for j := range d {
			img.SetRGBA(i, j, getInterColor(d[i][j]))
		}
	}
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	return png.Encode(f, img)
}
