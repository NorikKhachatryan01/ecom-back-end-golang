package models

import (
	"github.com/NorikKhachatryan01/go-ecommerce/pkg/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Product struct{
	gorm.Model
	ID          string     `gorm:json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       int        `json:"price"`
	Category    string  `json:"category"`
   }
   
   type Category struct{
	ID 	    string 	`json:"id"`
	Type    string  `json:"type"`
   }
   
   func init(){
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Product{})
   }

   func (p *Product) CreateProduct() *Product{
	db.NewRecord(p)
	db.Create(&p)
	return p
   }

   func GetAllProducts() []Product{
	var Products []Product
	db.Find(&Products)
	return Products
   }

   func GetProductById(Id int64)(*Product, *gorm.DB){
	var getProduct Product
	db:=db.Where("ID=?", Id).Find(&getProduct)
	return &getProduct,db
   }

   func DeleteProduct(ID int64) Product{
	var product Product
	db.Where("ID=?", ID).Delete(product)
	return product
   }