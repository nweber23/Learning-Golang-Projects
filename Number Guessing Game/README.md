# Subject: Number Guessing Game in Go

Goal
- Build a CLI-based number guessing game where the computer picks a random number and the player guesses it.

You'll Practice
- Reading user input from stdin with bufio.Scanner
- Generating random numbers with math/rand
- Handling difficulty levels and game state
- String parsing and validation
- Building a game loop with multiple rounds
- Tracking high scores across sessions

Run
1) cd "Number Guessing Game"
2) go run main.go
3) Follow the on-screen prompts to play

Game Rules
- Computer randomly selects a number between 1 and 100
- Player selects a difficulty level:
  - Easy: 10 chances
  - Medium: 5 chances
  - Hard: 3 chances
- Each incorrect guess provides feedback (number is greater/less)
- Type "hint" at any time to get a clue about the number's range
- Game ends when player guesses correctly or runs out of chances

Behavior
- Welcome message displayed on startup
- Difficulty selection menu presented at the start of each game
- After each guess, remaining chances are shown
- Hint system reveals whether the number is in lower half (1-50) or upper half (51-100)
- Win message shows number of attempts and elapsed time
- High scores tracked for each difficulty level
- Player can choose to play again after each round

Try
- Start the game:
  - go run main.go
- Select difficulty (e.g., enter "2" for Medium)
- Make guesses when prompted
- Use "hint" command if stuck
- Play multiple rounds by answering "yes" to play again prompt

Project Files
- main.go
