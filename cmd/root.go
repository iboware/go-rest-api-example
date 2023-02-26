package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/iboware/go-rest-api-example/config"
	"github.com/iboware/go-rest-api-example/db"
	"github.com/iboware/go-rest-api-example/pkg/handler"
	"github.com/iboware/go-rest-api-example/pkg/store"
	"github.com/spf13/cobra"
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
		router := mux.NewRouter()
		router.HandleFunc("/in-memory", kvHandler.Fetch).Methods("GET")
		router.HandleFunc("/in-memory", kvHandler.Create).Methods("POST")
		router.HandleFunc("/mdb", mdbHandler.Fetch).Methods("GET")

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
