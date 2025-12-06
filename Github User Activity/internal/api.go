package internal

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

const baseURL = "https://api.github.com"

type Event struct {
    Type      string    `json:"type"`
    Repo      Repo      `json:"repo"`
    CreatedAt time.Time `json:"created_at"`
    Payload   Payload   `json:"payload"`
}

type Repo struct {
    Name string `json:"name"`
}

type Payload struct {
    Action       string `json:"action"`
    Size         int    `json:"size"`
    DistinctSize int    `json:"distinct_size"`
    Ref          string `json:"ref"`
    RefType      string `json:"ref_type"`
    Number       int    `json:"number"`
    Forkee       struct {
        FullName string `json:"full_name"`
    } `json:"forkee"`
    Release struct {
        TagName string `json:"tag_name"`
    } `json:"release"`
    Commits []struct {
        Sha     string `json:"sha"`
        Message string `json:"message"`
        Distinct bool  `json:"distinct"`
    } `json:"commits"`
}

type RateInfo struct {
    Remaining int
    Reset     int64
    ResetTime time.Time
}

func FetchUserEvents(username string) ([]Event, RateInfo, error) {
    url := fmt.Sprintf("%s/users/%s/events", baseURL, username)

    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, RateInfo{}, err
    }
    req.Header.Set("User-Agent", "github-activity-cli")

    if token := os.Getenv("GITHUB_TOKEN"); token != "" {
        req.Header.Set("Authorization", "Bearer "+token)
    }

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return nil, RateInfo{}, err
    }
    defer resp.Body.Close()

    rate := parseRate(resp.Header)

    if resp.StatusCode == http.StatusNotFound {
        return nil, rate, fmt.Errorf("user '%s' not found", username)
    }
    if resp.StatusCode == http.StatusForbidden && rate.Remaining == 0 {
        return nil, rate, fmt.Errorf("rate limit exceeded")
    }
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, rate, fmt.Errorf("GitHub API error: %s", resp.Status)
    }

    var events []Event
    if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
        return nil, rate, err
    }
    return events, rate, nil
}

func parseRate(h http.Header) RateInfo {
    remaining := intHeader(h.Get("X-RateLimit-Remaining"))
    reset := int64Header(h.Get("X-RateLimit-Reset"))
    var resetTime time.Time
    if reset > 0 {
        resetTime = time.Unix(reset, 0)
    }
    return RateInfo{
        Remaining: remaining,
        Reset:     reset,
        ResetTime: resetTime,
    }
}

func intHeader(s string) int {
    var v int
    fmt.Sscanf(s, "%d", &v)
    return v
}

func int64Header(s string) int64 {
    var v int64
    fmt.Sscanf(s, "%d", &v)
    return v
}

func HumanizeEvent(e Event) string {
    switch e.Type {
    case "PushEvent":
        count := e.Payload.DistinctSize
        if count == 0 {
            count = e.Payload.Size
        }
        if count == 0 && e.Payload.Commits != nil {
            count = len(e.Payload.Commits)
        }

        if count == 0 {
            return fmt.Sprintf("Pushed updates to %s", e.Repo.Name)
        }
        if count == 1 {
            return fmt.Sprintf("Pushed 1 commit to %s", e.Repo.Name)
        }
        return fmt.Sprintf("Pushed %d commits to %s", count, e.Repo.Name)

    case "PublicEvent":
        return fmt.Sprintf("Made %s public", e.Repo.Name)

    case "IssuesEvent":
        switch e.Payload.Action {
        case "opened":
            return fmt.Sprintf("Opened a new issue in %s", e.Repo.Name)
        case "closed":
            return fmt.Sprintf("Closed an issue in %s", e.Repo.Name)
        case "reopened":
            return fmt.Sprintf("Reopened an issue in %s", e.Repo.Name)
        default:
            return fmt.Sprintf("Updated an issue in %s", e.Repo.Name)
        }

    case "IssueCommentEvent":
        return fmt.Sprintf("Commented on an issue in %s", e.Repo.Name)

    case "PullRequestEvent":
        switch e.Payload.Action {
        case "opened":
            return fmt.Sprintf("Opened a pull request in %s", e.Repo.Name)
        case "closed":
            return fmt.Sprintf("Closed a pull request in %s", e.Repo.Name)
        case "merged":
            return fmt.Sprintf("Merged a pull request in %s", e.Repo.Name)
        default:
            return fmt.Sprintf("Updated a pull request in %s", e.Repo.Name)
        }

    case "PullRequestReviewCommentEvent":
        return fmt.Sprintf("Reviewed a pull request in %s", e.Repo.Name)

    case "CreateEvent":
        if e.Payload.RefType == "repository" {
            return fmt.Sprintf("Created repository %s", e.Repo.Name)
        }
        return fmt.Sprintf("Created %s %s in %s", e.Payload.RefType, e.Payload.Ref, e.Repo.Name)

    case "DeleteEvent":
        if e.Payload.RefType == "repository" {
            return fmt.Sprintf("Deleted repository %s", e.Repo.Name)
        }
        return fmt.Sprintf("Deleted %s %s in %s", e.Payload.RefType, e.Payload.Ref, e.Repo.Name)

    case "ForkEvent":
        if e.Payload.Forkee.FullName != "" {
            return fmt.Sprintf("Forked %s to %s", e.Repo.Name, e.Payload.Forkee.FullName)
        }
        return fmt.Sprintf("Forked %s", e.Repo.Name)

    case "WatchEvent":
        return fmt.Sprintf("Starred %s", e.Repo.Name)

    case "ReleaseEvent":
        tag := e.Payload.Release.TagName
        if tag == "" {
            return fmt.Sprintf("Published a release in %s", e.Repo.Name)
        }
        return fmt.Sprintf("Published release %s in %s", tag, e.Repo.Name)

    default:
        return fmt.Sprintf("%s in %s", e.Type, e.Repo.Name)
    }
}