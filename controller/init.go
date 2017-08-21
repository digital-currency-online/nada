package controller

import (
	"log"
	"nada/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result map[string]interface{}

var (
	//Server web server engine - *gin.Engine
	Server = gin.New()
	errs   = make(map[string]Result, 0)
)

func init() {
	// Server = gin.New()
	if Server == nil {
		log.Fatalln("init web server error")
		return
	}
	// Global middleware
	Server.Use(gin.Logger())
	Server.Use(gin.Recovery())

}

func NewResult() (r Result) {
	r = make(map[string]interface{}, 0)
	r["Ok"] = false
	r["Err"] = ""
	return
}

func getToken() gin.HandlerFunc {
	tn := core.GlobalConfig.GetTokenName()
	return func(c *gin.Context) {
		tv := c.GetHeader(tn)
		if tv != "" {
			c.Set(core.DefaultInternalTokenName, tv)
			c.Next()
			return
		}
		var has bool
		if tv, has = c.GetQuery(tn); has {
			c.Set(core.DefaultInternalTokenName, tv)
			c.Next()
			return
		}
		r := NewResult()
		r["Err"] = "no token"
		c.AbortWithStatusJSON(http.StatusUnauthorized, r)
	}
}

func Run() {
	if Server == nil {
		log.Fatalln("web server not running")
		return
	}
	Server.StaticFile("/test/jwt", "./test/jwt.html")
	Server.Run(core.GlobalConfig.GetServerAddr())
}