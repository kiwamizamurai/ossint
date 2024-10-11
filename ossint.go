package ossint

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/aquasecurity/table"
)

const (
	githubAPIURL = "https://api.github.com"
)

type PullRequest struct {
	Title          string `json:"title"`
	Repo           string
	Stars          int
	URL            string `json:"html_url"`
	Status         string
	Number         int    `json:"number"`
	RepoURL        string `json:"repository_url"`
	Additions      int    `json:"additions"`
	Deletions      int    `json:"deletions"`
	ChangedFiles   int    `json:"changed_files"`
	FileExtensions []string
}

func Run(argv []string, outStream, errStream io.Writer) error {
	fs := flag.NewFlagSet(
		fmt.Sprintf("ossint (v%s rev:%s)", version, revision), flag.ContinueOnError)
	fs.SetOutput(errStream)
	v := fs.Bool("version", false, "display version")
	username := fs.String("username", "", "GitHub username")
	token := fs.String("token", "", "GitHub token")

	if err := fs.Parse(argv); err != nil {
		return err
	}

	if *v {
		return printVersion(outStream)
	}

	if *username == "" {
		return fmt.Errorf("GitHub username is required. Set it with -username flag")
	}

	if *token == "" {
		*token = os.Getenv("GITHUB_TOKEN")
	}

	if *token == "" {
		fmt.Fprintln(errStream, "GitHub token not provided. Attempting to use 'gh auth token'...")
		var err error
		*token, err = getGitHubTokenFromCLI()
		if err != nil {
			return fmt.Errorf("error getting GitHub token: %v\nPlease provide a token with -token flag or set GITHUB_TOKEN environment variable", err)
		}
	}

	return run(*username, *token, outStream, errStream)
}

func printVersion(out io.Writer) error {
	_, err := fmt.Fprintf(out, "ossint v%s (rev:%s)\n", version, revision)
	return err
}

func run(username, token string, outStream, errStream io.Writer) error {
	prs, err := getUserPRs(username, token)
	if err != nil {
		return fmt.Errorf("error getting user PRs: %v", err)
	}

	for i, pr := range prs {
		repoFullName := getRepoFullName(pr.RepoURL)
		stars, err := getRepoStars(repoFullName, token)
		if err != nil {
			fmt.Fprintf(errStream, "Error getting stars for repo %s: %v\n", repoFullName, err)
			continue
		}
		prs[i].Repo = repoFullName
		prs[i].Stars = stars

		status, err := getPRStatus(repoFullName, pr.Number, token)
		if err != nil {
			fmt.Fprintf(errStream, "Error getting PR status for %s/%d: %v\n", repoFullName, pr.Number, err)
			continue
		}
		prs[i].Status = status

		prDetails, err := getPRDetails(repoFullName, pr.Number, token)
		if err != nil {
			fmt.Fprintf(errStream, "Error getting PR details for %s/%d: %v\n", repoFullName, pr.Number, err)
			continue
		}
		prs[i].Additions = prDetails.Additions
		prs[i].Deletions = prDetails.Deletions
		prs[i].ChangedFiles = prDetails.ChangedFiles
		prs[i].FileExtensions = prDetails.FileExtensions
	}

	for i, pr := range prs {
		repoFullName := getRepoFullName(pr.RepoURL)
		stars, err := getRepoStars(repoFullName, token)
		if err != nil {
			fmt.Fprintf(errStream, "Error getting stars for repo %s: %v\n", repoFullName, err)
			continue
		}
		prs[i].Repo = repoFullName
		prs[i].Stars = stars

		status, err := getPRStatus(repoFullName, pr.Number, token)
		if err != nil {
			fmt.Fprintf(errStream, "Error getting PR status for %s/%d: %v\n", repoFullName, pr.Number, err)
			continue
		}
		prs[i].Status = status
	}

	sort.Slice(prs, func(i, j int) bool {
		return prs[i].Stars > prs[j].Stars
	})

	t := table.New(outStream)
	t.SetRowLines(false)
	t.SetBorders(false)
	t.SetDividers(table.UnicodeRoundedDividers)
	t.SetHeaderStyle(table.StyleBold)
	t.SetLineStyle(table.StyleBlue)
	t.SetHeaders("PR", "Stars", "Title", "Additions", "Deletions", "Changed Files", "File Extensions")

	for _, pr := range prs {
		if pr.Status != "Closed" && !strings.Contains(pr.Repo, username) {
			t.AddRow(pr.URL, strconv.Itoa(pr.Stars), pr.Title, strconv.Itoa(pr.Additions), strconv.Itoa(pr.Deletions), strconv.Itoa(pr.ChangedFiles), strings.Join(pr.FileExtensions, ", "))
		}
	}

	t.Render()

	return nil
}

func getGitHubTokenFromCLI() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getUserPRs(username, token string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/search/issues?q=author:%s+type:pr+is:public", githubAPIURL, username)
	body, err := makeRequest(url, token)
	if err != nil {
		return nil, err
	}

	var result struct {
		Items []PullRequest `json:"items"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

func getRepoStars(repoFullName, token string) (int, error) {
	url := fmt.Sprintf("%s/repos/%s", githubAPIURL, repoFullName)
	body, err := makeRequest(url, token)
	if err != nil {
		return 0, err
	}

	var result struct {
		StargazersCount int `json:"stargazers_count"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result.StargazersCount, nil
}

func getPRStatus(repoFullName string, prNumber int, token string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls/%d", githubAPIURL, repoFullName, prNumber)
	body, err := makeRequest(url, token)
	if err != nil {
		return "", err
	}

	var result struct {
		State  string `json:"state"`
		Merged bool   `json:"merged"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.State == "open" {
		return "Open", nil
	} else if result.Merged {
		return "Merged", nil
	}
	return "Closed", nil
}

func makeRequest(url, token string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func getRepoFullName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return strings.Join(parts[len(parts)-2:], "/")
}

func getPRDetails(repoFullName string, prNumber int, token string) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/pulls/%d", githubAPIURL, repoFullName, prNumber)
	body, err := makeRequest(url, token)
	if err != nil {
		return nil, err
	}

	var pr PullRequest
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return nil, err
	}

	// Get file extensions
	filesURL := fmt.Sprintf("%s/repos/%s/pulls/%d/files", githubAPIURL, repoFullName, prNumber)
	filesBody, err := makeRequest(filesURL, token)
	if err != nil {
		return nil, err
	}

	var files []struct {
		Filename string `json:"filename"`
	}
	err = json.Unmarshal(filesBody, &files)
	if err != nil {
		return nil, err
	}

	extensions := make(map[string]bool)
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		if ext != "" {
			extensions[ext] = true
		}
	}

	pr.FileExtensions = make([]string, 0, len(extensions))
	for ext := range extensions {
		pr.FileExtensions = append(pr.FileExtensions, ext)
	}

	return &pr, nil
}
