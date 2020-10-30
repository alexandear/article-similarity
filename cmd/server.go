package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"

	"github.com/devchallenge/article-similarity/internal/restapi"
	"github.com/devchallenge/article-similarity/internal/util"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			serv, err := restapi.NewArticleServer()
			defer util.Close(serv)
			if err != nil {
				return errors.WithStack(err)
			}

			return serv.Serve()
		},
	}
	cmd.PersistentFlags().AddFlagSet(s.config.Flags())

	return cmd
}
