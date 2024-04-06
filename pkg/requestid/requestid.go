package requestid

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

// Set sets request ID for request, got by context
func Set(ctx *gin.Context) {
	rid := ctx.GetHeader(HeaderXRequestID)
	if rid == "" {
		rid = uuid.NewString()
		ctx.Request.Header.Add(HeaderXRequestID, rid)
	}

	ctx.Header(HeaderXRequestID, rid)
	ctx.Next()
}

// Get returns the request's ID
func Get(c *gin.Context) string {
	return c.Writer.Header().Get(HeaderXRequestID)
}
