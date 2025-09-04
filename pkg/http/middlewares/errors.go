package middlewares

import (
	"net/http"
	"scheduler/internal/errors"

	"github.com/gin-gonic/gin"
)

func Errors(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err

		if e := errors.GetError(err); e != nil {
			c.JSON(e.Code, gin.H{
				"code":    e.Code,
				"message": e.Msg(),
				"success": false,
			})
			return
		}

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error":   "Internal Server Error",
				"message": "contact the admin",
				"success": false,
			},
		)
	}
}
