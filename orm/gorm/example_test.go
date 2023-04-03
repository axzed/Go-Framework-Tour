package gorm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

type Product struct {
	gorm.Model
	Code  string `gorm:"column(code)"`
	Price uint
}

func  (p Product) TableName() string {
	return "product_t"
}

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	println("before save")
	return
}

func (p *Product) AfterSave(tx *gorm.DB) (err error) {
	println("after save")
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	println("before create")
	return
}

func (p *Product) AfterCreate(tx *gorm.DB) (err error) {
	println("after create")
	// 刷新缓存
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	println("before update")
	return
}

func (p *Product) AfterUpdate(tx *gorm.DB) (err error) {
	println("after update")
	// 刷新缓存
	return
}

func (p *Product) BeforeDelete(tx *gorm.DB) (err error) {
	println("before update")
	return
}

func (p *Product) AfterDelete(tx *gorm.DB) (err error) {
	println("after update")
	return
}

func  (p *Product) AfterFind(tx *gorm.DB) (err error) {
	println("after find")
	return
}

func TestCRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 打印 SQL，但不执行
	db.DryRun = true

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})


	// Read
	var product Product
	db.First(&product, 1) // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)
}
