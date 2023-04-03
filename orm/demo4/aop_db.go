package orm

// type AopDBContext struct {
// 	query string
// 	args []any
// }
//
// type AopDBResult struct {
// 	row *sql.Row
// }
//
// type Handler func(ctx *AopDBContext) *AopDBResult
// type Middleware func(next Handler) Handler
// // sql.DB 上加 AOP 解决方案，我该怎么加
//
// type AopDB struct {
// 	db *sql.DB
// 	ms []Middleware
// }
//
// func (db *AopDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
// 	// 在这之前做点东西
// 	 var handler Handler = func(actx *AopDBContext) *AopDBResult {
// 		 row := db.db.QueryRowContext(ctx, actx.query, actx.args...)
// 		 return &AopDBResult{
// 			 row: row,
// 		 }
// 	 }
//
// 	 // 在这之后也做点东西
// 	for i := len(db.ms) - 1; i >= 0; i -- {
// 		handler = db.ms[i](handler)
// 	}
// 	res := handler(&AopDBContext{
//
// 	})
// 	return res.row
// }
