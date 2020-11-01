package model

type Article struct {
	ID           int
	Content      string
	DuplicateIDs []int
}

type DuplicateGroup struct {
	IDs []int
}
