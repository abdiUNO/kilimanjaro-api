package models

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"kilimanjaro-api/database/orm"
	"kilimanjaro-api/utils"
)

//
//var redisClient = redis.NewClient(&redis.Options{
//	Addr:     "localhost:6379",
//	Password: "", // no password set
//	DB:       0,  // use default DB
//})

type Product struct {
	orm.GormModel
	Title       string          `json:"title"`
	Description string          `sql:"type:longtext" json:"description"`
	Category    string          `json:"category"`
	Price       decimal.Decimal `json:"price" gorm:"type:numeric"`
	Image       string          `json:"image"`
	Vendor      *Vendor         `json:"vendor"`
	VendorID    string          `json:"vendorId"`
}

func (product *Product) Create() (*Product, *utils.Error) {
	err := GetDB().Create(&product).Error

	if err != nil {
		log.Println("RRRR")
		return &Product{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return product, nil
}

func FindProductById(productId string) (*Product, *utils.Error) {
	product := &Product{}
	err := GetDB().First(&product, "id = ?", productId).Error

	if err != nil {
		log.Println(err)
		if err == gorm.ErrRecordNotFound {
			return &Product{}, utils.NewError(utils.ENOTFOUND, "products not found", nil)
		}
		return &Product{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	if err != nil {
		log.Println(err)
		return product, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return product, nil
}

func FindAllProducts() (*[]Product, *utils.Error) {

	products := &[]Product{}

	if err := GetDB().Limit(10).Table("products").Preload("Vendor").Find(&products).Error; err != nil {
		log.Println(err)
		return products, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return products, nil
}

func QueryProduct(query string) (*[]Product, *utils.Error) {
	products := &[]Product{}

	//SELECT * from products where MATCH(name) AGAINST('Radio' IN NATURAL LANGUAGE MODE)
	log.Println(query)
	if err := GetDB().Raw(`
		SELECT
			*,
			MATCH(title) AGAINST (? IN BOOLEAN MODE) AS score
		FROM
			products
		WHERE
			MATCH(title) AGAINST (? IN BOOLEAN MODE) > 0
	`, query, query).Scan(&products).Error; err != nil {
		log.Println(err)
		return products, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return products, nil
}
