package messengermodel

type ReminderFrequency struct {
    ID       uint   `json:"id" gorm:"primary_key"`
    Name string `json:"name"`
    Frequency float32 `json:"frequency"`
    IsActive bool `json:"is_active"`
}
