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
	var intermediate []struct{ plus, mult int }

	for lineNr, values := range input.LineIsStrSliceReader(rd) {
		if lineNr == 0 {
			if build.DEBUG {
				fmt.Printf("First row (%d values)\n", len(values))
			}

			intermediate = make([]struct{ plus, mult int }, len(values))

			for i, value := range values {
				v, err := strconv.Atoi(value)
				if err != nil {
					return err
				}

				intermediate[i].plus = v
				intermediate[i].mult = v
			}

			continue
		}

		if values[0][0] == '+' || values[0][0] == '*' {
			if build.DEBUG {
				fmt.Printf("Last row\n")
			}

			for i, op := range values {
				switch op[0] {
				case '+':
					sum += intermediate[i].plus
				case '*':
					sum += intermediate[i].mult
				default:
					return fmt.Errorf("Invalid field %d: %s\n", i, op)
				}
			}

			break
		}

		if build.DEBUG {
			fmt.Printf("Row %d\n", lineNr)
		}

		for i, value := range values {
			v, err := strconv.Atoi(value)
			if err != nil {
				return err
			}

			intermediate[i].plus += v
			intermediate[i].mult *= v
		}
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
