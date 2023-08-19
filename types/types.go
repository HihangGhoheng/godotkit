package gdk_types

type PostgreSQLConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	SSLMode      string
}

type MongoConfig struct {
	Dsn                   string
	DatabaseName          string
	ReadPreference        string
	MinPoolSize           int
	MaxPoolSize           int
	MaxConnectionIdleTime int
}
