package models

import (
	"time"

	"github.com/tett23/mangrove/lib/mangrove_db"
)

// Program 番組情報
type Program struct {
	ID          int       `json:"id"`
	Count       int       `json:"int"`
	EpisodeName string    `json:"episode_name"`
	ChannelName string    `json:"channel_name"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	Title       string    `json:"title"`
	ChannelID   int       `json:"channel_id"`
}
type Programs []Program

// Search 番組情報を検索
func (p *Program) Search(channelName string, startAt, endAt time.Time) bool {
	return !mangrove_db.GetDB().Where("start_at = ? and channel_name = ?", startAt, channelName).First(p).RecordNotFound()

}
