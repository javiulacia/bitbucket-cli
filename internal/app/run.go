package app

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"bitbucket-cli/internal/bitbucket"
	"bitbucket-cli/internal/config"
)

func Run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "auth":
		return runAuth(args[1:])
	case "repo":
		return runRepo(args[1:])
	case "pr":
		return runPR(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runAuth(args []string) error {
	if len(args) == 0 {
		return errors.New("expected auth subcommand: login, status, logout")
	}
	switch args[0] {
	case "login":
		fs := flag.NewFlagSet("auth login", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		email := fs.String("email", "", "Bitbucket account email, used with --api-token")
		apiToken := fs.String("api-token", "", "Bitbucket API token")
		accessToken := fs.String("access-token", "", "Bitbucket OAuth/access token")
		workspace := fs.String("workspace", "", "Default workspace")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		cfg := &config.Config{Workspace: strings.TrimSpace(*workspace), GitProtocol: "https"}
		switch {
		case *apiToken != "":
			if strings.TrimSpace(*email) == "" {
				return errors.New("--email is required with --api-token")
			}
			cfg.AuthMethod = config.AuthMethodAPIToken
			cfg.Email = strings.TrimSpace(*email)
			cfg.APIToken = strings.TrimSpace(*apiToken)
		case *accessToken != "":
			cfg.AuthMethod = config.AuthMethodAccessToken
			cfg.AccessToken = strings.TrimSpace(*accessToken)
		default:
			return errors.New("provide either --api-token with --email, or --access-token")
		}

		client, err := bitbucket.NewClient(cfg)
		if err != nil {
			return err
		}
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			return fmt.Errorf("authentication check failed: %w", err)
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("Authenticated as %s\n", displayUser(user))
		return nil

	case "status":
		cfg, err := requireConfig()
		if err != nil {
			return err
		}
		client, err := bitbucket.NewClient(cfg)
		if err != nil {
			return err
		}
		user, err := client.CurrentUser(context.Background())
		if err != nil {
			return err
		}
		fmt.Printf("Authenticated as %s\n", displayUser(user))
		fmt.Printf("Auth method: %s\n", cfg.AuthMethod)
		if cfg.Workspace != "" {
			fmt.Printf("Default workspace: %s\n", cfg.Workspace)
		}
		return nil

	case "logout":
		return config.Delete()
	default:
		return fmt.Errorf("unknown auth subcommand %q", args[0])
	}
}

func runRepo(args []string) error {
	if len(args) == 0 {
		return errors.New("expected repo subcommand: list, view, create")
	}
	cfg, err := requireConfig()
	if err != nil {
		return err
	}
	client, err := bitbucket.NewClient(cfg)
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		fs := flag.NewFlagSet("repo list", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		ws, err := resolveWorkspace(cfg, *workspace)
		if err != nil {
			return err
		}
		repos, err := client.ListRepositories(context.Background(), ws)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(repos)
		}
		printRepositories(repos)
		return nil

	case "view":
		fs := flag.NewFlagSet("repo view", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 {
			return errors.New("usage: bb repo view <repo> [--workspace slug] [--json]")
		}
		ws, err := resolveWorkspace(cfg, *workspace)
		if err != nil {
			return err
		}
		repo, err := client.GetRepository(context.Background(), ws, fs.Arg(0))
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(repo)
		}
		printRepository(repo)
		return nil

	case "create":
		fs := flag.NewFlagSet("repo create", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		private := fs.Bool("private", true, "Create a private repository")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 {
			return errors.New("usage: bb repo create <repo> [--workspace slug] [--private=true|false] [--json]")
		}
		ws, err := resolveWorkspace(cfg, *workspace)
		if err != nil {
			return err
		}
		repo, err := client.CreateRepository(context.Background(), ws, fs.Arg(0), *private)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(repo)
		}
		printRepository(repo)
		return nil
	default:
		return fmt.Errorf("unknown repo subcommand %q", args[0])
	}
}

func runPR(args []string) error {
	if len(args) == 0 {
		return errors.New("expected pr subcommand: list, view, create")
	}
	cfg, err := requireConfig()
	if err != nil {
		return err
	}
	client, err := bitbucket.NewClient(cfg)
	if err != nil {
		return err
	}

	switch args[0] {
	case "list":
		fs := flag.NewFlagSet("pr list", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		repo := fs.String("repo", "", "Repository slug")
		state := fs.String("state", "OPEN", "PR state")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		ws, repoSlug, err := resolveWorkspaceRepo(cfg, *workspace, *repo)
		if err != nil {
			return err
		}
		prs, err := client.ListPullRequests(context.Background(), ws, repoSlug, *state)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(prs)
		}
		printPullRequests(prs)
		return nil

	case "view":
		fs := flag.NewFlagSet("pr view", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		repo := fs.String("repo", "", "Repository slug")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if fs.NArg() != 1 {
			return errors.New("usage: bb pr view <id> --repo slug [--workspace slug] [--json]")
		}
		id, err := bitbucket.ParsePullRequestID(fs.Arg(0))
		if err != nil {
			return err
		}
		ws, repoSlug, err := resolveWorkspaceRepo(cfg, *workspace, *repo)
		if err != nil {
			return err
		}
		pr, err := client.GetPullRequest(context.Background(), ws, repoSlug, id)
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(pr)
		}
		printPullRequest(pr)
		return nil

	case "create":
		fs := flag.NewFlagSet("pr create", flag.ContinueOnError)
		fs.SetOutput(os.Stdout)
		workspace := fs.String("workspace", "", "Workspace slug")
		repo := fs.String("repo", "", "Repository slug")
		title := fs.String("title", "", "Pull request title")
		description := fs.String("description", "", "Pull request description")
		source := fs.String("source", "", "Source branch")
		destination := fs.String("destination", "", "Destination branch")
		closeSource := fs.Bool("close-source-branch", false, "Close the source branch when merged")
		jsonOut := fs.Bool("json", false, "Print JSON")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if strings.TrimSpace(*title) == "" || strings.TrimSpace(*source) == "" {
			return errors.New("--title and --source are required")
		}
		ws, repoSlug, err := resolveWorkspaceRepo(cfg, *workspace, *repo)
		if err != nil {
			return err
		}
		pr, err := client.CreatePullRequest(context.Background(), ws, repoSlug, bitbucket.CreatePullRequestInput{
			Title:             strings.TrimSpace(*title),
			Description:       *description,
			SourceBranch:      strings.TrimSpace(*source),
			DestinationBranch: strings.TrimSpace(*destination),
			CloseSourceBranch: *closeSource,
		})
		if err != nil {
			return err
		}
		if *jsonOut {
			return printJSON(pr)
		}
		printPullRequest(pr)
		return nil
	default:
		return fmt.Errorf("unknown pr subcommand %q", args[0])
	}
}

func requireConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, errors.New("not authenticated; run `bb auth login` first")
	}
	return cfg, nil
}

