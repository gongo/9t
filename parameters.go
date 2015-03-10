package main

import "path"

var (
	// red, green, yellow, magenta, cyan
	ansiColorCodes = [...]int{31, 32, 33, 35, 36}
)

func getColorCode(index int) int {
	return ansiColorCodes[index%len(ansiColorCodes)]
}

func maximumNameLength(filenames []string) int {
	max := 0
	for _, name := range filenames {
		base := path.Base(name)
		if len(base) > max {
			max = len(base)
		}
	}
	return max
}
