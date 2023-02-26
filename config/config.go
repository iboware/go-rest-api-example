package config

import (
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type Config struct {
	MongoURI string
	Database string
	Table    string
	Port     int
}

// RegisterFlags adds the configuration flags to the given flag set.
func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.MongoURI, "mongouri", "u",
		"mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/?retryWrites=true",
		"MongoDB URI")
	f.StringVarP(&c.Database, "database", "d", "getircase-study", "Database name")
	f.StringVarP(&c.Table, "table", "t", "records", "Table Name")
	f.IntVarP(&c.Port, "port", "p", 8000, "Port")

	// this is to set the port environmental variable provided by heroku.
	port := os.Getenv("PORT")
	if port != "" {
		var err error
		c.Port, err = strconv.Atoi(port)
		if err != nil {
			panic("failed to parse port value")
		}
	}
}