func resolveWorkspace(cfg *config.Config, workspace string) (string, error) {
	if strings.TrimSpace(workspace) != "" {
		return strings.TrimSpace(workspace), nil
	}
	if strings.TrimSpace(cfg.Workspace) != "" {
		return strings.TrimSpace(cfg.Workspace), nil
	}
	return "", errors.New("workspace is required; pass --workspace or set a default during `bb auth login`")
}

func resolveWorkspaceRepo(cfg *config.Config, workspace, repo string) (string, string, error) {
	ws, err := resolveWorkspace(cfg, workspace)
	if err != nil {
		return "", "", err
	}
	if strings.TrimSpace(repo) == "" {
		return "", "", errors.New("repository is required; pass --repo")
	}
	return ws, strings.TrimSpace(repo), nil
}

func printUsage() {
	fmt.Println(`bitbucket-cli

Usage:
  bb auth login --email you@example.com --api-token TOKEN [--workspace slug]
  bb auth login --access-token TOKEN [--workspace slug]
  bb auth status
  bb auth logout

  bb repo list [--workspace slug] [--json]
  bb repo view <repo> [--workspace slug] [--json]
  bb repo create <repo> [--workspace slug] [--private=true|false] [--json]

  bb pr list --repo slug [--workspace slug] [--state OPEN] [--json]
  bb pr view <id> --repo slug [--workspace slug] [--json]
  bb pr create --repo slug --title "..." --source branch [--destination branch] [--description text] [--close-source-branch] [--workspace slug] [--json]`)
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func printRepositories(repos []bitbucket.Repository) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVISIBILITY\tMAIN\tURL")
	for _, repo := range repos {
		visibility := "public"
		if repo.IsPrivate {
			visibility = "private"
		}
		mainBranch := ""
		if repo.MainBranch != nil {
			mainBranch = repo.MainBranch.Name
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", repo.FullName, visibility, mainBranch, repo.Links.HTML.Href)
	}
	w.Flush()
}

func printRepository(repo *bitbucket.Repository) {
	visibility := "public"
	if repo.IsPrivate {
		visibility = "private"
	}
	mainBranch := ""
	if repo.MainBranch != nil {
		mainBranch = repo.MainBranch.Name
	}
	fmt.Printf("Name: %s\n", repo.FullName)
	fmt.Printf("Visibility: %s\n", visibility)
	fmt.Printf("SCM: %s\n", repo.SCM)
	if mainBranch != "" {
		fmt.Printf("Main branch: %s\n", mainBranch)
	}
	if repo.Description != "" {
		fmt.Printf("Description: %s\n", repo.Description)
	}
	if repo.Links.HTML.Href != "" {
		fmt.Printf("URL: %s\n", repo.Links.HTML.Href)
	}
}

func printPullRequests(prs []bitbucket.PullRequest) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATE\tSOURCE\tDESTINATION\tTITLE")
	for _, pr := range prs {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", pr.ID, pr.State, branchName(pr.Source), branchName(pr.Destination), pr.Title)
	}
	w.Flush()
}

func printPullRequest(pr *bitbucket.PullRequest) {
	fmt.Printf("ID: %d\n", pr.ID)
	fmt.Printf("Title: %s\n", pr.Title)
	fmt.Printf("State: %s\n", pr.State)
	fmt.Printf("Source: %s\n", branchName(pr.Source))
	fmt.Printf("Destination: %s\n", branchName(pr.Destination))
	if pr.Author != nil {
		fmt.Printf("Author: %s\n", displayUser(&pr.Author.User))
	}
	if pr.Description != "" {
		fmt.Printf("Description: %s\n", pr.Description)
	}
	if pr.Links.HTML.Href != "" {
		fmt.Printf("URL: %s\n", pr.Links.HTML.Href)
	}
}

func branchName(branch *bitbucket.PullRequestBranch) string {
	if branch == nil {
		return ""
	}
	return branch.Branch.Name
}

func displayUser(user *bitbucket.User) string {
	if user == nil {
		return ""
	}
	for _, value := range []string{user.DisplayName, user.Nickname, user.Username, user.AccountID, user.UUID} {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return "unknown user"
}
