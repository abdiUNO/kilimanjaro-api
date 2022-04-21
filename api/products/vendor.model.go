package products

import (
	"github.com/jinzhu/gorm"
	"kilimanjaro-api/database/orm"
)

type Vendor struct {
	orm.GormModel
	Title       string `json:"title"`
	Description string `sql:"type:longtext"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Image       string `json:"Image"`
}

func (vendor *Vendor) Create() error {
	err := GetDB().Create(&vendor).Error
	return err
}

func GetVendorByName(title string) *Vendor {
	vendor := &Vendor{}
	err := GetDB().Table("vendors").Where("title = ?", title).First(vendor).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}

	return vendor
}
