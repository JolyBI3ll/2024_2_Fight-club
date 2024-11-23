package domain

type Survey struct {
	ID    int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Title string `gorm:"type:text;size:3000;column:title" json:"title"`
}
