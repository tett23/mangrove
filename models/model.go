package models

import (
	"reflect"
	"regexp"

	"github.com/jinzhu/gorm"
	"github.com/tett23/mangrove/lib/mangrove_db"
)

// TransactionFunc トランザクションの関数
type TransactionFunc func(*DBContext) error

var db *gorm.DB

func FirstOrInit(s interface{}, keys ...interface{}) bool {
	keyNames := primaryKeys(s)
	findValues := make(map[string]interface{})
	for i, k := range keyNames {
		findValues[k] = keys[i]
	}

	return !db.FirstOrInit(s, findValues).RecordNotFound()
}

func FirstOrCreate(s interface{}, keys ...interface{}) error {
	context := currentContext()

	return context.FirstOrCreate(s, keys...)
}

func FindByID(s interface{}, key ...interface{}) bool {
	return !db.First(s, key...).RecordNotFound()
}

func FindByIDs(s interface{}, key ...interface{}) bool {
	// NOTE: PKey名取るためだけにScope使うのと、複合主キーに対応してない辺りがイマイチ
	scope := db.NewScope(s)
	return !db.Find(s, scope.PrimaryKey()+" in(?)", key).RecordNotFound()
}

// Save 保存
func Save(s interface{}) error {
	context := currentContext()

	return context.Save(s)
}

// Create 作成
func Create(s interface{}) error {
	context := currentContext()

	return context.Create(s)
}

// Delete 削除
func Delete(s interface{}) error {
	context := currentContext()

	return context.Delete(s)
}

// Transaction 引数のfuncでトランザクションを開始
func Transaction(f TransactionFunc) error {
	// TODO: 複数のgoroutineがtxを触ると死
	tx := mangrove_db.GetDB().Begin()

	if err := f((*DBContext)(tx)); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func currentContext() *DBContext {
	var context *gorm.DB
	context = db

	return (*DBContext)(context)
}

type DBContext gorm.DB

func (context *DBContext) FirstOrInit(s interface{}, keys ...interface{}) bool {
	return FirstOrInit(s, keys...)
}

func (context *DBContext) FirstOrCreate(s interface{}, keys ...interface{}) error {
	c := (*gorm.DB)(context)

	keyNames := primaryKeys(s)
	findValues := make(map[string]interface{})
	for i, k := range keyNames {
		findValues[k] = keys[i]
	}

	return c.FirstOrCreate(s, findValues).Error
}

func (context *DBContext) FindByID(s interface{}, key interface{}) bool {
	return FindByID(s, key)
}

func (context *DBContext) Save(s interface{}) error {
	c := (*gorm.DB)(context)

	return c.Save(s).Error
}

func (context *DBContext) Create(s interface{}) error {
	c := (*gorm.DB)(context)

	// TODO: bulk insertに対応
	sv := reflect.ValueOf(s)
	switch sv.Kind() {
	case reflect.Slice:
		for i := 0; i < sv.Len(); i++ {
			insertItem := sv.Index(i).Addr().Interface()
			if err := c.Create(insertItem).Error; err != nil {
				return err
			}
		}
	case reflect.Ptr:
		return c.Create(s).Error
	default:
		insertItem := sv.Addr().Interface()
		return c.Create(insertItem).Error
	}

	return nil
}

func (context *DBContext) Delete(s interface{}, where ...interface{}) error {
	c := (*gorm.DB)(context)

	return c.Delete(s, where...).Error
}

func primaryKeys(s interface{}) []string {
	value := reflect.ValueOf(s).Elem()
	nonPointerType := value.Type()
	re, _ := regexp.Compile(`primary_key:\s*?true`)

	length := value.NumField()
	primarykeys := make([]string, 0, 1)
	for i := 0; i < length; i++ {
		field := nonPointerType.Field(i)
		tag := field.Tag.Get("gorm")

		if re.Match([]byte(tag)) {
			primarykeys = append(primarykeys, field.Tag.Get("json"))
		}
	}

	// idはGORMのデフォルトprimary key
	if len(primarykeys) == 0 {
		primarykeys = append(primarykeys, "id")
	}

	return primarykeys
}
