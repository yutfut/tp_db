package forum

import "github.com/fasthttp/router"

func SetForumRouting(r *router.Router, h *Handlers) {
	r.POST("/api/forum/create", h.CreateForum)
	r.GET("/api/forum/{slug}/details", h.GetForumDetails)
	r.POST("/api/forum/{slug}/create", h.CreateThread)
	r.GET("/api/forum/{slug}/threads", h.GetForumThreads)
	r.GET("/api/forum/{slug}/users", h.ForumUsers)
}
