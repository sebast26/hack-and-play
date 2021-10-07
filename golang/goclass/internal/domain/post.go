package domain

import "fmt"

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (p Post) String() string {
	return fmt.Sprintf("this post has userid=%d id=%d title=%s body=%s\n", p.UserID, p.ID, p.Title[:12], p.Body[:12])
}
