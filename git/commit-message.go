package git

import (
	"io/ioutil"
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
	}
	return CommitMessage{
		Subject: subject,
		Body:    body,
	}
}

func (g *Git) GetInProgressCommitMessageEither() string {
	msg, err := g.GetInProgressCommitMessage(true)
	msg = parseOneLine(msg)
	if err != nil || msg == "" {
		msg, _ = g.GetInProgressCommitMessage(false)
	}

	return msg
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
			println(err)
		}
		return "", err
	}

	msg, err := ioutil.ReadFile(file)
	if err != nil {
		println(err)
		return "", err
	}

	return string(msg), nil
}
