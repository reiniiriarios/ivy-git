package git

import (
	"os"
	"path/filepath"
	"strings"
)

type CommitMessage struct {
	Subject string
	Body    string
}

func (g *Git) ParseCommitMessage(msg string) CommitMessage {
	msg = strings.ReplaceAll(msg, "\r\n", "\n")
	subject := msg
	body := ""
	if strings.Contains(msg, "\n") {
		parts := strings.SplitN(msg, "\n", 2)
		subject = parts[0]
		body = parts[1]
		if len(body) > 1 && body[:1] == "\n" {
			body = body[1:]
		}
		body = removeCommentLines(body, "#")
	}
	return CommitMessage{
		Subject: subject,
		Body:    body,
	}
}

func (g *Git) GetInProgressCommitMessage(merge bool) (string, error) {
	var file string
	if merge {
		file = "MERGE_MSG"
	} else {
		file = "COMMIT_EDITMSG"
	}
	file = filepath.Join(g.Repo.Directory, ".git", file)

	if _, err := os.Stat(file); err != nil {
		if !os.IsNotExist(err) {
			println("Error (GetInProgressCommitMessage, Stat)", err)
		}
		return "", err
	}

	msg, err := os.ReadFile(file)
	if err != nil {
		println("Error (GetInProgressCommitMessage, ReadFile)", err)
		return "", err
	}

	return string(msg), nil
}

func removeCommentLines(input string, starts_with string) string {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	trimmed := []string{}

	for i := range lines {
		if !strings.HasPrefix(lines[i], starts_with) {
			trimmed = append(trimmed, lines[i])
		}
	}

	return strings.TrimSpace(strings.Join(trimmed, "\n"))
}
