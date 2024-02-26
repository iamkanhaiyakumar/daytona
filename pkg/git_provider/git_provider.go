package git_provider

import (
	"errors"

	"github.com/daytonaio/daytona/cli/config"
	"github.com/daytonaio/daytona/common/types"
)

const personalNamespaceId = "<PERSONAL>"

type GitProvider interface {
	GetNamespaces() ([]GitNamespace, error)
	GetRepositories(namespace string) ([]GitRepository, error)
	GetUserData() (GitUser, error)
}

type GitUser struct {
	Id       string
	Username string
	Name     string
	Email    string
}

type GitNamespace struct {
	Id   string
	Name string
}

type GitRepository struct {
	FullName string
	Name     string
	Url      string
}

func GetGitProvider(providerId string, gitProviders []types.GitProvider) GitProvider {
	var chosenProvider *types.GitProvider
	for _, gitProvider := range gitProviders {
		if gitProvider.Id == providerId {
			chosenProvider = &gitProvider
		}
	}

	if chosenProvider == nil {
		return nil
	}

	switch providerId {
	case "github":
		return &GitHubGitProvider{
			token: chosenProvider.Token,
		}
	case "gitlab":
		return &GitLabGitProvider{
			token: chosenProvider.Token,
		}
	case "bitbucket":
		return &BitbucketGitProvider{
			username: chosenProvider.Username,
			token:    chosenProvider.Token,
		}
	default:
		return nil
	}
}

func GetUsernameFromToken(providerId string, gitProviders []config.GitProvider, token string) (string, error) {
	var gitProvider GitProvider

	switch providerId {
	case "github":
		gitProvider = &GitHubGitProvider{
			token: token,
		}
	case "gitlab":
		gitProvider = &GitLabGitProvider{
			token: token,
		}
	case "bitbucket":
		gitProvider = &BitbucketGitProvider{
			token: token,
		}
	default:
		return "", errors.New("provider not found")
	}

	gitUser, err := gitProvider.GetUserData()
	if err != nil {
		return "", errors.New("user not found")
	}

	return gitUser.Username, nil
}