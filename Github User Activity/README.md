# Subject: Github User Activity CLI

Goal
- Build a simple command line interface (CLI) to fetch the recent activity of a GitHub user and display it in the terminal.

You’ll Practice
- Consuming REST APIs (GitHub API)
- JSON decoding with standard libraries
- Command-line argument parsing
- Formatting and grouping console output

Run
1) cd "Github User Activity"
2) go build -o github-activity
3) ./github-activity <username>

Behavior
- Accepts a GitHub username as a CLI argument.
- Fetches recent events from `https://api.github.com/users/<username>/events`.
- Displays activity in a human-readable format (e.g., "Pushed updates to...", "Opened a new issue in...").
- Groups consecutive identical events to reduce noise (e.g., "(x3)").
- Handles errors gracefully (User not found, API rate limits).
- No external libraries used for API requests.

Try
- Build:        go build -o github-activity
- Run (yours):  ./github-activity you
- Error check:  ./github-activity non_existent_user_123