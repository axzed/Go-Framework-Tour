package main

import (
	"context"
	"gitee.com/geektime-geekbang/geektime-go/orm"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/opentelemetry"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/prometheus"
	"gitee.com/geektime-geekbang/geektime-go/orm/middleware/querylog"
	"go.uber.org/zap"
)

func initDB() *orm.DB {
	db, err := orm.Open("mysql", "root:root@tcp(localhost:13306)/userapp",
		orm.DBWithMiddleware(
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		defer func() {
			// 			if m:=recover(); m !=nil {
			//
			// 			}
			// 		}()
			// 		next(ctx, qc)
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	// 不带 WHERE 的语句不允许执行
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		if qc.Type == "INSERT" {
			// 			next(ctx, qc)
			// 			return
			// 		}
			// 		q, err := qc.Query()
			// 		// 更进一步，你可以在这里检测 WHERE 带不带索引列
			// 		if strings.Index(q.SQL, "WHERE") == -1 {
			// 			panic("查询语句不带 WHERE")
			// 		}
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	// 不允许删除
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		if qc.Type == "DELETE" {
			// 			panic("不能删除")
			// 		}
			// 		next(ctx, qc)
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		before := time.Now()
			// 		defer func() {
			// 			duration := before.Util()
			// 			if duration > time.Millisecond * 50 {
			// 				// 记录一下你的慢查询，或者告警
			// 			}
			// 		}()
			// 		next(ctx, qc)
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		// 如果是更新了 DB，UPDATE, DELETE 或者 INSERT
			// 		// 顺便更新缓存
			// 		// 这里就是 write-through
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		// 要针对不同的 model 来进行分发
			// 		// 你可以在这里，第一步查找缓存
			// 		// 第二步查找 DB 来集成缓存功能
			// 		// 把 read through 挪过来这里
			// 	}
			// },
			// func(next orm.HandleFunc) orm.HandleFunc {
			// 	return func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
			// 		if qc.Type == "SELECT" && qc.Model.TableName == "user"{
			// 			// 返回 mock 数据
			// 		}
			// 	}
			// },
			prometheus.MiddlewareBuilder{
			Name:        "userapp",
			Subsystem:   "orm",
			ConstLabels: map[string]string{"db": "userapp"}}.Build(),
			opentelemetry.MiddlewareBuilder{}.Build(),
			querylog.NewBuilder().LogFunc(func(sql string, args []any) {
				// 一般不建议记录参数，因为参数里面可能有一些加密信息，
				// 当然如果你确定自己是在开发环境下使用，那么你就可以记录参数
				zap.L().Info("query", zap.String("SQL", sql))
			}).Build()))
	if err != nil {
		panic(err)
	}
	return db
}
