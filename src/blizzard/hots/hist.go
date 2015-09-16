package main

import (
	"fmt"
	"sort"
)

type HistEntry struct {
	Name  string
	Count int
}
type Hist []HistEntry

func (h Hist) Len() int      { return len(h) }
func (h Hist) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func NewHist(counts map[string]int) Hist {
	h := make([]HistEntry, 0, len(counts))
	for name, count := range counts {
		h = append(h, HistEntry{name, count})
	}
	return h
}

func (h Hist) Print() {
	for _, he := range h {
		fmt.Printf("%d %s\n", he.Count, he.Name)
	}
}

func (h Hist) Desc() Hist {
	sort.Sort(Desc{h})
	return h
}

type Desc struct{ Hist }

func (h Desc) Less(i, j int) bool { return h.Hist[i].Count > h.Hist[j].Count }
