package cachedmysql

import (
	"database/sql"

	"github.com/go-redis/redis/v9"
)

type (
	Config struct {
		DB  DBConfig
		RDS RDSConfig
	}

	DBConfig struct {
		Driver     string
		DataSource string
	}

	RDSConfig redis.Options
)

func New(cfg *Config, opts ...Option) *DB {
	o := newOptions(opts...)
	db, err := sql.Open(cfg.DB.Driver, cfg.DB.DataSource)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &DB{
		DB: db,
		Cache: &cache{
			rds:    redis.NewClient((*redis.Options)(&cfg.RDS)),
			expiry: o.Expiry,
		},
	}
}
