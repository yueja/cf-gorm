package db

import (
	"errors"
	"gorm.io/gorm"
)

type Curd struct {
	db *gorm.DB
}

// DbCurd 获取指定数据库和数据行的操作实例会话
func DbCurd() *Curd {
	return &Curd{
		db: GetClient(nil),
	}
}

// SpecifyDbCurd 指定客户端获取指定数据库和数据行的操作实例会话 client：客户端名称
func SpecifyDbCurd(client string) *Curd {
	return &Curd{
		db: GetClient([]string{client}),
	}
}

// Save 插入数据，rowsAffected:返回插入记录的条数
func (c Curd) Save(table string, data interface{}) (rowsAffected int64, err error) {
	tx := c.db
	if table != "" {
		tx = tx.Table(table)
	}
	result := tx.Save(data)
	return result.RowsAffected, result.Error
}

// CreateInBatches 批量插入数据 rowsAffected:返回插入记录的条数
func (c Curd) CreateInBatches(table string, data interface{}, batchSize int) (rowsAffected int64, err error) {
	tx := c.db
	if table != "" {
		tx = tx.Table(table)
	}
	result := tx.CreateInBatches(data, batchSize)
	return result.RowsAffected, result.Error
}

// FindById rowsAffected:返回插入记录的条数,err:返回error
func (c Curd) FindById(table string, id int, resp interface{}) error {
	return c.db.Table(table).Where("id = ?", id).Scan(resp).Error
}

// FindList 查询列表
func (c Curd) FindList(table string, resp interface{}, params QueryParams) error {
	db := c.db.Table(table)
	for _, v := range params.Query {
		db = db.Where(v)
	}
	if params.Offset != 0 {
		db = db.Offset(params.Offset)
	}
	if params.Limit != 0 {
		db = db.Limit(params.Limit)
	}
	if params.Sort != "" {
		db = db.Order(params.Sort)
	}
	if params.Group != "" {
		db = db.Group(params.Group)
	}
	return db.Scan(resp).Error
}

// FindFirst 查询第一条数据
func (c Curd) FindFirst(table string, query []string, resp interface{}) error {
	tx := c.db.Table(table)
	for _, v := range query {
		tx = tx.Where(v)
	}
	return tx.First(resp).Error
}

// FindLast 查询最后一条数据
func (c Curd) FindLast(table string, query []string, resp interface{}) error {
	tx := c.db.Table(table)
	for _, v := range query {
		tx = tx.Where(v)
	}
	return tx.Last(resp).Error
}

// Update 更新数据
func (c Curd) Update(table string, params UpdateParams) (rowsAffected int64, err error) {
	tx := c.db.Table(table)
	for _, v := range params.Query {
		tx = tx.Where(v)
	}
	for k, v := range params.Update {
		tx = tx.Update(k, v)
	}
	return tx.RowsAffected, tx.Error
}

// Delete 删除数据
func (c Curd) Delete(table string, params DeleteParams) (rowsAffected int64, err error) {
	tx := c.db.Table(table)
	if len(params.Query) == 0 {
		err = errors.New("sensitive operation. It is forbidden to delete all table data")
		return
	}
	for _, v := range params.Query {
		tx = tx.Where(v)
	}
	tx.Delete(params.Query)
	return tx.RowsAffected, tx.Error
}
