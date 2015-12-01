# webdoc
An [goji](https://github.com/zenazn/goji) web.Mux wrapper that allows you to add documentation to routes.

This package is an simple wrapper around goji's web structure, adding functionality to document the API routes.

**Please keep in mind that this package in still under development

## Example
``` go
import (
        "fmt"
        "net/http"

        "github.com/donseba/webdoc"
        "github.com/zenazn/goji"
        "github.com/zenazn/goji/web"
)

var Doc map[string]webdoc.Info

// This is an simple documantation line for the function hello
// Hello will output the string after hello/
func hello(c web.C, w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
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
        wd.Get("/hello/:banana", hello, webdoc.Doc{Title: "Say hello", Description: "Say hello to :banana"})
        wd.Get("/routes", routes, webdoc.Doc{Title: "API Routes", Description: "Retrieve the list of API Routes in JSON format"})

		// Assing doc to an global variable
		Doc = wd.DocMap


        goji.Handle("/*", wd.Mux() )
        goji.Serve()
}
```

Running the example and pointing the browser to `/routes`

Would return the following :
``` json
{
	"/hello": {
		"routes": {
			"/:banana": {
				"method": "GET",
				"documentation": {
					"title": "Say Hello",
					"description": "Say hello to :banana"
				}
			},
			"/:name": {
				"method": "GET",
				"documentation": {
					"title": "Say Hello",
					"description": "Say hello to :name"
				}
			}
		}
	},
	"/routes": {
		"method": "GET",
		"documentation": {
			"Title": "API Routes",
			"Description": "Retrieve the list of API Routes in JSON format"
		}
	}
}

```