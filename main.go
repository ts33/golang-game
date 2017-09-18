package main


import (
	"fmt"
	"github.com/gin-gonic/gin"
)

const port int = 3000

func main() {

	// Creates a gin router with default middleware
	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	router.GET("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{ "message": "in get"} )
	})

	router.POST("/users", func(c *gin.Context){
		c.JSON(200, gin.H{ "message": "in post"} )
	})

	router.GET("/items", func(c *gin.Context) {
		c.JSON(200, gin.H{ "message": "in get"} )
	})

	router.POST("/items", func(c *gin.Context){
		c.JSON(200, gin.H{ "message": "in post"} )
	})

	router.DELETE("/all", func(c *gin.Context){
		c.JSON(200, gin.H{ "message": "in delete"} )
	})


	router.Run(fmt.Sprintf(":%d", port))
}
