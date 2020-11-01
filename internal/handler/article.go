package handler

type Article struct {
	ID      int    `json:"id" bson:"id"`
	Content string `json:"content" bson:"content"`
}
