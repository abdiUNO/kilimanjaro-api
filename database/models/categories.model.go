package models

import (
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
)

type Category struct {
	orm.GormModel
	Title    string    `json:"title"`
	Image    string    `json:"image"`
	Products []Product `json:"products" gorm:"foreignKey:CategoryId"`
}

func (category *Category) Create() (*Category, *utils.Error) {
	err := GetDB().Create(&category).Error

	if err != nil {
		return &Category{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return category, nil
}
func FindCategoryById(categoryId string) (*Category, *utils.Error) {

	category := &Category{}

	if err := GetDB().Table("Categories").Preload("Products.Vendor").Where("id = ?", categoryId).Find(&category).Error; err != nil {
		log.Println(err)
		return category, utils.NewError(utils.EINTERNAL, "category record not found", err)
	}

	return category, nil
}

func FindAllCategories() (*[]Category, *utils.Error) {

	categories := &[]Category{}

	if err := GetDB().Limit(10).Find(&categories).Error; err != nil {
		log.Println(err)
		return categories, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return categories, nil
}
