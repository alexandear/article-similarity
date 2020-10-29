package cmd

import (
	"github.com/spf13/cobra"
)

type App struct {
}

func (a *App) Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "article-similarity",
		Short: "Article Similarity",
	}
}
