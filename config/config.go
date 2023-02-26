package config

import "github.com/spf13/pflag"

type Config struct {
	MongoURI string
	Database string
	Table    string
}

// RegisterFlags adds the configuration flags to the given flag set.
func (c *Config) RegisterFlags(f *pflag.FlagSet) {
	f.StringVarP(&c.MongoURI, "mongouri", "u",
		"mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/?retryWrites=true",
		"MongoDB URI")
	f.StringVarP(&c.Database, "database", "d", "getircase-study", "Database name")
	f.StringVarP(&c.Table, "table", "t", "records", "Table Name")
}
