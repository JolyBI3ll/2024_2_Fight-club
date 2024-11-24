package domain

import (
	"context"
	"time"
)

type Chat struct {
	LastMessage  string    `gorm:"type:text;size:1000;column:lastMessage" json:"lastMessage"`
	LastDate     time.Time `gorm:"type:timestamp;column:lastDate" json:"lastDate"`
	AuthorName   string    `gorm:"type:text;size:255;column:authorName" json:"authorName"`
	AuthorAvatar string    `gorm:"type:text;size:255;column:authorAvatar" json:"authorAvatar"`
	AuthorUUID   string    `gorm:"column:authorUuid;not null" json:"authorUuid"`
}

type Message struct {
	ID         int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	SenderID   string    `gorm:"column:senderId;not null" json:"senderId"`
	ReceiverID string    `gorm:"column:receiverId;not null" json:"receiverId"`
	Content    string    `gorm:"type:text;size:1000;column:content" json:"content"`
	CreatedAt  time.Time `gorm:"type:timestamp;column:createdAt;default:CURRENT_TIMESTAMP" json:"createdAt"`
	Sender     User      `gorm:"foreignkey:SenderID;references:UUID" json:"-"`
	Receiver   User      `gorm:"foreignkey:ReceiverID;references:UUID" json:"-"`
}

type ChatRepository interface {
	GetChats(ctx context.Context, userID string, lastUpdateTime time.Time) ([]*Chat, error)
	SendNewMessage(ctx context.Context, receiver string, sender string, message string) error
	GetMessages(ctx context.Context, userID1 string, userID2 string, lastSentTime time.Time) ([]*Message, error)
}
