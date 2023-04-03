package demo

// type UserService struct {
// 	db *sql.DB
// }
//
// // 调用者来注入 db
// // 所以调用者可以注入 mock 的 db，方便测试
// func NewUserService(db *sql.DB) *UserService {
// 	return &UserService{
// 		db: db,
// 	}
// }
//
// // 调用者完全没办法传入一个 sqlmock
// func NewUserService(driver string, dsn string) *UserService {
// 	db, err := sql.Open(driver, dsn)
// 	return &UserService{
// 		db: db,
// 	}
// }
//
// const (
// 	driver = "mysql"
// 	dsn    = ""
// )
//
// func NewUserService() *UserService {
// 	db, err := sql.Open(driver, dsn)
// 	return &UserService{
// 		db: db,
// 	}
// }
//
// func main() {
// 	db := sql.Open()
//
// 	userService := NewUserService(db)
// 	uc := UserController{
// 		us: userService,
// 	}
// 	g := gin.Default()
// 	g.POST("/logi", uc.Login)
// 	g.Run(":8081")
// }
//
// type UserController struct {
// 	us *UserService
// }
//
// func (*UserController) Login(ctx *gin.Context) {
//
// }
