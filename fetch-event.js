class FetchEvent {
  constructor(request, writeRes) {
    this.request = request
    this.__writeRes = writeRes
  }

  respondWith(responsePromise) {
    Promise.resolve(responsePromise).then(response => {
      // FIXME headers
      this.__writeRes(response.body, response.status)
    })
    // FIXME catch
  }
}
