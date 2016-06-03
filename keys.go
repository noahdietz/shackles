package shackles

import (
  "net/http"
  "net/http/httputil"
  "net/url"
)

// given a target and chains of inbound/outbound plugins, we make a ReverseProxy
func BuildRev(target *url.URL, inbound InChain, outbound OutChain) *httputil.ReverseProxy {
  rev := httputil.NewSingleHostReverseProxy(target)

  // override ReverseProxy.Transport default with ours
  rev.Transport = &Shackles{inbound, outbound}

  // return it for use
  return rev
}

// provide the Transport.RoundTrip function to handle requests the proxy receives
func (s *Shackles) RoundTrip(req *http.Request) (resp *http.Response, err error) {
  // apply inbound request plugins (if any)
  req, err = s.applyInboundLinks(req)
  if err != nil {
    return nil, err
  }

  // use default RoundTrip function handle the actual request/response
  resp, err = http.DefaultTransport.RoundTrip(req)
  if err != nil {
    return nil, err
  }

  // block until the entire body has been read
  _, err = httputil.DumpResponse(resp, true)
  if err != nil {
    return nil, err
  }

  // apply outbound response plugins (if any)
  resp, err = s.applyOutboundLinks(resp)
  if err != nil {
    return nil, err
  }

  // pass response back to client
  return resp, nil
}

// applies any outbound response plugins to the given response
func (s *Shackles) applyOutboundLinks(resp *http.Response) (mod *http.Response, err error) {
  mod = resp

  for o := range s.outbound {
    mod, err = s.outbound[o](mod)
    if err != nil {
      return nil, err
    }
  }

  return mod, nil
}

// applies any inbound request plugins to the given request
func (s *Shackles) applyInboundLinks(req *http.Request) (mod *http.Request, err error) {
  mod = req

  for i := range s.inbound {
    mod, err = s.inbound[i](mod)
    if err != nil {
      return nil, err
    }
  }

  return mod, nil
}
