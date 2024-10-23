package api

import "github.com/gin-gonic/gin"

type Api struct {
	ListenAddr string
}

func NewApi(listenAddr string) *Api {
	return &Api{ListenAddr: listenAddr}
}

func (a *Api) ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (a *Api) Run() {
	r := gin.Default()

	r.GET("/ping", a.ping)

	r.Run(a.ListenAddr)
}
