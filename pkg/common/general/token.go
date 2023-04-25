package general

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func GetTokenFromRequest(g *gin.Context) string {
	token := g.Request.Header.Get("Authorization")

	// normally Authorization the_token_xxx
	strArr := strings.Split(token, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}
