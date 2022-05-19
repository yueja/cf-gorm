package db

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
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

func (c Curd) setQueryParams(db *gorm.DB, params QueryParams) {
	c.setSelectFiled(db, params.SelectFiled)
	for _, v := range params.Query {
		db.Where(v)
	}
	for _, v := range params.Not {
		db.Not(v)
	}
	for _, v := range params.Or {
		db.Or(v)
	}
	if params.Offset != 0 {
		db.Offset(params.Offset)
	}
	if params.Limit != 0 {
		db.Limit(params.Limit)
	}
	if params.Sort != "" {
		db.Order(params.Sort)
	}
	if params.Group != "" {
		db.Group(params.Group)
	}
}

func (c Curd) setSelectFiled(db *gorm.DB, SelectFiled []string) {
	if len(SelectFiled) == 0 {
		db.Select("*")
	}
	if len(SelectFiled) != 0 {
		var selectFiled string
		for _, v := range SelectFiled {
			if v != "" {
				continue
			}
			selectFiled += v + ","
		}
		selectFiled = strings.TrimRight(selectFiled, ",")
		if selectFiled != "" {
			db.Select(selectFiled)
		}
	}
}

// Save 插入数据，rowsAffected:返回插入记录的条数
func (c Curd) Save(table string, data interface{}) (rowsAffected int64, err error) {
	db := c.db
	if table != "" {
		db = db.Table(table)
	}
	db.Save(data)
	return db.RowsAffected, db.Error
}

// Create 插入数据，rowsAffected:返回插入记录的条数
func (c Curd) Create(table string, data interface{}) (rowsAffected int64, err error) {
	db := c.db
	if table != "" {
		db.Table(table)
	}
	db.Create(data)
	return db.RowsAffected, db.Error
}

// CreateInBatches 批量插入数据 rowsAffected:返回插入记录的条数
func (c Curd) CreateInBatches(table string, data interface{}, batchSize int) (rowsAffected int64, err error) {
	db := c.db
	if table != "" {
		db.Table(table)
	}
	db.CreateInBatches(data, batchSize)
	return db.RowsAffected, db.Error
}

// FindById rowsAffected:返回插入记录的条数,err:返回error
func (c Curd) FindById(table string, id int, resp interface{}) error {
	return c.db.Table(table).Where("id = ?", id).Scan(resp).Error
}

// FindList 查询列表
func (c Curd) FindList(table string, resp interface{}, params QueryParams) error {
	db := c.db.Table(table)
	c.setQueryParams(db, params)
	return db.Scan(resp).Error
}

// FindCount 查询总数
func (c Curd) FindCount(table string, params QueryParams) (count int64, err error) {
	db := c.db.Table(table)
	c.setQueryParams(db, params)
	err = db.Count(&count).Error
	return
}

// FindFirst 查询第一条数据
func (c Curd) FindFirst(table string, params QueryParams, resp interface{}) error {
	db := c.db.Table(table)
	c.setQueryParams(db, params)
	return db.First(resp).Error
}

// FindLast 查询最后一条数据
func (c Curd) FindLast(table string, params QueryParams, resp interface{}) error {
	db := c.db.Table(table)
	c.setQueryParams(db, params)
	return db.Last(resp).Error
}

// Update 更新数据
func (c Curd) Update(table string, params UpdateParams) (rowsAffected int64, err error) {
	db := c.db.Table(table)
	for _, v := range params.Query {
		db.Where(v)
	}
	for k, v := range params.Update {
		db.Update(k, v)
	}
	return db.RowsAffected, db.Error
}

// Delete 删除数据
func (c Curd) Delete(table string, params DeleteParams) (rowsAffected int64, err error) {
	db := c.db.Table(table)
	if len(params.Query) == 0 {
		err = errors.New("sensitive operation. It is forbidden to delete all table data")
		return
	}
	for _, v := range params.Query {
		db.Where(v)
	}
	db.Delete(params.Query)
	return db.RowsAffected, db.Error
}

// UpsertFiledDefault 更新或新增（指定字段更新为默认值）
// Columns必须是唯一索引，主键索引或者其他唯一联合索引都行
// 当Column不是唯一索引，参数有主键，则默认以主键Upsert，否则新增
func (c Curd) UpsertFiledDefault(table string, params UpsertFiledDefaultParams) (rowsAffected int64, err error) {
	db := c.db.Table(table)
	var column []clause.Column
	for _, v := range params.ConflictFiled {
		column = append(column, clause.Column{Name: v})
	}
	db.Clauses(clause.OnConflict{
		Columns:   column,
		DoUpdates: clause.Assignments(params.Update), // 冲突时, 更新指定字段为默认值
	}).Create(params.Model)
	return db.RowsAffected, db.Error
}

// UpsertFiled 更新或新增（指定字段）
// Columns必须是唯一索引，主键索引或者其他唯一联合索引都行
// 当Column不是唯一索引，参数有主键，则默认以主键Upsert，否则新增
func (c Curd) UpsertFiled(table string, params UpsertParams) (rowsAffected int64, err error) {
	db := c.db.Table(table)
	var column []clause.Column
	for _, v := range params.ConflictFiled {
		column = append(column, clause.Column{Name: v})
	}
	db.Clauses(clause.OnConflict{
		Columns:   column,
		DoUpdates: clause.AssignmentColumns(params.UpdateFiled), // 冲突时, 更新指定字段
	}).Create(params.Model)
	return db.RowsAffected, db.Error
}

// Upsert 更新或新增（全部字段）
// Columns必须是唯一索引，主键索引或者其他唯一联合索引都行
// 当Column不是唯一索引，参数有主键，则默认以主键Upsert，否则新增
func (c Curd) Upsert(table string, params UpsertParams) (rowsAffected int64, err error) {
	db := c.db.Table(table)
	var column []clause.Column
	for _, v := range params.ConflictFiled {
		column = append(column, clause.Column{Name: v})
	}
	db.Clauses(clause.OnConflict{
		Columns:   column,
		UpdateAll: true, // 冲突时, 更新所有字段
	}).Create(params.Model)
	return db.RowsAffected, db.Error
}

// Distinct 从模型中选择不相同的值
func (c Curd) Distinct(table string, resp interface{}, params DistinctParams) error {
	db := c.db.Table(table).Distinct(params.DistinctFiled)
	c.setQueryParams(db, params.QueryParams)
	return db.Find(resp).Error
}

// Joins 联合查询
func (c Curd) Joins(table string, resp interface{}, params JoinsParams) error {
	db := c.db.Table(table)
	c.setQueryParams(db, params.QueryParams)
	for _, v := range params.Joins {
		db.Joins(v)
	}
	return db.Scan(resp).Error
}
