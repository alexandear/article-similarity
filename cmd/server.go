package cmd

import (
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"

	"github.com/devchallenge/article-similarity/internal/restapi"
)

type AppServer struct {
	*App
	config Config
}

func (a *App) AppServer() cmder.Cmder {
	return &AppServer{App: a}
}

func (s *AppServer) Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP server",
		RunE:  restapi.RunEArticleServer,
	}
	cmd.PersistentFlags().AddFlagSet(s.config.Flags())

	return cmd
}
