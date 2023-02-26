package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iboware/go-rest-api-example/config"
	"github.com/iboware/go-rest-api-example/db"
	"github.com/iboware/go-rest-api-example/pkg/handler"
	"github.com/iboware/go-rest-api-example/pkg/store"
	"github.com/spf13/cobra"
)

var cfg = config.Config{}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		err := http.ListenAndServe(":8000", router)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	cfg.RegisterFlags(serverCmd.Flags())
}
