package webdoc

import (
	"fmt"
	"net/http"
	"strings"

	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
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
	Routes  map[string]Info `json:"routes,omitempty"`
	Methods map[string]*Doc `json:"methods,omitempty"`
}

type Mux struct {
	webmux *goji.Mux       // goji's web.Mux
	DocMap map[string]Info // Doc Map
}

//New creates a new webdoc.Mux
func New() *Mux {
	return &Mux{
		webmux: goji.NewMux(),
		DocMap: make(map[string]Info),
	}
}

// NewSub creates a new goji.SubMux()
func NewSub() *Mux {
	return &Mux{
		webmux: goji.SubMux(),
		DocMap: make(map[string]Info),
	}
}


//Mux returns the goji's *web.Mux
func (m *Mux) Mux() *goji.Mux {
	return m.webmux
}

// muxMap adds the route to the webdoc.DocMap
func muxMap(m *Mux, method string, pattern string, handler goji.Handler, doc *Doc) {
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
				Methods: map[string]*Doc{strings.ToUpper(method): doc},
			}
		} else {
			ma := make(map[string]*Doc)
			for k, v := range m.DocMap[workString].Methods {
				ma[k] = v
			}
			ma[strings.ToUpper(method)] = doc

			m.DocMap[workString] = Info{
				Methods: ma,
			}
		}

		return
	}

	parts := strings.Split(workString, "/")

	// Add missing URL Param to doc.
	// This will save some typing.
	for _, v := range parts {
		if string(v[0]) == ":" {
			if doc == nil {
				doc = &Doc{
					URLParams: make(map[string]string),
				}
			}

			if doc.URLParams == nil {
				doc.URLParams = make(map[string]string)
			}

			if "" == doc.URLParams[v] {
				doc.URLParams[v[1:]] = "string"
			}
		}
	}

	_ = addParts(m.DocMap, parts, strings.ToUpper(method), doc)
}

// muxMap adds the nested routes to the webdoc.DocMap
func handleMap(m *Mux, submux *Mux, pattern string) {
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
	target := addParts(m.DocMap, parts, "", nil)

	for k, v := range submux.DocMap {
		target.Routes[k] = v
	}
}

func addParts(target map[string]Info, parts []string, method string, doc *Doc) Info {
	key := "/" + parts[0]

	if _, ok := target[key]; !ok {
		if len(parts[1:]) == 0 {
			target[key] = Info{
				Routes:  make(map[string]Info),
				Methods: map[string]*Doc{strings.ToUpper(method): doc},
			}
		} else {
			target[key] = Info{
				Routes: make(map[string]Info),
			}
		}
	} else {
		if len(parts[1:]) == 0 {
			ma := make(map[string]*Doc)
			for k, v := range target[key].Methods {
				ma[k] = v
			}

			ma[method] = doc

			target[key] = Info{
				Routes:  make(map[string]Info),
				Methods: ma,
			}
		}
	}

	if len(parts[1:]) > 0 {
		return addParts(target[key].Routes, parts[1:], method, doc)
	}

	return target[key]
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.webmux.ServeHTTP(w, r)
}

func (m *Mux) ServeHTTPC(c context.Context, w http.ResponseWriter, r *http.Request) {
	m.webmux.ServeHTTPC(c, w, r)
}

func (m *Mux) Use(middleware func(http.Handler) http.Handler) {
	m.webmux.Use(middleware)
}

func (m *Mux) UseC(middleware func(goji.Handler) goji.Handler) {
	m.webmux.UseC(middleware)
}

func (m *Mux) Handle(pattern string, submux *Mux) {
	handleMap(m, submux, pattern)
	m.webmux.Handle(pat.New(pattern), submux.Mux())
}

func (m *Mux) Delete(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "delete", pattern, handler, doc)
	m.webmux.HandleC(pat.Delete(pattern), handler)
}

func (m *Mux) Get(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "get", pattern, handler, doc)

	m.webmux.HandleC(pat.Get(pattern), handler)
}

func (m *Mux) Head(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "head", pattern, handler, doc)
	m.webmux.HandleC(pat.Head(pattern), handler)
}

func (m *Mux) Options(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "options", pattern, handler, doc)
	m.webmux.HandleC(pat.Options(pattern), handler)
}

func (m *Mux) Patch(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "patch", pattern, handler, doc)
	m.webmux.HandleC(pat.Patch(pattern), handler)
}

func (m *Mux) Post(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "post", pattern, handler, doc)
	m.webmux.HandleC(pat.Post(pattern), handler)
}

func (m *Mux) Put(pattern string, handler goji.Handler, doc *Doc) {
	muxMap(m, "put", pattern, handler, doc)
	m.webmux.HandleC(pat.Put(pattern), handler)
}
