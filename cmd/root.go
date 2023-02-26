package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/iboware/go-rest-api-example/config"
	"github.com/iboware/go-rest-api-example/db"
	_ "github.com/iboware/go-rest-api-example/docs"
	"github.com/iboware/go-rest-api-example/pkg/handler"
	"github.com/iboware/go-rest-api-example/pkg/store"
	"github.com/spf13/cobra"
	httpSwagger "github.com/swaggo/http-swagger"
)

var cfg = config.Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-rest-api-example",
	Short: "Starts rest api server",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("starting server")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Initialize Mongo Client
		mongoClient := db.NewMongoDBConnection(ctx, cfg.MongoURI)
		defer mongoClient.Disconnect(ctx)

		// Initialize Mongo Store
		mongoStore := store.NewMongoStore(mongoClient, cfg.Database, cfg.Table)

		// Initialize handlers
		kvHandler := handler.NewKeyValueHandler()
		mdbHandler := handler.NewMDBHandler(ctx, mongoStore)

		// Map handlers to routers
		router := chi.NewRouter()
		router.Get("/in-memory", kvHandler.Fetch)
		router.Post("/in-memory", kvHandler.Create)
		router.Get("/mdb", mdbHandler.Fetch)
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf(":%d/swagger/doc.json", cfg.Port)),
		))

		err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), router)
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cfg.RegisterFlags(rootCmd.Flags())
}
