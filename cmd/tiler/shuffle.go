// +build go1.10

package main

import (
	"math/rand"

	"github.com/thomersch/grandine/lib/tile"
)

func shuffleWork(w []tile.ID) {
	rand.Shuffle(len(w), func(i, j int) {
		w[i], w[j] = w[j], w[i]
	})
}
