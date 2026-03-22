package cmd

import "bitbucket-cli/internal/app"

func Execute(args []string) error {
	return app.Run(args)
}
