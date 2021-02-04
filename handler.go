package v8gohttp

import (
	_ "embed"
	"io"
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

	//go:embed call-handler.js
	callHandlerJs string
)

type bodyReader http.Request

func (r *bodyReader) callback(info *v8go.FunctionCallbackInfo) *v8go.Value {
	b, err := io.ReadAll((*http.Request)(r).Body)
	if err != nil {
		panic(err)
	}

	vm, err := info.Context().Isolate()
	if err != nil {
		panic(err)
	}

	v, err := v8go.NewValue(vm, string(b))
	if err != nil {
		panic(err)
	}

	return v
}

func Handler(handler string) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		vm, err := v8go.NewIsolate()
		if err != nil {
			panic(err)
		}
		// FIXME too early...
		// defer vm.Dispose()
		req.Context()

		headers, err := v8go.NewObjectTemplate(vm)
		if err != nil {
			panic(err)
		}
		for k := range req.Header {
			if err := headers.Set(k, req.Header.Get(k)); err != nil {
				panic(err)
			}
		}

		bodyReader, err := v8go.NewFunctionTemplate(vm, (*bodyReader)(req).callback)
		if err != nil {
			panic(err)
		}

		callback, err := v8go.NewFunctionTemplate(vm, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			// FIXME headers
			body, status := info.Args()[0].String(), info.Args()[1].Int32()
			if status != 0 {
				res.WriteHeader(int(status))
			}
			if _, err := res.Write(([]byte)(body)); err != nil {
				panic(err)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}

		infoTmpl, err := v8go.NewObjectTemplate(vm)
		if err != nil {
			panic(err)
		}
		for k, v := range map[string]interface{}{
			"url":        req.RequestURI,
			"method":     req.Method,
			"headers":    headers,
			"bodyReader": bodyReader,
			"callback":   callback,
		} {
			if err := infoTmpl.Set(k, v); err != nil {
				panic(err)
			}
		}

		ctx, err := v8go.NewContext(vm)
		if err != nil {
			panic(err)
		}
		defer ctx.Close()

		if _, err := ctx.RunScript(responseJs, "response.js"); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(handler, "handler.js"); err != nil {
			panic(err)
		}

		// FIXME check handler

		if _, err := ctx.RunScript(requestJs, "request.js"); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(fetchEventJs, "fetch-event.js"); err != nil {
			panic(err)
		}

		info, err := infoTmpl.NewInstance(ctx)
		if err != nil {
			panic(err)
		}
		if err := ctx.Global().Set("info", info); err != nil {
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
