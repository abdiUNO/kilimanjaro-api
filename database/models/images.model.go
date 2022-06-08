package models

import (
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
)

type Image struct {
	orm.GormModel
	Url string `json:"url"`
}

func (image *Image) Create() (*Image, *utils.Error) {
	err := GetDB().Create(&image).Error

	if err != nil {
		log.Println("RRRR")
		return &Image{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return image, nil
}