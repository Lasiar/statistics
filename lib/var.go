package lib

import (
	"database/sql"
	"github.com/go-redis/redis"
)

var (
	RedisIpDB	*redis.Client
	RedisStatDB *redis.Client
	Config      Configure
	ClickDB     *sql.DB
	PsqlDB		*sql.DB
)
