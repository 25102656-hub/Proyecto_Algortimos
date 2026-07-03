package rtree

import "math"

type Rectangulo struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func NuevoRectangulo(x, y float64) Rectangulo {

	return Rectangulo{
		MinX: x,
		MinY: y,
		MaxX: x,
		MaxY: y,
	}

}

func (r Rectangulo) Area() float64 {

	return (r.MaxX - r.MinX) * (r.MaxY - r.MinY)

}

func (r Rectangulo) Unir(otro Rectangulo) Rectangulo {

	return Rectangulo{

		MinX: math.Min(r.MinX, otro.MinX),
		MinY: math.Min(r.MinY, otro.MinY),

		MaxX: math.Max(r.MaxX, otro.MaxX),
		MaxY: math.Max(r.MaxY, otro.MaxY),
	}

}

func (r Rectangulo) Crecimiento(otro Rectangulo) float64 {

	nuevo := r.Unir(otro)

	return nuevo.Area() - r.Area()

}

func (r Rectangulo) Intersecta(otro Rectangulo) bool {

	if r.MaxX < otro.MinX {
		return false
	}

	if r.MinX > otro.MaxX {
		return false
	}

	if r.MaxY < otro.MinY {
		return false
	}

	if r.MinY > otro.MaxY {
		return false
	}

	return true

}
