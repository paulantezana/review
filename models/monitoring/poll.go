package monitoring

import "time"

type Poll struct {
    ID          uint      `json:"id" gorm:"primary_key"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    StartDate   time.Time `json:"start_date"`
    EndDate     time.Time `json:"end_date"`
    Message     string    `json:"message"`
    Weather     string    `json:"weather"`
    State       bool      `json:"state" gorm:"default:'true'"`

    Questions []Question `json:"questions"`
}
