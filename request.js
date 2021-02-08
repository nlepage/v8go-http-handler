class Request {
  constructor(url, { method = 'GET', headers = {}, readBody }) {
    this.url = url
    this.method = method
    // FIXME should be a Headers instance
    this.headers = headers
    this.__readBody = readBody
  }

  async text() {
    return this.__readBody()
  }

  async json() {
    return JSON.parse(this.__readBody())
  }

  async arrayBuffer() {
    return Uint8Array.from(this.__readBody(), c => c.charCodeAt(0)).buffer
  }

  formData() {
    throw new Error('not implemented')
  }
}
