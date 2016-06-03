package shackles

import (
  "net/http"
)

// interface for inbound/outbound request plugins
type InLink func(req *http.Request) (*http.Request, error)
type OutLink func(res *http.Response) (*http.Response, error)

// typed array for inbound/outbound plugin sequence
type InChain []InLink
type OutChain []OutLink

// construct holding plugin sequences
type Shackles struct {
  inbound InChain
  outbound OutChain
}

// variadic constructor for inbound plugin sequence
func NewInChain(in ...InLink) InChain {
  return append(([]InLink)(nil), in...)
}

// variadic constructor for outbound plugin sequence
func NewOutChain(out ...OutLink) OutChain {
  return append(([]OutLink)(nil), out...)
}