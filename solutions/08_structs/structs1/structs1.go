package structs1

type Rectangle struct {
	Width, Height float64
}

func Area(r Rectangle) float64 {
	return r.Width * r.Height
}
