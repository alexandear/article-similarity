package cmd

import (
	"github.com/go-openapi/loads"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	cmder "github.com/yaegashi/cobra-cmder"

	"github.com/devchallenge/article-similarity/internal/restapi"
	"github.com/devchallenge/article-similarity/internal/restapi/operations"
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
			swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
			if err != nil {
				return errors.Wrap(err, "failed to embedded spec")
			}
			cmd.Long = swaggerSpec.Spec().Info.Description

			api := operations.NewArticleSimilarityAPIAPI(swaggerSpec)
			serv := server{Server: restapi.NewServer(api)}
			defer util.Close(serv)

			if err := serv.Serve(); err != nil {
				return errors.WithStack(err)
			}

			return nil
		},
	}
	cmd.PersistentFlags().AddFlagSet(s.config.Flags())

	return cmd
}

type server struct {
	*restapi.Server
}

func (s server) Close() error {
	return s.Shutdown()
}
