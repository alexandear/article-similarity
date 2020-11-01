package model

type Article struct {
	ID           int
	Content      string
	DuplicateIDs []int
}
