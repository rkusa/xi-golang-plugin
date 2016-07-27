package main

import "strconv"

var (
	color1 = parseHex("ff323333")
	color2 = parseHex("ff183792")
	color3 = parseHex("ffed6a43")
	color4 = parseHex("ff96989f")
	color5 = parseHex("ffa71e77")
)

func parseHex(s string) int64 {
	n, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		panic(err)
	}
	return n
}
