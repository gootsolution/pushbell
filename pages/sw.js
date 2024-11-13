self.addEventListener("install", event => {
    event.waitUntil(self.skipWaiting())
})

self.addEventListener("activate", event => {
    event.waitUntil(self.clients.claim())
})

self.addEventListener("push", event => {
    console.log("[SW]: ", event.data.text())
})