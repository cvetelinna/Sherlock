package aggregate

import (
	"fmt"
	"github_api/api"
)

type UserAggregation struct {
	Username             string
	ReposCount           int
	FollowersCount       int
	ForksCount           int
	LanguageDistribution langDistribution
}

type aggregator struct {
	client api.Client
}

type langDistribution map[string]float64

func New(client api.Client) *aggregator {
	return &aggregator{client: client}
}

func (a *aggregator) AggregateUser(username string) (*UserAggregation, error) {
	userData, err := a.client.UserData(username)
	if err != nil {
		return nil, err
	}
	repos, err := a.client.UserRepositories(username)
	if err != nil {
		return nil, err
	}

	var res UserAggregation
	res.Username = userData.Login
	res.ReposCount = len(repos)
	res.FollowersCount = userData.Followers

	languageUsages := make(api.RepositoryLanguageData)
	for _, repo := range repos {
		res.ForksCount += repo.ForksCount
		languageData, err := a.client.LanguageData(username, repo.Name)
		if err != nil {
			return nil, err
		}
		for lang, count := range languageData {
			languageUsages[lang] += count
		}
	}

	res.LanguageDistribution = a.languageUsageDistribution(languageUsages)

	return &res, nil
}

func (a *aggregator) languageUsageDistribution(langCounts api.RepositoryLanguageData) langDistribution {
	res := make(langDistribution)
	sum := 0
	for _, count := range langCounts {
		sum += count
	}
	for lang, count := range langCounts {
		res[lang] = float64(count) / float64(sum)
	}
	return res
}

func (a *aggregator) Print(data *UserAggregation) {
	fmt.Printf("Username: %s\n", data.Username)
	fmt.Printf("Repos count: %d\n", data.ReposCount)
	fmt.Printf("Followers count: %d\n", data.FollowersCount)
	fmt.Printf("Forks count: %d\n", data.ForksCount)

	for lang, dist := range data.LanguageDistribution {
		fmt.Printf("%s : %.4f\n", lang, dist)
	}
}
