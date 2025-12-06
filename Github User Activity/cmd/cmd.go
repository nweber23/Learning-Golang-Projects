package cmd

import (
    "fmt"

    "github.com/nweber/github-activity/internal"
)

func Run(username string) error {
    events, rate, err := internal.FetchUserEvents(username)
    if err != nil {
        if rate.Remaining == 0 {
            return fmt.Errorf("GitHub API rate limit exceeded. Reset at %v", rate.ResetTime)
        }
        return err
    }

    if len(events) == 0 {
        fmt.Println("No recent activity found.")
        return nil
    }

    var lastMsg string
    var count int

    for _, e := range events {
        currentMsg := internal.HumanizeEvent(e)

        if currentMsg == lastMsg {
            count++
        } else {
            if lastMsg != "" {
                printEvent(lastMsg, count)
            }
            lastMsg = currentMsg
            count = 1
        }
    }
    if lastMsg != "" {
        printEvent(lastMsg, count)
    }

    return nil
}

func printEvent(msg string, count int) {
    if count > 1 {
        fmt.Printf("- %s (x%d)\n", msg, count)
    } else {
        fmt.Printf("- %s\n", msg)
    }
}
