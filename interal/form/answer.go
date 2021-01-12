package form


type Answer struct {
	Content string `form:"content" binding:"required"`
}
