package main

import (
	"fmt"
	"sort"

	"heliosdrift/render"
)

func colorStats(fb *render.FB, x0, y0, x1, y1 int, label string) {
	counts := map[uint32]int{}
	total := 0
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			counts[fb.Buf[y*fb.W+x]]++
			total++
		}
	}
	type kv struct {
		c uint32
		n int
	}
	var list []kv
	for c, n := range counts {
		list = append(list, kv{c, n})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].n > list[j].n })
	out := ""
	for i := 0; i < len(list) && i < 6; i++ {
		c := list[i].c
		out += fmt.Sprintf(" #%02x%02x%02x:%.0f%%", byte(c), byte(c>>8), byte(c>>16), float64(list[i].n)/float64(total)*100)
	}
	fmt.Printf("%-16s colours=%-3d%s\n", label, len(counts), out)
}
