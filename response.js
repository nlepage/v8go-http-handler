class Response {
  constructor(body, { status, headers } = {}) {
    if (typeof body !== 'string') throw new Error('not implemented')
    this.body = body
    this.status = status
    // FIXME should be a Headers instance
    this.headers = headers
  }
}
