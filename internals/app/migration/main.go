package mg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"github.com/vapusdata-ecosystem/vapusai/core/app/migration/migrations"
	apppkgs "github.com/vapusdata-ecosystem/vapusai/core/app/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/core/data-platform/connectors/databases"
	"github.com/vapusdata-ecosystem/vapusai/core/models"
	llm_observability "github.com/vapusdata-ecosystem/vapusai/core/observability/llm"
	// test "github.com/vapusdata-ecosystem/vapusai/examples/datamesh"
)

func LoadMigrate(ctx context.Context, sourceCreds *models.DataSourceCredsParams, logger zerolog.Logger) error {

	// Initialize SQLite database connection
	// logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	// ctx := context.Background()
	// fmt.Println("Username: ", sourceCreds.DataSourceCreds.GenericCredentialModel.Username)
	// fmt.Println("Password: ", sourceCreds.DataSourceCreds.GenericCredentialModel.Password)
	// fmt.Println("DB: ", sourceCreds.DataSourceCreds.DB)
	// fmt.Println("URL: ", sourceCreds.DataSourceCreds.URL)
	// fmt.Println("Port: ", sourceCreds.DataSourceCreds.Port)
	// fmt.Println("DataSourceEngine: ", sourceCreds.DataSourceEngine)
	// fmt.Println("DataStorelist: ", sourceCreds.DataStorelist)
	client, err := databases.New(
		ctx, databases.WithInApp(true), databases.WithLogger(logger),
		databases.WithDataSourceCredsParams(&models.DataSourceCredsParams{
			DataSourceCreds: &models.DataSourceSecrets{
				GenericCredentialModel: &models.GenericCredentialModel{
					Username: sourceCreds.DataSourceCreds.GenericCredentialModel.Username,
					Password: sourceCreds.DataSourceCreds.GenericCredentialModel.Password,
				},
				DB:      sourceCreds.DataSourceCreds.DB,
				URL:     sourceCreds.DataSourceCreds.URL,
				Port:    int64(sourceCreds.DataSourceCreds.Port),
				Version: "",
			},
			DataSourceEngine:  sourceCreds.DataSourceEngine,
			DataSourceType:    sourceCreds.DataSourceType,
			DataSourceService: sourceCreds.DataSourceService,
			DataStorelist:     sourceCreds.DataStorelist,
		}),
	)

	if err != nil {
		logger.Err(err).Msgf("Error while connecting to the Metadata")
		return err
	}
	// fmt.Printf("I am calling Postgress Connector for Migration --->>> %v", client)

	// client.RunMigrateClient()

	// client.AddQueryHook(bundebug.NewQueryHook(
	// 	bundebug.WithEnabled(false),
	// 	bundebug.FromEnv(),
	// ))
	// fmt.Println("DB instance: ", client.PostgresClient.DB)
	dmstore := &apppkgs.VapusStore{
		BeDataStore: &apppkgs.BeDataStore{
			Db: client,
		},
	}
	app := &cli.App{
		Name: "migrate",
		Commands: []*cli.Command{
			NewDBCommand(dmstore, migrate.NewMigrator(client.PostgresClient.DB, migrations.Migrations), logger),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func NewDBCommand(dmstore *apppkgs.VapusStore, migrator *migrate.Migrator, logger zerolog.Logger) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "database migrations",
		Subcommands: []*cli.Command{

			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context) //nolint:errcheck

					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no new migrations to run (database is up to date)\n")
						return nil
					}
					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					if err := migrator.Lock(c.Context); err != nil {
						return err
					}
					defer migrator.Unlock(c.Context) //nolint:errcheck

					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}
					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					return migrator.Lock(c.Context)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					return migrator.Unlock(c.Context)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(c.Context, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "create_tx_sql",
				Usage: "create up and down transactional SQL migrations",
				Action: func(c *cli.Context) error {
					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateTxSQLMigrations(c.Context, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created transaction migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ms, err := migrator.MigrationsWithStatus(c.Context)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())
					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					group, err := migrator.Migrate(c.Context, migrate.WithNopMigration())
					if err != nil {
						return err
					}
					if group.IsZero() {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}
					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
			{
				Name:  "update_ai_pricing",
				Usage: "Update AI pricing in the database",
				Action: func(c *cli.Context) error {
					// Call your Go function
					if err := llm_observability.UpdateAllAIModelPrices(context.Background(), dmstore, logger); err != nil {
						return fmt.Errorf("failed to update AI pricing: %w", err)
					}
					fmt.Println("AI pricing update completed successfully.")
					return nil
				},
			},
		},
	}
}
