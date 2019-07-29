package utilities

import (
	"os"
	"time"
)

type File struct {
	Name string    `json:"name"`
	Size uint      `json:"size"`
	Date time.Time `json:"date"`
}

// Exists reports whether the named file or directory exists.
func FileExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
