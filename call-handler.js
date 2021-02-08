(() => {
  const { url, method, headers, readBody, writeRes } = reqCtx
  const e = new FetchEvent(new Request(url, { method, headers, readBody }), writeRes)
  delete reqCtx
  
  // FIXME manage reject if async
  handler(e, writeRes)
})()
