package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type GameState struct {
	secretNumber int
	chances      int
	attempts     int
	difficulty   string
	startTime    time.Time
}

type HighScores struct {
	Easy   int
	Medium int
	Hard   int
}

var scores = HighScores{
	Easy:   999,
	Medium: 999,
	Hard:   999,
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	displayWelcome()
	for {
		game := startNewGame(scanner)
		playGame(&game, scanner)
		fmt.Println("\nDo you want to play again?: ")
		response := readInput(scanner)
		if !strings.EqualFold(response, "yes") && !strings.EqualFold(response, "y") {
			fmt.Println("\n Thanks for playing! Goodbye!")
			break
		}
		fmt.Println()
	}
}

func displayWelcome() {
	fmt.Println("\n╔════════════════════════════════════════╗")
	fmt.Println("║  Welcome to the Number Guessing Game!  ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println("\nI'm thinking of a number between 1 and 100.")
	fmt.Println("Can you guess it?")
}

func selectDifficulty(scanner *bufio.Scanner) (string, int) {
	fmt.Println("Please select the difficulty level:")
	fmt.Println("1. Easy (10 chances)")
	fmt.Println("2. Medium (5 chances)")
	fmt.Println("3. Hard (3 chances)")
	fmt.Print("Enter your choice (1-3): ")

	for {
		choice := readInput(scanner)
		switch choice {
		case "1":
			fmt.Println("Great! You have selected Easy difficulty level.")
			return "Easy", 10
		case "2":
			fmt.Println("Great! You have selected Medium difficulty level.")
			return "Medium", 5
		case "3":
			fmt.Println("Great! You have selected Hard difficulty level.")
			return "Hard", 3
		default:
			fmt.Print("Invalid choice. Please enter 1, 2, or 3: ")
		}
	}
}

func startNewGame(scanner *bufio.Scanner) GameState {
	difficulty, chances := selectDifficulty(scanner)

	return GameState{
		secretNumber: rand.Intn(100) + 1,
		chances:      chances,
		attempts:     0,
		difficulty:   difficulty,
		startTime:    time.Now(),
	}
}

func playGame(game *GameState, scanner *bufio.Scanner) {
	fmt.Println("\nLet's start the game!")
	fmt.Printf("You have %d chances to guess the number.\n\n", game.chances)

	for game.chances > 0 {
		fmt.Printf("Chances remaining: %d\n", game.chances)
		fmt.Print("Enter your guess (or 'hint' for a clue): ")

		input := readInput(scanner)

		if strings.EqualFold(input, "hint") {
			displayHint(game)
			continue
		}

		guess, err := strconv.Atoi(input)
		if err != nil || guess < 1 || guess > 100 {
			fmt.Println("Invalid input. Please enter a number between 1 and 100.")
			continue
		}

		game.attempts++
		game.chances--

		if guess == game.secretNumber {
			displayWinMessage(game)
			updateHighScore(game)
			return
		}

		if guess < game.secretNumber {
			fmt.Printf("Incorrect! The number is greater than %d\n", guess)
		} else {
			fmt.Printf("Incorrect! The number is less than %d\n", guess)
		}
		fmt.Println()
	}

	displayLoseMessage(game)
}

func displayHint(game *GameState) {
	if game.secretNumber < 50 {
		fmt.Println("Your secret number is in the lower half (1-50)")
	} else {
		fmt.Println("Your secret number is in the upper half (1-50)")
	}
	fmt.Println()
}

func displayWinMessage(game *GameState) {
	elapsed := time.Since(game.startTime).Seconds()
	fmt.Printf("\nCongratulations! You guessed the correct number (%d) in %d attempts!\n", game.secretNumber,
		game.attempts)
	fmt.Printf("Time taken: %.1f seconds\n", elapsed)
	displayStats(game)
}

func displayLoseMessage(game *GameState) {
	fmt.Println("\nYou lost the game!")
	fmt.Printf("The correct number was : %d\n", game.secretNumber)
}

func displayStats(_ *GameState) {
	fmt.Println("\nHigh Scores by Difficulty:")
	fmt.Printf("  Easy:   %d attempts\n", scores.Easy)
	fmt.Printf("  Medium: %d attempts\n", scores.Medium)
	fmt.Printf("  Hard:   %d attempts\n", scores.Hard)
}

func updateHighScore(game *GameState) {
	switch game.difficulty {
	case "Easy":
		if game.attempts < scores.Easy {
			scores.Easy = game.attempts
			fmt.Printf("New High Score for Easy!\n")
		}
	case "Medium":
		if game.attempts < scores.Medium {
			scores.Medium = game.attempts
			fmt.Printf("New High Score for Medium!\n")
		}
	case "Hard":
		if game.attempts < scores.Hard {
			scores.Hard = game.attempts
			fmt.Printf("New High Score for Hard!\n")
		}
	}
}

func readInput(scanner *bufio.Scanner) string {
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
