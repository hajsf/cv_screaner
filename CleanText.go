package main

import "strings"

func trimEmptyLines(b string) string {
	strs := strings.Split(b, "\n")
	str := ""
	for _, s := range strs {
		if len(strings.TrimSpace(s)) == 0 {
			continue
		}
		str += s + "\n"
	}
	str = strings.TrimSuffix(str, "\n")

	return str
}
