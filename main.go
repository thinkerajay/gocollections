package main

import "log"

type LineItem struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func change(lmap map[string]*LineItem, id string) *LineItem {

	return lmap[id]
}

func main() {
	lmap := map[string]*LineItem{
		"445": &LineItem{ID: 4, Name: "one", Status: "active"},
	}

	l := change(lmap, "5")
	log.Println(l)
}
