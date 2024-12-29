package ginmethods

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ReadAndPrintBody(c *gin.Context) ([]byte, error) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read body"})
		return nil, err
	}
	fmt.Println("Request Body:", string(bodyBytes))
	return bodyBytes, nil
}
