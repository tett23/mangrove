package models

import (
	"github.com/tett23/mangrove/lib/mangrove_db"
)

type Video struct {
	ID           int    `json:"id" gorm:"primary_key:true"`
	Name         string `json:"name"`
	OutputName   string `json:"output_name"`
	OriginalName string `json:"original_name"`
	Program      string `json:"program" gorm:"type:text"`
}
type Videos []Video

func (videos *Videos) Latest() error {
	err := mangrove_db.GetDB().Order("created_at desc").Limit(1000).Find(videos).Error

	return err
}
