package v8gohttp

import (
	"io"
	"net/http"

	"rogchap.com/v8go"
)

type requestContext struct {
	vm  *v8go.Isolate
	res http.ResponseWriter
	req *http.Request
}

func newRequestContext(vm *v8go.Isolate, res http.ResponseWriter, req *http.Request) *requestContext {
	return &requestContext{vm, res, req}
}

func (reqCtx *requestContext) instance(ctx *v8go.Context) (*v8go.Object, error) {
	tmpl, err := v8go.NewObjectTemplate(reqCtx.vm)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Set("url", reqCtx.req.RequestURI); err != nil {
		return nil, err
	}

	if err := tmpl.Set("method", reqCtx.req.Method); err != nil {
		return nil, err
	}

	headers, err := reqCtx.headers()
	if err != nil {
		return nil, err
	}
	if err := tmpl.Set("headers", headers); err != nil {
		return nil, err
	}

	readBody, err := v8go.NewFunctionTemplate(reqCtx.vm, reqCtx.readBody)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Set("readBody", readBody); err != nil {
		return nil, err
	}

	writeRes, err := v8go.NewFunctionTemplate(reqCtx.vm, reqCtx.writeRes)
	if err != nil {
		return nil, err
	}
	if err := tmpl.Set("writeRes", writeRes); err != nil {
		return nil, err
	}

	return tmpl.NewInstance(ctx)
}

func (reqCtx *requestContext) headers() (*v8go.ObjectTemplate, error) {
	headers, err := v8go.NewObjectTemplate(reqCtx.vm)
	if err != nil {
		return nil, err
	}

	for k := range reqCtx.req.Header {
		if err := headers.Set(k, reqCtx.req.Header.Get(k)); err != nil {
			return nil, err
		}
	}

	return headers, nil
}

func (reqCtx *requestContext) readBody(info *v8go.FunctionCallbackInfo) *v8go.Value {
	b, err := io.ReadAll(reqCtx.req.Body)
	if err != nil {
		panic(err)
	}

	v, err := v8go.NewValue(reqCtx.vm, string(b))
	if err != nil {
		panic(err)
	}

	return v
}

func (reqCtx *requestContext) writeRes(info *v8go.FunctionCallbackInfo) *v8go.Value {
	// FIXME headers
	body, status := info.Args()[0].String(), info.Args()[1].Int32()

	if status != 0 {
		reqCtx.res.WriteHeader(int(status))
	}

	if _, err := reqCtx.res.Write(([]byte)(body)); err != nil {
		panic(err)
	}

	return nil
}
