package common

import "math/rand/v2"

// ShuffleSlice shuffles a slice of strings in place using math/rand.
func ShuffleSlice(slice []string) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
