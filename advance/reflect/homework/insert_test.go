package homework

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInsertStmt(t *testing.T) {
	testCases := []struct {
		name     string
		entity   interface{}
		wantErr  error
		wantSQL  string
		wantArgs []interface{}
	}{
		{
			name:    "nil",
			wantErr: errInvalidEntity,
		},
		{
			name:    "empty struct",
			entity:  Empty{},
			wantErr: errInvalidEntity,
		},
		{
			name:   "simple struct",
			entity: BaseEntity{},
			// 稍微注意一下，这里我们用反射取出来的不是存粹的 nil，而是一个带了类型的 nil
			wantArgs: []interface{}{int64(0), (*int64)(nil)},
			wantSQL:  "INSERT INTO `BaseEntity`(`CreateTime`,`UpdateTime`) VALUES(?,?);",
		},
		{
			name:     "pointer",
			entity:   &BaseEntity{CreateTime: 123, UpdateTime: ptrInt64(234)},
			wantArgs: []interface{}{int64(123), ptrInt64(234)},
			wantSQL:  "INSERT INTO `BaseEntity`(`CreateTime`,`UpdateTime`) VALUES(?,?);",
		},
		{
			name: "multiple pointer",
			entity: func() interface{} {
				entity := &BaseEntity{CreateTime: 123, UpdateTime: ptrInt64(234)}
				return &entity
			}(),
			wantErr: errInvalidEntity,
		},
		{
			// 组合
			name: "composition",
			entity: User{
				BaseEntity: BaseEntity{
					CreateTime: 123,
					UpdateTime: ptrInt64(456),
				},
				Id:       789,
				NickName: sql.NullString{String: "Tom", Valid: true},
			},
			wantArgs: []interface{}{int64(123), ptrInt64(456), uint64(789),
				sql.NullString{String: "Tom", Valid: true}, (*sql.NullInt32)(nil)},
			wantSQL: "INSERT INTO `User`(`CreateTime`,`UpdateTime`,`Id`,`NickName`,`Age`) VALUES(?,?,?,?,?);",
		},
		{
			name: "deep composition",
			entity: &Buyer{
				User: User{
					BaseEntity: BaseEntity{
						CreateTime: 123,
						UpdateTime: ptrInt64(456),
					},
					Id:       789,
					NickName: sql.NullString{String: "Tom", Valid: true},
					Age:      &sql.NullInt32{Int32: 18, Valid: true},
				},
				Address: "China",
			},
			wantArgs: []interface{}{int64(123), ptrInt64(456), uint64(789),
				sql.NullString{String: "Tom", Valid: true}, &sql.NullInt32{Int32: 18, Valid: true}, "China"},
			wantSQL: "INSERT INTO `Buyer`(`CreateTime`,`UpdateTime`,`Id`,`NickName`,`Age`,`Address`) VALUES(?,?,?,?,?,?);",
		},
		{
			name: "multiple composition",
			entity: &Customer{
				Buyer: Buyer{
					User: User{
						BaseEntity: BaseEntity{
							CreateTime: 123,
							UpdateTime: ptrInt64(456),
						},
						Id:       789,
						NickName: sql.NullString{String: "Tom", Valid: true},
						Age:      &sql.NullInt32{Int32: 18, Valid: true},
					},
					Address: "China",
				},
				BaseEntity: BaseEntity{
					CreateTime: 987,
					UpdateTime: ptrInt64(654),
				},
				Company: "DM",
			},
			wantArgs: []interface{}{int64(123), ptrInt64(456), uint64(789),
				sql.NullString{String: "Tom", Valid: true}, &sql.NullInt32{Int32: 18, Valid: true}, "China", "DM"},
			wantSQL: "INSERT INTO `Customer`(`CreateTime`,`UpdateTime`,`Id`,`NickName`,`Age`,`Address`,`Company`) VALUES(?,?,?,?,?,?,?);",
		},
		{
			// 使用指针的组合，我们不会深入解析，会出现很奇怪的结果
			name: "pointer composition",
			entity: InvalidUser{
				BaseEntity: &BaseEntity{},
				Address:    "China",
			},
			// &BaseEntity{} 这个参数发送到 driver 那里，会出现无法解析的情况
			wantArgs: []interface{}{&BaseEntity{}, "China"},
			wantSQL:  "INSERT INTO `InvalidUser`(`BaseEntity`,`Address`) VALUES(?,?);",
		},
		{
			name:   "not embed field",
			entity: Seller{User: User{}},
			// 顺便测试一下单个字段
			wantArgs: []interface{}{User{}},
			wantSQL:  "INSERT INTO `Seller`(`User`) VALUES(?);",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, args, err := InsertStmt(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if tc.wantErr != nil {
				// 预期会有错误返回，就不需要进一步校验其它两个返回值了
				return
			}
			assert.Equal(t, tc.wantSQL, query)
			assert.Equal(t, tc.wantArgs, args)
		})
	}
}

type Empty struct {
}

type BaseEntity struct {
	CreateTime int64
	UpdateTime *int64
}

type User struct {
	BaseEntity
	Id       uint64
	NickName sql.NullString
	Age      *sql.NullInt32
}

type Buyer struct {
	User
	Address string
}

type InvalidUser struct {
	*BaseEntity
	Address string
}

type Customer struct {
	Buyer
	BaseEntity
	Company string
}

// Seller 注意和 Buyer 的区别。在 Seller 这种模式下，我们会认为，它们是一个关联关系，比如说一对一关系
// 而在 Buyer 的模式下，我们会认为它们是同一张表。
type Seller struct {
	User User
}

func ptrInt64(val int64) *int64 {
	return &val
}
