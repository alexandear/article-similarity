package model

type Article struct {
	ID           int
	Content      string
	DuplicateIDs []int
	IsUnique     bool
}

type DuplicateGroup struct {
	IDs []int
}
