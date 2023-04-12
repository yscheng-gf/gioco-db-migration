package config

type Conf struct {
	Postgres PostgresConf `mapstructure:"postgres"`
	Mongo    MongoConf    `mapstructure:"mongo"`
}
type PostgresConf struct {
	Host     string
	Port     uint64
	User     string
	Password string
	DBName   string
	SslMode  string
}

type MongoConf struct {
	Host string
}
