package main

import(
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

type HammerServer struct {
	router *gin.Engine
	db *xorm.Engine
}

var hs *HammerServer

func main() {
	defer func() {
		if p := recover();p != nil {
			fmt.Println(p)
		}
	}()

	hs = new(HammerServer)
	hs.router = gin.Default()
	hs.NewRouter()

	hs.router.Run("0.0.0.0:8001")
}


