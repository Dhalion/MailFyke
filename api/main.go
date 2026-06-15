package main

import (
	"github.com/chris/MailFyke/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	cfg := config.Load()

	root := &cobra.Command{Use: "mailfyke"}
	root.AddCommand(
		newServerCmd(cfg),
		newMigrateCmd(cfg),
	)
	_ = root.Execute()
}

func newServerCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start the HTTP and SMTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(cfg)
		},
	}
}

func newMigrateCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMigrate(cfg)
		},
	}
}
