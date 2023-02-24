package orchestrator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubActions(t *testing.T) {
	t.Run("BranchBuild", func(t *testing.T) {
		defer resetEnv(os.Environ())
		os.Clearenv()
		os.Unsetenv("GITHUB_HEAD_REF")
		os.Setenv("GITHUB_ACTIONS", "true")
		os.Setenv("GITHUB_REF", "refs/heads/feat/test-gh-actions")
		os.Setenv("GITHUB_RUN_ID", "42")
		os.Setenv("GITHUB_SHA", "abcdef42713")
		os.Setenv("GITHUB_SERVER_URL", "github.com")
		os.Setenv("GITHUB_REPOSITORY", "foo/bar")

		p, _ := NewOrchestratorSpecificConfigProvider()

		assert.False(t, p.IsPullRequest())
		assert.Equal(t, "github.com/foo/bar/actions/runs/42", p.GetBuildURL())
		assert.Equal(t, "feat/test-gh-actions", p.GetBranch())
		assert.Equal(t, "refs/heads/feat/test-gh-actions", p.GetReference())
		assert.Equal(t, "abcdef42713", p.GetCommit())
		assert.Equal(t, "github.com/foo/bar", p.GetRepoURL())
		assert.Equal(t, "GitHubActions", p.OrchestratorType())
	})

	t.Run("PR", func(t *testing.T) {
		defer resetEnv(os.Environ())
		os.Clearenv()
		os.Setenv("GITHUB_HEAD_REF", "feat/test-gh-actions")
		os.Setenv("GITHUB_BASE_REF", "main")
		os.Setenv("GITHUB_REF", "refs/pull/42/merge")

		p := GitHubActionsConfigProvider{}
		c := p.GetPullRequestConfig()

		assert.True(t, p.IsPullRequest())
		assert.Equal(t, "feat/test-gh-actions", c.Branch)
		assert.Equal(t, "main", c.Base)
		assert.Equal(t, "42", c.Key)
	})

	// to be fixed
	// t.Run("Test log receiving", func(t *testing.T) {
	// 	defer resetEnv(os.Environ())
	// 	os.Clearenv()
	// 	os.Unsetenv("GITHUB_HEAD_REF")
	// 	os.Setenv("GITHUB_ACTIONS", "true")
	// 	os.Setenv("GITHUB_REF_NAME", "feat/test-gh-actions")
	// 	os.Setenv("GITHUB_REF", "refs/heads/feat/test-gh-actions")
	// 	os.Setenv("GITHUB_RUN_ID", "42")
	// 	os.Setenv("GITHUB_SHA", "abcdef42713")
	// 	os.Setenv("GITHUB_REPOSITORY", "foo/bar")
	// 	os.Setenv("GITHUB_URL", "https://github.com/")
	// 	p := func() OrchestratorSpecificConfigProviding {
	// 		g := GitHubActionsConfigProvider{}
	// 		g.client = piperHttp.Client{}
	// 		g.client.SetOptions(piperHttp.ClientOptions{
	// 			MaxRequestDuration:        5 * time.Second,
	// 			Token:                     "TOKEN",
	// 			TransportSkipVerification: true,
	// 			UseDefaultTransport:       true, // need to use default transport for http mock
	// 			MaxRetries:                -1,
	// 		})
	// 		return &g
	// 	}()
	// 	httpmock.Activate()
	// 	defer httpmock.DeactivateAndReset()
	// 	httpmock.RegisterResponder("GET", "https://api.github.com/repos/foo/bar/actions/runs/42/jobs",
	// 		func(req *http.Request) (*http.Response, error) {
	// 			return httpmock.NewJsonResponse(200, struct {
	// 				Jobs []struct {
	// 					Id string `json:"id"`
	// 				} `json:"jobs"`
	// 			}{
	// 				Jobs: []struct {
	// 					Id string `json:"id"`
	// 				}{
	// 					{
	// 						Id: "123",
	// 					},
	// 					{
	// 						Id: "124",
	// 					},
	// 					{
	// 						Id: "125",
	// 					},
	// 				},
	// 			})
	// 		},
	// 	)
	// 	logs := []string{
	// 		"log_record1\n",
	// 		"log_record2\n",
	// 	}
	// 	httpmock.RegisterResponder("GET", "https://api.github.com/repos/foo/bar/actions/jobs/123/logs",
	// 		func(req *http.Request) (*http.Response, error) {
	// 			return httpmock.NewStringResponse(200, logs[0]), nil
	// 		},
	// 	)
	// 	httpmock.RegisterResponder("GET", "https://api.github.com/repos/foo/bar/actions/jobs/124/logs",
	// 		func(req *http.Request) (*http.Response, error) {
	// 			return httpmock.NewStringResponse(200, logs[1]), nil
	// 		},
	// 	)
	// 	actual, _ := p.GetLog()
	// 	fmt.Println(string(actual))
	// 	assert.Equal(t, strings.Join(logs, ""), string(actual))
	// })

	// to be deleted
	// t.Run("Test log receiving 2", func(t *testing.T) {
	// 	defer resetEnv(os.Environ())
	// 	os.Clearenv()
	// 	os.Unsetenv("GITHUB_HEAD_REF")
	// 	os.Setenv("GITHUB_ACTIONS", "true")
	// 	os.Setenv("GITHUB_REF_NAME", "feat/test-gh-actions")
	// 	os.Setenv("GITHUB_REF", "refs/heads/feat/test-gh-actions")
	// 	os.Setenv("GITHUB_SHA", "abcdef42713")
	// 	os.Setenv("GITHUB_REPOSITORY", "project-piper/azure-demo-k8s-node")
	// 	os.Setenv("GITHUB_URL", "https://github.tools.sap/")
	// 	os.Setenv("GITHUB_TOKEN", "")
	// 	os.Setenv("GITHUB_RUN_ID", "1738520")
	// 	p := func() OrchestratorSpecificConfigProviding {
	// 		g := GitHubActionsConfigProvider{}
	// 		g.client = piperHttp.Client{}
	// 		g.client.SetOptions(piperHttp.ClientOptions{
	// 			MaxRequestDuration:        60 * time.Second,
	// 			Token:                     "TOKEN",
	// 			TransportSkipVerification: true,
	// 			UseDefaultTransport:       true, // need to use default transport for http mock
	// 			MaxRetries:                -1,
	// 		})
	// 		return &g
	// 	}()

	// 	b, err := p.GetLog()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	if err := os.WriteFile("test_log.log", b, 777); err != nil {
	// 		log.Fatal(err)
	// 	}
	// })
}
