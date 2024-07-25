package main

func main() {
	k := 10
	ds := 0.68
	d, err := MakeTerrain(k, ds)
	if err != nil {
		panic(err)
	}

	if err := SaveAsGray16(d, "gray16.png"); err != nil {
		panic(err)
	}

	if err := SaveAsColor(d, "color.png"); err != nil {
		panic(err)
	}
}
