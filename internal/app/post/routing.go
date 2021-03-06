package post

import "github.com/fasthttp/router"

func SetPostRouting(r *router.Router, h *Handlers) {
	r.GET("/api/post/{id}/details", h.PostDetails)
	r.POST("/api/post/{id}/details", h.UpdatePost)
}
