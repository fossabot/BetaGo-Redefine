package xmodel

type WordMatchType string

const (
	MatchTypeSubStr WordMatchType = "substr"
	MatchTypeRegex  WordMatchType = "regex"
	MatchTypeFull   WordMatchType = "full"
)

type ReplyType string

type ReplyNType struct {
	Reply     string    `json:"reply" gorm:"primaryKey;index"`
	ReplyType ReplyType `json:"reply_type" gorm:"primaryKey;index;default:text"`
}

const (
	ReplyTypeText ReplyType = "text"
	ReplyTypeImg  ReplyType = "img"
)
