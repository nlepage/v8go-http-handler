class FetchEvent {
  constructor(request, callback) {
    this.request = request
    this.callback = callback
  }

  async respondWith(responsePromise) {
    const response = await responsePromise
    // FIXME headers
    this.callback(response.body, response.status, response.statusText)
  }
}
