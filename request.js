class Request {
  constructor(url, { method = 'GET', headers = {}, bodyReader }) {
    this.url = url
    this.method = method
    // FIXME should be a Headers instance
    this.headers = headers
    this.__bodyReader = bodyReader
  }

  async text() {
    return this.__bodyReader()
  }

  async json() {
    return JSON.parse(this.__bodyReader())
  }

  async arrayBuffer() {
    return Uint8Array.from(this.__bodyReader(), c => c.charCodeAt(0)).buffer
  }

  formData() {
    throw new Error('not implemented')
  }
}
