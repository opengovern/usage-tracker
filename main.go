package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opengovern/og-util/pkg/httpserver"
	"github.com/opengovern/og-util/pkg/koanf"
	"github.com/opengovern/og-util/pkg/postgres"
	"github.com/opengovern/usage-tracker/api"
	"github.com/opengovern/usage-tracker/config"
	"github.com/opengovern/usage-tracker/db/model"
	"github.com/opengovern/usage-tracker/db/repo"
	"github.com/opengovern/usage-tracker/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)


func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := Command().ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Command() *cobra.Command {
	cnf := koanf.Provide("information", config.InformationConfig{})

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			logger, err := zap.NewProduction()
			if err != nil {
				return err
			}

			logger = logger.Named("information")

			cmd.SilenceUsage = true

			db, err := postgres.NewClient(&postgres.Config{
				Host:    cnf.Postgres.Host,
				Port:    cnf.Postgres.Port,
				User:    cnf.Postgres.Username,
				Passwd:  cnf.Postgres.Password,
				DB:      cnf.Postgres.DB,
				SSLMode: cnf.Postgres.SSLMode,
			}, logger)
			if err != nil {
				return err
			}
			// create citext extension if not exists
			err = db.Exec("CREATE EXTENSION IF NOT EXISTS citext").Error
			if err != nil {
				logger.Error("failed to create citext extension", zap.Error(err))
				return err
			}
			err = db.AutoMigrate(&model.CspmUsage{})

			cspmUsageRepo := repo.NewCspmUsageRepo(db)

			informationService := service.NewInformationService(cnf, logger, cspmUsageRepo)
			return httpserver.RegisterAndStart(
				ctx,
				logger,
				cnf.Http.Address,
				api.New(cnf, logger, informationService),
			)
		},
	}

	return cmd
}
