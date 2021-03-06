package response

type CreatedItem struct {
	Message string `json:"message" binding:"required"`
	Item    *Item  `json:"item" binding:"required"`
}

type Item struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required,max=255"`
}
