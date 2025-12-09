package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int
	points := make([]Point, 0, 1024)

	for lineNr, line := range input.LineReader(rd) {
		var x, y int
		_, err := fmt.Sscanf(line, "%d,%d", &x, &y)
		if err != nil {
			return fmt.Errorf("line %d failed to parse. %w", lineNr, err)
		}

		p := Point{x, y}

		for _, p2 := range points {
			sum = max(sum, area(p, p2))
		}

		points = append(points, p)
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func area(a, b Point) int {
	return (AbsInt(b.X-a.X) + 1) * (AbsInt(b.Y-a.Y) + 1)
}

func AbsInt(a int) int {
	if a < 0 {
		return -a
	}

	return a
}

func main() {
	f, err := os.Open("sample-input-1.txt")
	if err != nil {
		panic(err)
	}

	if err := Solution(context.Background(), logger.CreateIterationLogger(context.Background()), f, os.Stdout); err != nil {
		panic(err)
	}
}

type Point struct {
	X, Y int
}
