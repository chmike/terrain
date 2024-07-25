package main

import (
	"fmt"
	"math/rand/v2"
)

// MakeTerrain returns a square terrain height map of size 2^k+1. The parameter ds
// defines the terrain roughness and must be in the range between 0 and 1.
func MakeTerrain(k int, ds float64) ([][]uint16, error) {
	if k >= 32 {
		return nil, fmt.Errorf("expect size smaller than 2^32+1, got 2^%d+1", k)
	}
	if ds <= 0 || ds >= 1 {
		return nil, fmt.Errorf("expect roughness in the range 0 to 1, got %f", ds)
	}
	n := 1<<k + 1
	d := make([][]uint16, n)
	for i := range d {
		d[i] = make([]uint16, n)
	}

	d[0][0] = uint16(rand.Uint32() / 16)
	d[n-1][0] = uint16(rand.Uint32() / 16)
	d[0][n-1] = uint16(rand.Uint32() / 16)
	d[n-1][n-1] = uint16(rand.Uint32() / 16)

	for w, s := n-1, 1.; w > 1; w, s = w/2, s*ds {
		singleDiamondSquareStep(d, w, s)
	}
	return d, nil
}

type offset struct {
	di, dj int
}

// singleDiamondSquareStep computes the height using width w, and roughness ds.
func singleDiamondSquareStep(d [][]uint16, w int, ds float64) {
	n := len(d)
	v := w / 2
	diamond := [4]offset{{-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	square := [4]offset{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// diamond step
	for i := v; i < n; i += w {
		for j := v; j < n; j += w {
			average(d, i, j, v, diamond, ds)
		}
	}

	// square step rows
	for i := v; i < n; i += w {
		for j := 0; j < n; j += w {
			average(d, i, j, v, square, ds)
		}
	}

	// square step cols
	for i := 0; i < n; i += w {
		for j := v; j < n; j += w {
			average(d, i, j, v, square, ds)
		}
	}
}

// average returns the average of the points at given offsets.
func average(d [][]uint16, i, j, v int, offsets [4]offset, ds float64) {
	var r, k int
	n := len(d)
	for _, o := range offsets {
		u := i + o.di*v
		v := j + o.dj*v
		if u >= 0 && u < n && v >= 0 && v < n {
			r += int(d[u][v])
			k++
		}
	}
	d[i][j] = randomize(r/k, ds)
}

func randomize(v int, ds float64) uint16 {
	v += int(ds * float64(rand.IntN(0x10000)-0x8000))
	if v < 0 {
		v = 0
	} else if v > 0xFFFF {
		v = 0xFFFF
	}
	return uint16(v)
}
