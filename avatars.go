package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// https://api.github.com/users/<username>
type GitHubUserData struct {
	Login     string `json:"login"`
	Id        uint64 `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

var githubUrlRegex = regexp.MustCompile(`^(\d+)\+([^@]+)@users.noreply.github.com$`)

const AVATAR_SIZE = 128 // pixels

func (a *App) GetAvatarUrl(email string) string {
	// Anon GitHub email addresses are not on Gravatar (directly):
	matches := githubUrlRegex.FindStringSubmatch(email)
	if len(matches) > 2 {
		return getGitHubAvatarUrl(matches[1], matches[2])
	}
	// Gravatar
	hash := md5.Sum([]byte(email))
	// d=identicon will return an anon avatar if none is found
	// d=404 will cause an error
	// otherwise, the gravatar logo will load if not found
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon&size=%d", hash, AVATAR_SIZE)
}

func getGitHubAvatarUrl(user_id string, username string) string {
	// Bot email addresses are /in/<id> where <id> does not match the user id.
	// As such, we must fetch this data.

	// todo: This requires an API key. :(
	/*
		if strings.HasSuffix(username, "[bot]") {
			println("bot", username)
			data, err := getGitHubUserData(username)
			if err != nil {
				println(err.Error())
				return fmt.Sprintf("https://avatars.githubusercontent.com/u/%s?s=64&v=4", user_id)
			}
			return data.AvatarUrl
		}
	*/
	return fmt.Sprintf("https://avatars.githubusercontent.com/u/%s?s=%d&v=4", user_id, AVATAR_SIZE)
}

func getGitHubUserData(username string) (GitHubUserData, error) {
	var data GitHubUserData
	resp, err := http.Get("https://api.github.com/users/" + username)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	if resp.StatusCode != http.StatusOK {
		return data, errors.New(resp.Status)
	}
	json.Unmarshal(body, &data)
	return data, nil
}
