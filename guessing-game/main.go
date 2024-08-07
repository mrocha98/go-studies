package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
)

const maximumAttempts = 10
const minRandom = 1
const maxRandom = 100

func main() {
	fmt.Println("======\tGUESSING GAME\t======")
	fmt.Printf("You have %d attempts to guess the number.\n", maximumAttempts)
	fmt.Printf("Tip: it's an integer between %d and %d\n\n", minRandom, maxRandom)

	scanner := bufio.NewScanner(os.Stdin)
	attempts := [maximumAttempts]uint64{}
	random := rand.Uint64N(maxRandom) + 1

	for i := range maximumAttempts {
		fmt.Print("Type a number:\t")
		scanner.Scan()
		attempt, err := strconv.ParseUint(strings.TrimSpace(scanner.Text()), 10, 64)
		if err != nil {
			panic("Invalid guess")
		}

		attempts[i] = attempt
		switch {
		case attempt == random:
			fmt.Printf(
				"You win!\n"+
					"You made it in %d attempts.\n"+
					"Attempts: %v\n",
				i+1, attempts[:i+1],
			)
			return
		case attempt < random:
			fmt.Println("Wrong! Try a greater number!")
		case attempt > random:
			fmt.Println("Wrong! Try a lesser number!")
		}
	}

	fmt.Printf("Game Over!\n"+"The number was %d\n"+"You attempts: %v", random, attempts)
}
