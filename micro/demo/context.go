package demo

import "context"


func CtxWithOneway(ctx context.Context) context.Context {
	return context.WithValue(ctx, "oneway", true)
}

func isOneway(ctx context.Context) bool {
	return ctx.Value("oneway") != nil
}
