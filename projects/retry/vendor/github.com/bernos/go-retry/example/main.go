package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bernos/go-retry"
)

func main() {
	r := retry.Retry(
		numberOrBust(128, 255),
		retry.MaxRetries(5),
		retry.BaseDelay(time.Millisecond))

	value, err := r()

	if err != nil {
		fmt.Printf("Bad luck! %s\n", err.Error())
	} else {
		fmt.Printf("Jackpot! You rolled a %d\n", value)
	}
}

// numberOrBust creates a func that chooses a random number between 0 and maxNumber
// and returns an error if that random number does not match the value of magicNumber
func numberOrBust(magicNumber int, maxNumber int) func() (interface{}, error) {
	return func() (interface{}, error) {
		guess := rand.Intn(maxNumber)

		if guess == magicNumber {
			return "Got it!", nil
		}

		return nil, fmt.Errorf("Want %d, got %d", magicNumber, guess)
	}
}
