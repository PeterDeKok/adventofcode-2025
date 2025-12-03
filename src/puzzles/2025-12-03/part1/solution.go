package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	for _, line := range input.LineReader(rd) {
		left := 0
		right := 0

		if build.DEBUG {
			fmt.Printf("bank: %s\n", line)
		}

		for i, c := range line[:len(line)-1] {
			l := int(c) - '0'
			r := int(line[i+1]) - '0'

			if l > left {
				left = l
				right = r
			} else if r > right {
				right = r
			}
		}

		sum += 10*left + right
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
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
