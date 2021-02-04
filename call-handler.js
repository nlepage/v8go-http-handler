const { url, method, headers, bodyReader, callback } = info
const e = new FetchEvent(new Request(url, { method, headers, bodyReader }), callback)
delete info

// FIXME manage reject if async
handler(e)
