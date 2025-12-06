package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/build"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/input"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tools/logger"
	"sort"
	"strconv"
)

func Solution(_ context.Context, _ *logger.IterationLogger, rd io.Reader, w io.Writer) error {
	var sum int

	ranges := make([]struct{ l, r int }, 0, 1024)
	procLine := procRangeLine

	for lineNr, line := range input.LineReader(rd) {
		if len(line) == 0 {
			if build.DEBUG {
				fmt.Printf("ranges:\n")
				for i, pair := range ranges {
					fmt.Printf(" %3d: %20d - %-20d\n", i, pair.l, pair.r)
				}
			}

			procLine = procIngredientLine
			continue
		}

		if err := procLine(lineNr, line, &ranges, &sum); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(strconv.Itoa(sum))); err != nil {
		return err
	}

	return nil
}

func procRangeLine(lineNr int, line string, ranges *[]struct{ l, r int }, _ *int) error {
	var l, r int

	_, err := fmt.Sscanf(line, "%d-%d", &l, &r)
	if err != nil {
		return err
	}

	if build.DEBUG {
		fmt.Printf("range (line %d): %d - %d\n", lineNr, l, r)
	}

	shouldAppend := true

	var i int
	for ; i < len(*ranges); i++ {
		pair := (*ranges)[i]

		// -1 & +1, so adjacency is also seen as overlap
		if r < (pair.l - 1) {
			// *ranges is sorted, r stays the same, pair.l is strictly increasing.
			// This condition will always be true from now on.
			i = len(*ranges) // Also skip the 'next' overlap check
			break
		} else if l > (pair.r + 1) {
			continue
		}

		if build.DEBUG {
			fmt.Printf(" > overlaps with range %d (%d - %d)\n", i, pair.l, pair.r)
		}

		shouldAppend = false

		if l >= pair.l && r <= pair.r {
			// Fully overlapping; skip checking for other overlaps
			i = len(*ranges)
			break
		}

		if l < pair.l {
			(*ranges)[i].l = l
		}

		if r > pair.r {
			(*ranges)[i].r = r
		}

		// Problem is, there could be more overlaps...
		// Good thing is, its only with the *ranges AFTER this one.
		// Bad thing, we can't continue with the 'remainder' only, as it might be on 2 sides.
		// So we continue iterating in the next for-loop, reducing overlaps we find, nil-ing zero length remainders.
		// However, it still leaves adjacent pairs if any were reduced, but not removed.
		break
	}

	for i++; i < len(*ranges); i++ {
		pair := (*ranges)[i]

		if l > pair.r {
			// No overlap
			continue
		} else if r < pair.l {
			// *ranges is sorted, r stays the same, pair.l is strictly increasing.
			// This condition will always be true from now on.
			break
		}

		if build.DEBUG {
			fmt.Printf(" > overlaps again with range %d (%d - %d)\n", i, pair.l, pair.r)
		}

		// Note; fully contained overlap (l > pair.l && r < pair.r) is not possible.
		// This loop would not trigger if there was no other overlap before
		// and any other overlap would then automatically overlap this one and been merged or negated.
		// An exact overlap is still possible though.

		if l <= pair.l {
			(*ranges)[i].l = r + 1
		} else if r >= pair.r {
			(*ranges)[i].r = l - 1
		}

		if pair.l > pair.r {
			// And... invalid, remove
			copy((*ranges)[i:], (*ranges)[i+1:])
			*ranges = (*ranges)[:len(*ranges)-1]
			i--
		}
	}

	if shouldAppend {
		index := sort.Search(len(*ranges), func(j int) bool { return (*ranges)[j].l >= l && ((*ranges)[j].l != l || (*ranges)[j].r >= r) })
		*ranges = append(*ranges, struct{ l, r int }{}) // Will be overwritten on the next line!
		copy((*ranges)[index+1:], (*ranges)[index:])
		(*ranges)[index] = struct{ l, r int }{l, r}
	}

	return nil
}

func procIngredientLine(lineNr int, line string, ranges *[]struct{ l, r int }, sum *int) error {
	if build.DEBUG {
		fmt.Printf("ingredient (line %d): %s\n", lineNr, line)
	}

	ids := *ranges

	id, err := strconv.Atoi(line)
	if err != nil {
		return err
	}

	index := sort.Search(len(ids), func(j int) bool { return ids[j].l >= id })

	if index > 0 && id >= ids[index-1].l && id <= ids[index-1].r {
		if build.DEBUG {
			fmt.Printf(" > FRESH\n")
		}

		*sum++

		return nil
	}

	if build.DEBUG {
		fmt.Printf(" > spoiled\n")
	}

	return nil
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
