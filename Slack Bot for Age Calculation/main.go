package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shomali11/slacker"
)

func printCommandEvent() {
	for event := range bot.CommandEvents() {
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {
	os.Setenv("SLACK_BOT_TOKEN", "")
	os.Setenv("SLACK_APP_TOKEN", "")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvent(bot.CommandEvents())

	bot.Command("calculate age <birthyear>", &slacker.CommandDefinition{
		Description: "Calculate age from birth year",
		Example:     "calculate age 1990",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			birthYearStr := request.Param("birthyear")
			birthYear, err := strconv.Atoi(birthYearStr)
			if err != nil {
				response.Reply("Please provide a valid year.")
				return
			}

			currentYear := time.Now().Year()
			age := currentYear - birthYear

			reply := fmt.Sprintf("You are %d years old.", age)
			response.Reply(reply)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}