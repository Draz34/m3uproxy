package db

/*
	"category_id": "101",
    "category_name": "VOD FR",
    "parent_id": 0
*/
type MovieCategorie struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentID     int    `json:"parent_id"`
}
