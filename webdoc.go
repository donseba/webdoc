package webdoc

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zenazn/goji/web"
)

type Doc struct {
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	In          interface{}       `json:"in,omitempty"`
	Out         interface{}       `json:"out,omitempty"`
	FormValue   map[string]string `json:"form_value,omitempty"`
	URLParams   map[string]string `json:"url_params,omitempty"`
}

type Info struct {
	Routes        map[string]Info `json:"routes,omitempty"`
	Method        string          `json:"method,omitempty"`
	Documentation *Doc            `json:"documentation,omitempty"`
}

type Mux struct {
	webmux *web.Mux        // goji's web.Mux
	DocMap map[string]Info // Doc Map
}

//New creates a new webdoc.Mux
func New() *Mux {
	return &Mux{
		webmux: web.New(),
		DocMap: make(map[string]Info),
	}
}

//Mux returns the goji's *web.Mux
func (m *Mux) Mux() *web.Mux {
	return m.webmux
}

// muxMap adds the route to the webdoc.DocMap
func muxMap(m *Mux, method string, pattern web.PatternType, handler web.HandlerType, doc Doc) {
	sPattern := fmt.Sprintf("%+v", pattern)

	// now we need an string to work with.
	// this string should not contain leading or trailing slashes.
	workString := sPattern
	workString = strings.Trim(workString, "/")

	// it the workstring is empty we are on root.
	if workString == "" {
		workString = "/"
		// make the new root groop if it doesnt exists.
		if _, ok := m.DocMap[workString]; !ok {
			m.DocMap[workString] = Info{
				Method:        strings.ToUpper(method),
				Documentation: &doc,
			}
		}

		return
	}

	parts := strings.Split(workString, "/")
	_ = addParts(m.DocMap, parts, Info{
		Method:        strings.ToUpper(method),
		Documentation: &doc,
	})

}

// muxMap adds the nested routes to the webdoc.DocMap
func handleMap(m *Mux, submux *Mux, pattern web.PatternType) {
	sPattern := fmt.Sprintf("%+v", pattern)

	// check if the group string ends with `/*`
	// if it does, remove that part.
	sPattern = strings.TrimRight(sPattern, "/*")

	// now we need an string to work with.
	// this string should not contain leading or trailing slashes.
	workString := sPattern
	workString = strings.Trim(workString, "/")

	// it the workstring is empty we are on root.
	if workString == "" {
		workString = "/"
		// make the new root groop if it doesnt exists.
		if _, ok := m.DocMap[workString]; !ok {
			m.DocMap[workString] = Info{
				Routes: make(map[string]Info),
			}
		}

		// append the children
		for k, v := range submux.DocMap {
			m.DocMap[workString].Routes[k] = v
		}

		return
	}

	parts := strings.Split(workString, "/")
	target := addParts(m.DocMap, parts, Info{})

	for k, v := range submux.DocMap {
		target.Routes[k] = v
	}
}

func addParts(target map[string]Info, parts []string, info Info) Info {
	key := "/" + parts[0]

	if _, ok := target[key]; !ok {
		if len(parts[1:]) == 0 {
			target[key] = Info{
				Routes:        make(map[string]Info),
				Method:        info.Method,
				Documentation: info.Documentation,
			}
		} else {
			target[key] = Info{
				Routes: make(map[string]Info),
			}
		}
	}

	if len(parts[1:]) > 0 {
		return addParts(target[key].Routes, parts[1:], info)
	}

	return target[key]
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.webmux.ServeHTTP(w, r)
}

func (m *Mux) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	m.webmux.ServeHTTPC(c, w, r)
}

func (m *Mux) Use(middleware web.MiddlewareType) {
	m.webmux.Use(middleware)
}

func (m *Mux) Insert(middleware, before web.MiddlewareType) error {
	return m.webmux.Insert(middleware, before)
}

func (m *Mux) Abandon(middleware web.MiddlewareType) error {
	return m.webmux.Abandon(middleware)
}

func (m *Mux) Router(c *web.C, h http.Handler) http.Handler {
	return m.webmux.Router(c, h)
}

func (m *Mux) Handle(pattern web.PatternType, submux *Mux) {
	handleMap(m, submux, pattern)
	m.webmux.Handle(pattern, submux.Mux())
}

func (m *Mux) Connect(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "connect", pattern, handler, doc)
	m.webmux.Connect(pattern, handler)
}

func (m *Mux) Delete(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "delete", pattern, handler, doc)
	m.webmux.Delete(pattern, handler)
}

func (m *Mux) Get(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "get", pattern, handler, doc)
	m.webmux.Get(pattern, handler)
}

func (m *Mux) Head(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "head", pattern, handler, doc)
	m.webmux.Head(pattern, handler)
}

func (m *Mux) Options(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "options", pattern, handler, doc)
	m.webmux.Options(pattern, handler)
}

func (m *Mux) Patch(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "patch", pattern, handler, doc)
	m.webmux.Patch(pattern, handler)
}

func (m *Mux) Post(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "post", pattern, handler, doc)
	m.webmux.Post(pattern, handler)
}

func (m *Mux) Put(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "put", pattern, handler, doc)
	m.webmux.Put(pattern, handler)
}

func (m *Mux) Trace(pattern web.PatternType, handler web.HandlerType, doc Doc) {
	muxMap(m, "trace", pattern, handler, doc)
	m.webmux.Trace(pattern, handler)
}
