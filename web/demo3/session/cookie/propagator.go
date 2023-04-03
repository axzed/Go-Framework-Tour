package cookie

import (
	"net/http"
)

type CookieOption func(cookie *http.Cookie)

type PropagatorOption func(p *Propagator)

type Propagator struct {
	cookieName string
	cookieOpt  CookieOption
}

func WithCookieOption(opt CookieOption) PropagatorOption {
	return func(p *Propagator) {
		p.cookieOpt = opt
	}
}

func NewPropagator(cookieName string, opts ...PropagatorOption) *Propagator {
	res := &Propagator{
		cookieName: cookieName,
		cookieOpt: func(cookie *http.Cookie) {

		},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	c, err := req.Cookie(p.cookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func (p *Propagator) Inject(id string, resp http.ResponseWriter) error {
	c := &http.Cookie{
		Name:  p.cookieName,
		Value: id,
	}
	p.cookieOpt(c)
	http.SetCookie(resp, c)
	return nil
}

func (p *Propagator) Remove(resp http.ResponseWriter) error {
	http.SetCookie(resp, &http.Cookie{
		Name:   p.cookieName,
		MaxAge: -1,
	})
	return nil
}
