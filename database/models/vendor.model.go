package models

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
)

func GetDB() *gorm.DB {
	return orm.DBCon
}

type Vendor struct {
	orm.GormModel
	Title       string    `json:"title"`
	Description string    `sql:"type:longtext" json:"description"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Image       string    `json:"image"`
	Catalog     []Product `json:"catalog"`
}

func (vendor *Vendor) Create() (*Vendor, *utils.Error) {
	err := GetDB().Create(&vendor).Error
	if err != nil {
		return &Vendor{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return vendor, nil
}

func FindAllVendors() (*[]Vendor, *utils.Error) {

	vendors := &[]Vendor{}

	if err := GetDB().Limit(10).Find(&vendors).Error; err != nil {
		log.Println(err)
		return vendors, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return vendors, nil
}

func FindVendorById(vendorID string) (*Vendor, *utils.Error) {

	vendor := &Vendor{}

	if err := GetDB().First(&vendor, "id = ?", vendorID).Error; err != nil {
		log.Println(err)
		return vendor, utils.NewError(utils.EINTERNAL, "vendor record not found", err)
	}

	return vendor, nil
}

func GetVendorByName(title string) *Vendor {
	vendor := &Vendor{}
	err := GetDB().Table("vendors").Where("title = ?", title).First(vendor).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}

	return vendor
}
