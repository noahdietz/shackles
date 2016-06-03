# shackles
extremely simple/bare bones reverse proxy with custom plugin support in Go lang

### What
shackles is a simple reverse proxy mechanism that accommodates custom plugin chains for
inbound requests and outbound responses. Following the interface defined in [shackles.go](https://github.com/noahdietz/shackles/blob/master/shackles.go)
one can make their own request/response modifiers and build custom chains without having to
change any target application logic.

### Why
I wanted to learn about the networking packages built into the Go language. The project is inspired by
this cool and simple middleware package called [Alice](https://github.com/justinas/alice).

### Example
```go
package main

import (
  "net/http"
  "net/url"
  "log"
  "github.com/noahdietz/shackles"
)

// simple inbound request plugin that modifies header
func reqPluginOne(req *http.Request) (*http.Request, error) {
  log.Println("In request plugin one for", req.URL.Path)
  req.Header.Add("Request-plugin-one", "test1")

  return req, nil
}

// duplicate of first inbound request plugin, to show chaining
func reqPluginTwo(req *http.Request) (*http.Request, error) {
  log.Println("In request plugin two for", req.URL.Path)
  req.Header.Add("Request-plugin-two", "test2")

  return req, nil
}

// simple outbound response plugin that modifies header
func resPluginOne(res *http.Response) (*http.Response, error) {
  log.Println("In response plugin one with status", res.Status)
  res.Header.Add("Response-plugin-one", "test1")

  return res, nil
}

// duplicate of first outbound request plugin, to show chaining
func resPluginTwo(res *http.Response) (*http.Response, error) {
  log.Println("In response plugin two with status", res.Status)
  res.Header.Add("Response-plugin-two", "test2")

  return res, nil
}

func main() {
  // provide URL of target application
  u, _ := url.Parse("http://localhost:8080")

  // build plugin chains
  in := shackles.NewInChain(reqPluginOne, reqPluginTwo)
  out := shackles.NewOutChain(resPluginOne, resPluginTwo)

  // build shackles reverse proxy
  rev := shackles.BuildRev(u, in, out)

  // run shackles reverse proxy on port 3000
  log.Fatal(http.ListenAndServe(":3000", rev))
}
```

Enjoy!