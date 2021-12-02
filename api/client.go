package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	UserData(username string) (*UserData, error)
	UserRepositories(username string) ([]*RepositoryData, error)
	LanguageData(username, repo string) (RepositoryLanguageData, error)
}

type ghClient struct {
	http          *http.Client
	user          string
	userRepos     string
	repoLanguages string
}

func NewClient() *ghClient {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	c := &ghClient{
		http:          httpClient,
		user:          "https://api.github.com/users/%s",
		userRepos:     "https://api.github.com/users/%s/repos",
		repoLanguages: "https://api.github.com/repos/%s/%s/languages",
	}

	return c
}

func (c *ghClient) UserData(username string) (*UserData, error) {
	url := fmt.Sprintf(c.user, username)
	req, err := c.request(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userData UserData
	err = json.NewDecoder(res.Body).Decode(&userData)
	if err != nil {
		return nil, err
	}
	return &userData, nil
}

func (c *ghClient) UserRepositories(username string) ([]*RepositoryData, error) {
	url := fmt.Sprintf(c.userRepos, username)
	req, err := c.request(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(bytes))

	userRepos := make([]*RepositoryData, 0)
	err = json.Unmarshal(bytes, &userRepos)
	if err != nil {
		return nil, err
	}
	return userRepos, nil
}

func (c *ghClient) LanguageData(username, repo string) (RepositoryLanguageData, error) {
	url := fmt.Sprintf(c.repoLanguages, username, repo)
	req, err := c.request(http.MethodGet, url)
	if err != nil {
		return nil, err
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	repoData := make(RepositoryLanguageData)
	err = json.NewDecoder(res.Body).Decode(&repoData)
	if err != nil {
		return nil, err
	}
	return repoData, nil
}

func (c *ghClient) request(method, u string) (*http.Request, error) {
	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}

	token := "some-auth-token"
	req.SetBasicAuth("username", token)
	return req, nil
}
