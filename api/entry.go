package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/newtoallofthis123/sqlite_http/db"
)

type Api struct {
	ListenAddr string
	Db         *db.Db
}

func NewApi(listenAddr string) *Api {
	d, err := db.NewDb("test.db")
	if err != nil {
		panic(err)
	}
	err = d.Discover()
	if err != nil {
		panic(err)
	}

	return &Api{ListenAddr: listenAddr, Db: d}
}

func (a *Api) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (a *Api) handleCommandSend(c *gin.Context) {
	cmd := c.PostForm("cmd")

	if cmd == "" {
		c.JSON(400, gin.H{
			"error": "no command provided",
		})
		return
	}

	if strings.HasPrefix(cmd, "SELECT") {
		rows, err := a.Db.RunQuery(cmd)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"rows": rows,
		})
		return
	} else {
		r, err := a.Db.RunExec(cmd)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		rowsEffected, _ := r.RowsAffected()

		c.JSON(200, gin.H{
			"message": rowsEffected,
		})
		return
	}
}

func (a *Api) handleTables(c *gin.Context) {
	c.JSON(200, gin.H{
		"tables": a.Db.RowsInfo,
	})
}

func (a *Api) Run() {
	r := gin.Default()

	r.GET("/ping", a.ping)
	r.POST("/send", a.handleCommandSend)
	r.GET("/tables", a.handleTables)

	r.Run(a.ListenAddr)
}
