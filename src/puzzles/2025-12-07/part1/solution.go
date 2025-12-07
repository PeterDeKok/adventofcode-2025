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
	var beams []bool

	for lineNr, line := range input.LineReader(rd) {
		if lineNr == 0 {
			beams = make([]bool, len(line))

			for i, r := range line {
				if r == 'S' {
					beams[i] = true
					if build.DEBUG {
						fmt.Printf("│")
					}
				} else if build.DEBUG {
					fmt.Printf(".")
				}
			}

			if build.DEBUG {
				fmt.Printf("\n")
			}

			continue
		}

		for i, r := range line {
			if r == '^' && beams[i] {
				sum++
				beams[i-1] = true
				beams[i] = false
				beams[i+1] = true
				if build.DEBUG {
					fmt.Printf("\033[D│^")
				}
			} else if build.DEBUG {
				if beams[i] {
					fmt.Printf("│")
				} else {
					fmt.Printf(".")
				}
			}
		}

		if build.DEBUG {
			fmt.Printf("\n")
		}
	}

	if build.DEBUG {
		fmt.Printf("\n\n")
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
