(async () => {
  const { url, method, headers, readBody, writeRes } = reqCtx
  const e = new FetchEvent(new Request(url, { method, headers, readBody }), writeRes)
  delete reqCtx
  
  try {
    await handler(e, writeRes)
  } catch (err) {
    writeRes(`${err}`, 500)
  }
})()
