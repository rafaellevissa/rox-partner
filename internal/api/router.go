package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/trades/summary", GetTradeSummary(db))
	return r
}
