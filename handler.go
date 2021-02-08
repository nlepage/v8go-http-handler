package v8gohttp

import (
	_ "embed"
	"net/http"

	"rogchap.com/v8go"
)

var (
	//go:embed request.js
	requestJs string

	//go:embed response.js
	responseJs string

	//go:embed fetch-event.js
	fetchEventJs string

	libJs = requestJs + responseJs + fetchEventJs

	//go:embed call-handler.js
	callHandlerJs string
)

func Handler(handler string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		vm, err := v8go.NewIsolate()
		if err != nil {
			panic(err)
		}
		// FIXME too early...
		// defer vm.Dispose()

		ctx, err := v8go.NewContext(vm)
		if err != nil {
			panic(err)
		}
		defer ctx.Close()

		if _, err := ctx.RunScript(libJs, "lib.js"); err != nil {
			panic(err)
		}

		// FIXME outside?
		if _, err := ctx.RunScript(handler, "handler.js"); err != nil {
			panic(err)
		}
		// FIXME very short timeout?
		// FIXME check handler

		reqCtx := newRequestContext(vm, res, req)

		reqCtxObj, err := reqCtx.instance(ctx)
		if err != nil {
			panic(err)
		}
		if err := ctx.Global().Set("reqCtx", reqCtxObj); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(callHandlerJs, "call-handler.js"); err != nil {
			panic(err)
		}
	})
}

func Handle(pattern, handler string) {
	http.Handle(pattern, Handler(handler))
}
