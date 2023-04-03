package main

import (
	"errors"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen/annotation"
	"gitee.com/geektime-geekbang/geektime-go/advance/template/gen/http"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"os"
	"path/filepath"
	"testing"
)

func Test_parseServiceDefinition(t *testing.T) {
	pkg := "main"
	tests := []struct {
		name    string
		typ     annotation.Type
		want    http.ServiceDefinition
		wantErr error
	}{
		{
			/*
				测试这种定义
				type UserService interface {
				    Get(ctx context.Context, req *GetUserReq)(*GetUserResp, error)
				}
			*/
			name: "no annotations",
			typ: annotation.Type{
				Annotations: annotation.Annotations[*ast.TypeSpec]{
					Node: &ast.TypeSpec{
						Name: ident("UserService"),
					},
				},
				Fields: []annotation.Field{
					{
						Annotations: annotation.Annotations[*ast.Field]{
							Node: &ast.Field{
								Names: idents("Get"),
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{},
											{Type: &ast.StarExpr{X: ident("GetUserReq")}},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{Type: &ast.StarExpr{X: ident("GetUserResp")}},
											{},
										},
									},
								},
							},
						},
					},
				},
			},
			want: http.ServiceDefinition{
				Package: pkg,
				Name:    "UserService",
				Methods: []http.ServiceMethod{
					{
						Name:         "Get",
						Path:         "/Get",
						ReqTypeName:  "GetUserReq",
						RespTypeName: "GetUserResp",
					},
				},
			},
		},
		{
			/*
				测试这种定义，缺了 context 参数
				type UserService interface {
				    Get(req *GetUserReq)(*GetUserResp, error)
				}
			*/
			name: "invalid parameter",
			typ: annotation.Type{
				Annotations: annotation.Annotations[*ast.TypeSpec]{
					Node: &ast.TypeSpec{
						Name: ident("UserService"),
					},
				},
				Fields: []annotation.Field{
					{
						Annotations: annotation.Annotations[*ast.Field]{
							Node: &ast.Field{
								Names: idents("Get"),
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{Type: &ast.StarExpr{X: ident("GetUserReq")}},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: errors.New("gen: 方法必须接收两个参数，其中第一个参数是 context.Context，第二个参数请求"),
		},
		{
			/*
				测试这种定义
				type UserService interface {
				    Get(ctx context.Context, req *GetUserReq) *GetUserResp
				}
			*/
			name: "invalid result",
			typ: annotation.Type{
				Annotations: annotation.Annotations[*ast.TypeSpec]{
					Node: &ast.TypeSpec{
						Name: ident("UserService"),
					},
				},
				Fields: []annotation.Field{
					{
						Annotations: annotation.Annotations[*ast.Field]{
							Node: &ast.Field{
								Names: idents("Get"),
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{},
											{Type: &ast.StarExpr{X: ident("GetUserReq")}},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{Type: &ast.StarExpr{X: ident("GetUserResp")}},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: errors.New("gen: 方法必须返回两个参数，其中第一个返回值是响应，第二个返回值是error"),
		},
		// 你可以进一步测试非指针的情况

		// 使用了注解
		{
			/*
				测试这种定义
				type UserService interface {
				    Get(ctx context.Context, req *GetUserReq)(*GetUserResp, error)
				}
			*/
			name: "with annotations",
			typ: annotation.Type{
				Annotations: annotation.Annotations[*ast.TypeSpec]{
					Node: &ast.TypeSpec{
						Name: ident("UserService"),
					},
					Ans: []annotation.Annotation{{"ServiceName", "MyUserService"}},
				},
				Fields: []annotation.Field{
					{
						Annotations: annotation.Annotations[*ast.Field]{
							Node: &ast.Field{
								Names: idents("Get"),
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{},
											{Type: &ast.StarExpr{X: ident("GetUserReq")}},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{Type: &ast.StarExpr{X: ident("GetUserResp")}},
											{},
										},
									},
								},
							},
							Ans: []annotation.Annotation{{Key: "Path", Value: "/user"}},
						},
					},
				},
			},
			want: http.ServiceDefinition{
				Package: pkg,
				Name:    "MyUserService",
				Methods: []http.ServiceMethod{
					{
						Name:         "Get",
						Path:         "/user",
						ReqTypeName:  "GetUserReq",
						RespTypeName: "GetUserResp",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseServiceDefinition(pkg, tt.typ)
			assert.Equal(t, tt.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, got, tt.want)
		})
	}
}

func ident(name string) *ast.Ident {
	return &ast.Ident{Name: name}
}

func idents(names ...string) []*ast.Ident {
	res := make([]*ast.Ident, 0, len(names))
	for _, n := range names {
		res = append(res, &ast.Ident{Name: n})
	}
	return res
}

// 注意，如果你在 http 里面写的模板跟我不同，那么你需要修改相应的 .txt 文件
// 即期望结果
func Test_gen(t *testing.T) {
	// 我们使用 testdata 目录来测试
	err := gen("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name     string
		genFile  string
		wantFile string
	}{
		{
			name:     "order",
			genFile:  "./testdata/my_order_service_gen.go",
			wantFile: "./testdata/my_order_service_gen.txt",
		},
		{
			name:     "user",
			genFile:  "./testdata/user_service_gen.go",
			wantFile: "./testdata/user_service_gen.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genOrder, err := os.ReadFile(tc.genFile)
			if err != nil {
				t.Fatal(err)
			}
			wantOrder, err := os.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, wantOrder, genOrder)
		})
	}
}

func Test_scanFiles(t *testing.T) {
	// 这里我们简单测试一下 testdata 就可以
	files, err := scanFiles("./testdata")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatal(err)
	}
	// 我的实现里面是用绝对地址的，所以我的测试用例也是绝对地址，如果你用相对地址，那么这个测试用例你需要修改
	user, err := filepath.Abs("./testdata/user_service.go")
	if err != nil {
		t.Fatal(err)
	}
	order, err := filepath.Abs("./testdata/order_service.go")
	if err != nil {
		t.Fatal(err)
	}
	// 在我机器上读取到的顺序是 order 在前，user 在后。按照道理来说，在你的机器上也会一样
	assert.Equal(t, []string{order, user}, files)
}

func Test_parseFiles(t *testing.T) {
	// 同样用 testdata 来测试
	// 在实际工作中，这种读写文件的测试，在逼不得已的时候，都放过去 testdata
	// 因为 parseServiceDefinition 的测试确保了我们构造 service definition 是对的
	// 而且 annotation 包的测试确保了我们的 AST 操作是对的
	// 所以我们只需要简单测试一下
	// 如果寻求严谨，那么也可以深入测试

	files, err := scanFiles("./testdata")
	if err != nil {
		t.Fatal(err)
	}

	defs, err := parseFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, []http.ServiceDefinition{
		{
			Package: "testdata",
			Name:    "MyOrderService",
			Methods: []http.ServiceMethod{
				{
					Name:         "Create",
					Path:         "/Create",
					ReqTypeName:  "CreateOrderReq",
					RespTypeName: "CreateOrderResp",
				},
			},
		},
		{
			Package: "testdata",
			Name:    "UserService",
			Methods: []http.ServiceMethod{
				{
					Name:         "Get",
					Path:         "/Get",
					ReqTypeName:  "GetUserReq",
					RespTypeName: "GetUserResp",
				},
				{
					Name:         "Update",
					Path:         "/user/update",
					ReqTypeName:  "UpdateUserReq",
					RespTypeName: "UpdateUserResp",
				},
			},
		},
	}, defs)
}
