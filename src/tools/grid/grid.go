package grid

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/direction"
)

type RuneCell[V any] interface {
	*V
	Rune() rune
	Bytes() []byte

	Y() int
	X() int

	MoveTo(y, x int)
}

type Grid[V any, PtrV RuneCell[V]] struct {
	g [][]PtrV
	W int
	H int
}

func CreateGrid[V any, PtrV RuneCell[V]]() *Grid[V, PtrV] {
	return &Grid[V, PtrV]{
		g: make([][]PtrV, 0),
	}
}

func CreateFixedGrid[V any, PtrV RuneCell[V]](h, w int, fn func(y, x int) (*V, error)) (*Grid[V, PtrV], error) {
	g := &Grid[V, PtrV]{
		g: make([][]PtrV, h),
		W: w,
		H: h,
	}

	var err error

	for y := 0; y < h; y++ {
		g.g[y] = make([]PtrV, w)

		for x := 0; x < w; x++ {
			g.g[y][x], err = fn(y, x)

			if err != nil {
				return nil, err
			}
		}
	}

	return g, err
}

func (g *Grid[V, PtrV]) AddRow(line string, fn func(x int, r rune) (*V, error)) error {
	rs := []rune(line)
	w := len(rs)

	if len(g.g) == 0 {
		g.W = w
	} else if w != g.W {
		return fmt.Errorf("invalid line for grid width, got %d, want %d", w, g.W)
	}

	row := make([]PtrV, w)

	for x, r := range rs {
		if v, err := fn(x, r); err != nil {
			return err
		} else {
			row[x] = v
		}
	}

	g.H++
	g.g = append(g.g, row)

	return nil
}

func (g *Grid[V, PtrV]) Get(y, x int) *V {
	if y < 0 || y >= g.H || x < 0 || x >= g.W {
		return nil
	}

	return g.g[y][x]
}

func (g *Grid[V, PtrV]) Swap(a, b PtrV) error {
	if a.Y() < 0 || a.Y() >= g.H || a.X() < 0 || a.X() >= g.W {
		return errors.New("failed to swap: a is out of bounds")
	}

	if b.Y() < 0 || b.Y() >= g.H || b.X() < 0 || b.X() >= g.W {
		return errors.New("failed to swap: b is out of bounds")
	}

	ay, ax := a.Y(), a.X()
	by, bx := b.Y(), b.X()

	g.g[ay][ax] = b
	g.g[by][bx] = a

	a.MoveTo(by, bx)
	b.MoveTo(ay, ax)

	return nil
}

func (g *Grid[V, PtrV]) Iter() iter.Seq[*V] {
	return func(yield func(*V) bool) {
		for _, row := range g.g {
			for _, v := range row {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (g *Grid[V, PtrV]) Square(c PtrV) (t, r, b, l PtrV) {
	t = g.Get(c.Y()-1, c.X())
	r = g.Get(c.Y(), c.X()+1)
	b = g.Get(c.Y()+1, c.X())
	l = g.Get(c.Y(), c.X()-1)

	return
}

func (g *Grid[V, PtrV]) FprintRaw(w io.Writer) error {
	str := make([]byte, 0, g.H*g.W+g.H+100)

	for _, l := range g.g {
		for _, c := range l {
			str = append(str, c.Bytes()...)
		}

		str = append(str, '\n')
	}

	if _, err := w.Write(str); err != nil {
		return err
	}

	return nil
}

func (g *Grid[V, PtrV]) Fprint(w io.Writer) error {
	str := make([]byte, 0, g.H*g.W+g.H+100)
	str = append(str, []byte(fmt.Sprintf("\033[s\033[%dF", g.H))...)

	for _, l := range g.g {
		for _, c := range l {
			str = append(str, c.Bytes()...)
		}

		str = append(str, '\n')
	}

	str = append(str, []byte("\033[u")...)

	if _, err := w.Write(str); err != nil {
		return err
	}

	return nil
}

func (g *Grid[V, PtrV]) Lines() []string {
	strs := make([]string, 0, g.H)

	for _, l := range g.g {
		str := make([]byte, 0, g.W)
		for _, c := range l {
			str = append(str, c.Bytes()...)
		}

		strs = append(strs, string(str))
	}

	return strs
}

type BaseCell struct {
	y, x int
	r    rune
}

func CreateBaseCell(y, x int, r rune) BaseCell {
	return BaseCell{
		y: y,
		x: x,
		r: r,
	}
}

func (c *BaseCell) Rune() rune {
	return c.r
}

func (c *BaseCell) String() string {
	return fmt.Sprintf("(%d, %d: %s)", c.y, c.x, string(c.r))
}

func (c *BaseCell) Bytes() []byte {
	return []byte(string(c.r))
}

func (c *BaseCell) Y() int {
	return c.y
}

func (c *BaseCell) X() int {
	return c.x
}

func (c *BaseCell) MoveTo(y, x int) {
	c.y = y
	c.x = x
}

func (c *BaseCell) North() (y, x int) {
	return c.y - 1, c.x
}

func (c *BaseCell) East() (y, x int) {
	return c.y, c.x + 1
}

func (c *BaseCell) South() (y, x int) {
	return c.y + 1, c.x
}

func (c *BaseCell) West() (y, x int) {
	return c.y, c.x - 1
}

func (c *BaseCell) Dir(d direction.Dir) (y, x int) {
	switch d {
	case direction.Up:
		return c.North()
	case direction.Right:
		return c.East()
	case direction.Down:
		return c.South()
	case direction.Left:
		return c.West()
	}

	panic("invalid cell Direction")
}
