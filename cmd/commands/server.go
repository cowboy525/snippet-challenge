package commands

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/topoface/snippet-challenge/api"
	"github.com/topoface/snippet-challenge/app"
	"github.com/topoface/snippet-challenge/config"
	"github.com/topoface/snippet-challenge/mlog"
	"github.com/topoface/snippet-challenge/viper"
	"github.com/topoface/snippet-challenge/web"
)

var serverCmd = &cobra.Command{
	Use:          "serve",
	Short:        "Run the ErniePJT server",
	RunE:         serverCmdF,
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(serverCmd)
	RootCmd.RunE = serverCmdF
}

func serverCmdF(command *cobra.Command, args []string) error {
	configDSN := viper.GetString("config")

	disableConfigWatch, _ := command.Flags().GetBool("disableconfigwatch")

	interruptChan := make(chan os.Signal, 1)

	configStore, err := config.NewStore(configDSN, !disableConfigWatch)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}

	return runServer(configStore, interruptChan)
}

func runServer(configStore config.Store, interruptChan chan os.Signal) error {
	options := []app.Option{
		app.ConfigStore(configStore),
	}
	server, err := app.NewServer(options...)
	if err != nil {
		mlog.Critical(err.Error())
		return err
	}
	defer server.Shutdown()

	api.Init(server, server.AppOptions, server.Router)
	web.New(server, server.AppOptions, server.Router)

	serverErr := server.Start()
	if serverErr != nil {
		mlog.Critical(serverErr.Error())
		return serverErr
	}

	// wait for kill signal before attempting to gracefully shutdown
	// the running service
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan

	return nil
}
