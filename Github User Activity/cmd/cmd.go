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

    for _, e := range events {
        fmt.Println(internal.HumanizeEvent(e))
    }
    return nil
}
