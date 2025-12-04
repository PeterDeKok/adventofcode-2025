package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/direction"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/grid"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
	"time"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	g, err := parseInput(rd)
	if err != nil {
		return err
	}

	for {
		var iterSum int

		for c := range g.Iter() {
			if c.Accessible == true {
				c.Accessible = false
			}

			if c.Roll && countAdjacent(g, c) < 4 {
				iterSum++
				c.Accessible = true
				c.Roll = false
			}
		}

		sum += iterSum

		if build.DEBUG {
			time.Sleep(time.Millisecond * 500)
			_ = g.Fprint(w)
		}

		if iterSum == 0 {
			break
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func countAdjacent(g *grid.Grid[Cell, *Cell], c *Cell) int {
	var rollCount int
	var a *Cell

	if a = g.Get(c.Y()+direction.Up.Y(), c.X()+direction.Left.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y()+direction.Up.Y(), c.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y()+direction.Up.Y(), c.X()+direction.Right.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y(), c.X()+direction.Left.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y(), c.X()+direction.Right.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y()+direction.Down.Y(), c.X()+direction.Left.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y()+direction.Down.Y(), c.X()); a != nil && a.Roll {
		rollCount++
	}

	if a = g.Get(c.Y()+direction.Down.Y(), c.X()+direction.Right.X()); a != nil && a.Roll {
		rollCount++
	}

	return rollCount
}

func parseInput(rd io.Reader) (*grid.Grid[Cell, *Cell], error) {
	g := grid.CreateGrid[Cell]()

	for y, line := range input.LineReader(rd) {
		err := g.AddRow(line, func(x int, r rune) (*Cell, error) {
			c := CreateCell(y, x, r)

			return c, nil
		})

		if build.DEBUG {
			fmt.Println("")
		}

		if err != nil {
			return g, fmt.Errorf("failed to parse line %d: %v", y, err)
		}
	}

	if build.DEBUG {
		// TODO Instead use the argument logger
		if err := g.Fprint(os.Stdout); err != nil {
			return g, err
		}
	}

	return g, nil
}

type Cell struct {
	grid.BaseCell
	Roll       bool
	Accessible bool
}

func CreateCell(y, x int, r rune) *Cell {
	return &Cell{
		BaseCell: grid.CreateBaseCell(y, x, r),
		Roll:     r == '@',
	}
}

func (c *Cell) Bytes() []byte {
	return []byte(string(c.Rune()))
}

func (c *Cell) Rune() rune {
	if c.Accessible {
		return 'X'
	}

	if c.Roll {
		return '@'
	}

	return '.'
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	if err := Solution(context.Background(), logger.CreateIterationLogger(context.Background()), f, os.Stdout); err != nil {
		panic(err)
	}
}
