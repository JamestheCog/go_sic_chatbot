// A file I came up with to store utility functions that've to do with
// the generation of unique IDs that've to do with session management.

package utils

import (
	"fmt"
	"math/rand/v2"
)

// -- Constants and structs --
const charSet = "abcdefghijklmnopqrstuvwxyz" +
	"1234567890" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// -- END --

// Given a positive integer, generate an ID of that said positive integer's
// length.  We gotta replace this with something more official
// if the app ever blows up, but I don't see this happening.
func GenerateID(idLen int) (string, error) {
	if idLen <= 0 {
		return "", fmt.Errorf("Cannot generate an ID of zero or negative length, man.")
	}

	toReturn := make([]byte, idLen)
	for i := range toReturn {
		toReturn[i] = charSet[rand.IntN(len(charSet))]
	}
	return string(toReturn), nil
}
