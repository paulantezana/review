package messengermodel

type MessageRecipient struct {
    ID       uint   `json:"id" gorm:"primary_key"`
    IsRead bool `json:"is_read"`

    RecipientID uint `json:"recipient_id"`
    RecipientGroupID uint `json:"recipient_group_id"`
    MessageID uint `json:"message_id"`
}
