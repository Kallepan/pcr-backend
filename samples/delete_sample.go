package samples

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/kaka/pcr-backend/common/database"
)

func DeleteSample(ctx *gin.Context) {
	tagesnummer := ctx.Param("tagesnummer")

	query := `DELETE FROM samples WHERE tagesnummer = $1`

	_, err := database.Instance.Exec(query, tagesnummer)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
