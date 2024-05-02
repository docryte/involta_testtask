package app

import (
	"main/internal/config"
	"main/internal/database"
	"main/internal/transport"
)

func Run() {
	cfg := config.GetEnvironment()
	reindexerDB := database.ConnectReindexer(cfg.REINDEXER_URL)
	database.ConnectRedis(cfg.REDIS_URL, cfg.REDIS_PASSWORD, cfg.REDIS_DATABASE)
	transport.StartRouting(reindexerDB)
}