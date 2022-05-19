package test

import (
	"errors"
	"fmt"
	"github.com/yueja/cf-gorm/mysql/db"
	"testing"
)

type User struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func Test_Save(t *testing.T) {
	user := []User{{Name: "张三", Age: 12, Phone: "11111"}, {Name: "李四", Age: 21, Phone: "22222"}}
	var (
		rowsAffected int64
		err          error
	)
	rowsAffected, err = db.DbCurd().Save("t_user", user)
	fmt.Println("数据：", user, rowsAffected, err)
}

func Test_Save2(t *testing.T) {
	user := []User{{Name: "张三", Age: 12, Phone: "11111"}, {Name: "李四", Age: 21, Phone: "22222"}}
	var (
		rowsAffected int64
		err          error
	)
	rowsAffected, err = db.SpecifyDbCurd("default").Save("t_user", user)
	fmt.Println("数据：", user, rowsAffected, err)
}

func Test_Create(t *testing.T) {
	user := make([]map[string]interface{}, 0)
	user = append(user, map[string]interface{}{
		"name":  "张三001",
		"age":   12,
		"phone": "1111",
	})
	user = append(user, map[string]interface{}{
		"name":  "张三002",
		"age":   15,
		"phone": "2222",
	})
	var (
		rowsAffected int64
		err          error
	)
	rowsAffected, err = db.DbCurd().Create("t_user", user)
	fmt.Println("数据：", user, rowsAffected, err)
}

func Test_CreateInBatches(t *testing.T) {
	user := []User{{Name: "张三", Age: 12, Phone: "11111"}, {Name: "李四", Age: 21, Phone: "22222"}}
	var (
		rowsAffected int64
		err          error
	)
	rowsAffected, err = db.SpecifyDbCurd("default").CreateInBatches("t_user", user, 2)
	fmt.Println("数据：", user, rowsAffected, err)
}

func Test_FindById(t *testing.T) {
	var (
		user User
		err  error
	)
	err = db.DbCurd().FindById("t_user", 6, &user)
	fmt.Println("数据：", user, err)
}

func Test_FindList(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id > 1", "age > 19")
	var (
		userList []User
		params   = db.QueryParams{
			SelectFiled: []string{"name", "age"},
			Query:       data,
			Offset:      0,
			Limit:       2,
			Sort:        "age desc",
			//Group:  "phone",
		}
	)
	err := db.DbCurd().FindList("t_user", &userList, params)
	fmt.Println("数据：", userList, err)
}

func Test_FindCount(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id > 1", "age > 2")
	var (
		params = db.QueryParams{
			Query: data,
			Sort:  "age desc",
			//Group:  "phone",
		}
		count int64
		err   error
	)
	count, err = db.DbCurd().FindCount("t_user", params)
	fmt.Println("数据：", count, err)
}

func Test_FindFirst(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id >= 4")
	var user User
	err := db.DbCurd().FindFirst("t_user", db.QueryParams{Query: data}, &user)
	fmt.Println("数据：", user, err)
}

func Test_FindLast(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id <10")
	var user User
	err := db.DbCurd().FindLast("t_user", db.QueryParams{Query: data}, &user)
	fmt.Println("数据：", user, err)
}

func Test_Update(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id > 1", "age > 19")
	var params1 = db.UpdateParams{
		Query: data,
		Update: map[string]interface{}{
			"age": "99",
		},
	}
	var (
		rowsAffected int64
		err          error
	)
	rowsAffected, err = db.DbCurd().Update("t_user", params1)
	fmt.Println("数据：", rowsAffected, err)
}

func Test_Delete(t *testing.T) {
	data := make([]string, 0)
	data = append(data, "id > 1", "age > 19")
	var (
		rowsAffected int64
		err          error
		user         User
		params       = db.DeleteParams{
			Query: data,
			Model: &user,
		}
	)
	rowsAffected, err = db.DbCurd().Delete("t_user", params)
	fmt.Println("数据", rowsAffected, err)
}

func Test_Tx(t *testing.T) {
	var (
		data  = make([]string, 0)
		data1 = make([]string, 0)
	)
	data = append(data, "id = 10")
	data1 = append(data1, "id = 12")
	var (
		rowsAffected int64
		user         User
		params       = db.DeleteParams{
			Query: data,
			Model: &user,
		}
		params1 = db.DeleteParams{
			Query: data1,
			Model: &user,
		}
	)
	err := db.WithTransaction(func(tx *db.Curd) (err error) {
		rowsAffected, err = tx.Delete("t_user", params)
		err = errors.New("测试错误")
		if err != nil {
			return err
		}
		fmt.Println("数据", rowsAffected, err)

		rowsAffected, err = tx.Delete("t_user", params1)
		fmt.Println("数据1", rowsAffected, err)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func Test_UpsertFiledDefault(t *testing.T) {
	var (
		rowsAffected int64
		err          error
		user         = User{
			Id:    18,
			Name:  "哈哈哈哈001",
			Age:   21,
			Phone: "1320",
		}
		params = db.UpsertFiledDefaultParams{
			ConflictFiled: []string{"name", "age"},
			Update: map[string]interface{}{
				"name": "我是更新名字",
				"age":  22,
			},
			Model: &user,
		}
	)
	rowsAffected, err = db.DbCurd().UpsertFiledDefault("t_user", params)
	fmt.Println("数据", rowsAffected, user, err)
}

func Test_UpsertFiled(t *testing.T) {
	var (
		rowsAffected int64
		err          error
		user         = User{
			Id:    12,
			Name:  "我是更新名字",
			Age:   21,
			Phone: "666",
		}
		params = db.UpsertParams{
			ConflictFiled: []string{"name", "age"},
			UpdateFiled:   []string{"name", "age", "phone"},
			Model:         &user,
		}
	)
	rowsAffected, err = db.DbCurd().UpsertFiled("t_user", params)
	fmt.Println("数据", rowsAffected, user, err)
}

func Test_Upsert(t *testing.T) {
	var (
		rowsAffected int64
		err          error
		user         = User{
			Name:    "我是更新名字",
			Age:     21,
			Phone:   "99",
			Address: "36",
		}
		params = db.UpsertParams{
			ConflictFiled: []string{"name", "age"},
			Model:         &user,
		}
	)
	rowsAffected, err = db.DbCurd().Upsert("t_user", params)
	fmt.Println("数据", rowsAffected, user, err)
}
