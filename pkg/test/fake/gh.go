package fake

import (
	"fmt"
	"os"
	"strings"
)

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Gh struct {
	Table   map[string]Response
	LogPath string
}

func (g Gh) ParseArgs(osArgs []string) ([]string, error) {
	sep := -1
	for i, a := range osArgs {
		if a == "--" {
			sep = i
			break
		}
	}
	if sep < 0 || sep+1 >= len(osArgs) {
		return nil, fmt.Errorf("missing -- separator")
	}
	return osArgs[sep+1:], nil
}

func (g Gh) Log(ghArgs []string) error {
	if g.LogPath == "" {
		return nil
	}

	f, err := os.OpenFile(g.LogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := f.WriteString(strings.Join(ghArgs, " ") + "\n"); err != nil {
		return err
	}
	return nil
}

func (g Gh) Key(ghArgs []string) (string, bool) {
	parts := ghArgs
	if len(parts) > 0 && parts[0] == "gh" {
		parts = parts[1:]
	}
	if len(parts) < 2 {
		return "", false
	}
	base := parts[0] + " " + parts[1]
	if base == "api graphql" {
		if op := graphqlOp(parts[2:]); op != "" {
			return base + " " + op, true
		}
	}
	return base, true
}

// graphqlOps maps GraphQL query text substrings to operation keys.
// More specific strings must appear before any that are substrings of them
// (e.g. "deletePullRequestReviewComment" before "deletePullRequestReview").
var graphqlOps = []string{
	"headRefOid",
	"addPullRequestReviewThread",
	"submitPullRequestReview",
	"addPullRequestReview",
	"deletePullRequestReviewComment",
	"updatePullRequestReviewComment",
	"deletePullRequestReview",
}

func graphqlOp(ghArgs []string) string {
	for i, a := range ghArgs {
		if (a == "-f" || a == "-F") && i+1 < len(ghArgs) {
			v := ghArgs[i+1]
			if strings.HasPrefix(v, "query=") {
				q := v[len("query="):]
				for _, op := range graphqlOps {
					if strings.Contains(q, op) {
						return op
					}
				}
				break
			}
		}
	}
	return ""
}

func (g Gh) Find(key string) (Response, bool) {
	resp, ok := g.Table[key]
	return resp, ok
}

func (g Gh) Write(resp Response) {
	if resp.Stdout != "" {
		fmt.Print(resp.Stdout)
	}
	if resp.Stderr != "" {
		fmt.Fprint(os.Stderr, resp.Stderr)
	}
}
