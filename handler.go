package v8gohttp

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"rogchap.com/v8go"
)

const handlerEvalTimeout = 100 * time.Millisecond

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
		defer vm.Dispose()

		ctx, err := v8go.NewContext(vm)
		if err != nil {
			panic(err)
		}
		defer ctx.Close()

		if _, err := ctx.RunScript(libJs, "lib.js"); err != nil {
			panic(err)
		}

		if err := evalHandler(req.Context(), ctx, handler); err != nil {
			panic(err)
		}

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

		<-reqCtx.Done()

		if reqCtx.Err() != context.Canceled {
			panic(reqCtx.Err())
		}
	})
}

func Handle(pattern, handler string) {
	http.Handle(pattern, Handler(handler))
}

func evalHandler(parent context.Context, v8ctx *v8go.Context, handler string) error {
	ctx, cancel := context.WithTimeout(parent, handlerEvalTimeout)
	var err error

	go func() {
		defer cancel()
		_, err = v8ctx.RunScript(handler, "handler.js")
	}()

	<-ctx.Done()

	if err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("Handler evaluation took too long %w", ctx.Err())
	}

	handlerRef, err := v8ctx.Global().Get("handler")
	if err != nil {
		return fmt.Errorf("Could not get handler reference %w", ctx.Err())
	}

	if !handlerRef.IsFunction() {
		return fmt.Errorf("Handler reference is not a function")
	}

	return nil
}
