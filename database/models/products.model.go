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
	Category    Category        `json:"category"`
	CategoryId  string          `json:"categoryId"`
	Price       decimal.Decimal `json:"price" gorm:"type:numeric"`
	ImageOne    string          `json:"imageOne"`
	ImageTwo    string          `json:"imageTwo"`
	ImageThree  string          `json:"imageThree"`
	Images      []string        `json:"images" gorm:"-"`
	Vendor      *User           `json:"vendor"`
	VendorID    string          `json:"vendorId"`
}

func (product *Product) AfterFind(tx *gorm.DB) (err error) {
	if product.ImageOne != "" {
		product.Images = append(product.Images, product.ImageOne)
	}
	if product.ImageTwo != "" {
		product.Images = append(product.Images, product.ImageTwo)
	}
	if product.ImageThree != "" {
		product.Images = append(product.Images, product.ImageThree)
	}

	return
}

func (product *Product) Create() (*Product, *utils.Error) {
	err := GetDB().Create(&product).Error

	if product.ImageOne != "" {
		product.Images = append(product.Images, product.ImageOne)
	}
	if product.ImageTwo != "" {
		product.Images = append(product.Images, product.ImageTwo)
	}
	if product.ImageThree != "" {
		product.Images = append(product.Images, product.ImageThree)
	}

	category, _ := FindCategoryById(product.CategoryId)
	product.Category = *category

	if err != nil {
		log.Println("RRRR")
		return &Product{}, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return product, nil
}

func FindProductById(productId string) (*Product, *utils.Error) {
	product := &Product{}
	err := GetDB().Table("Products").Preload("Categories").Where("id = ?", productId).Find(&product).Error

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

func FindProductsByVendor(userId string) (*[]Product, *utils.Error) {
	product := &[]Product{}
	err := GetDB().Table("products").Preload("Vendor").Preload("Category").Find(&product, "vendor_id = ?", userId).Error

	if err != nil {
		log.Println(err)
		return product, utils.NewError(utils.EINTERNAL, "internal database error", err)
	}

	return product, nil
}

func FindAllProducts() (*[]Product, *utils.Error) {

	products := &[]Product{}

	if err := GetDB().Table("products").Preload("Vendor").Preload("Category").Order("created_at desc").Find(&products).Error; err != nil {
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
