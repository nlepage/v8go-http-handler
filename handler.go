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

func Handler(js string) http.Handler {
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

		global, err := v8go.NewObjectTemplate(vm)
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
			if err := global.Set(k, v); err != nil {
				panic(err)
			}
		}

		ctx, err := v8go.NewContext(vm, global)
		if err != nil {
			panic(err)
		}
		defer ctx.Close()

		if _, err := ctx.RunScript(requestJs, "request.js"); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(fetchEventJs, "fetch-event.js"); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(responseJs, "response.js"); err != nil {
			panic(err)
		}

		if _, err := ctx.RunScript(`
			const e = new FetchEvent(new Request(url, { method, headers, bodyReader }), callback)

			async function handle(e) {
				const { name } = await e.request.json()
				e.respondWith(new Response('Hello ' + name + '!'))
			}

			handle(e)
		`, "test.js"); err != nil {
			panic(err)
		}
	})
}
