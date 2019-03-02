package utilities

import "time"

type File struct {
	Name string    `json:"name"`
	Size uint      `json:"size"`
	Date time.Time `json:"date"`
}
