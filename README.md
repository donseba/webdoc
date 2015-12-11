# webdoc
An [goji](https://github.com/goji/goji) Mux wrapper that allows you to add documentation to routes. 

Please not to `go get goji.io`

This package is an simple wrapper around goji's mux struct, adding functionality to document the API routes.

**Please keep in mind that this package in still under development

## Example
``` go
import (
    "fmt"
    "net/http"
    
    "goji.io"
    "goji.io/pat"

    "github.com/donseba/webdoc"
)

var Doc map[string]webdoc.Info


//hello will say hello to :name
func hello(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    name := pat.Param(ctx, "name")
    fmt.Fprintf(w, "Hello, %s!", name)
}

//routes will output an json of routes
func routes(c context.Context, w http.ResponseWriter, r *http.Request) {
    out, err := json.Marshal(core.DocMap)
    if err != nil {
        panic(err)
    }

    fmt.Fprint(w, out)
}

func main() {
    // instead of web.New()
    wd := webdoc.New()
    
    wd.Get("/hello/:name", hello, webdoc.Doc{Title: "Say hello", Description: "Say hello to :name"})
    wd.Get("/routes", routes, webdoc.Doc{Title: "API Routes", Description: "Retrieve the list of API Routes in JSON format"})
    
    // Assing doc to an global variable
    Doc = wd.DocMap

    mux := goji.NewMux()
    mux.Handle(pat.New("/*"), wd.Mux())
        
    http.ListenAndServe("localhost:8080", mux)
}
```

Running the example and pointing the browser to `/routes`

Would return the following :
``` json
{
	"/": {
		"methods": {
			"GET": {
				"title": "Index",
				"description": "Index Page"
			}
		}
	},
	"/hello": {
		"routes": {
			"/:name": {
				"methods": {
					"GET": {
						"documentation": {
							"title": "Say Hello",
							"description": "Say hello to :name"
						}
					},
					"PUT": {
						"documentation": {
							"title": "Put some data",
							"description": "Say hello to :name"
						}
					}
				}
			}
		}
	},
	"/routes": {
		"methods": {
			"GET": {
				"Title": "API Routes",
				"Description": "Retrieve the list of API Routes in JSON format"
			}
		}
	}
}
```
