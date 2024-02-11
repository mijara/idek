package idek

// OnFinish funcs will get executed as the last step before rendering to client
type OnFinishFunc func(ctx *Context, statusCode int, output any)
