# encoder [![wercker status](https://app.wercker.com/status/170727695eeb0c8fef3220cd7585c855 "wercker status")](https://app.wercker.com/project/bykey/170727695eeb0c8fef3220cd7585c855)

This is a simple wrapper to json.Marshal and xml.Marshal, which adds ability to skip some fields of structure. Unlike 'render' package it doesn't write anything, just returns marshalled data. It's useful for things like passwords, statuses, activation codes, etc... 

E.g.:

```go
type Some struct {
	Login    string        `json:"login"`
	Password string        `json:"-"`
}
```

Field 'Password' won't be exported.

#### Usage.
It's pretty straightforward:

```go
m.Use(func(c martini.Context, w http.ResponseWriter) {
	c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
})
```

Here is a ready to use example:

```go
package main

import (
	"github.com/martini-contrib/encoder"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

type Some struct {
	Login    string `json:"login"`
	Password string `json:"-"`
}

func main() {
	m := martini.New()
	route := martini.NewRouter()

	// map json encoder
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	route.Get("/test", func(enc encoder.Encoder) (int, []byte) {
		result := &Some{"awesome", "hidden"}
		return http.StatusOK, encoder.Must(enc.Encode(result))
	})

	m.Action(route.Handle)

	log.Println("Waiting for connections...")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}
```

XML:

```go
package main

import (
	"github.com/martini-contrib/encoder"
	"github.com/codegangsta/martini"
	"log"
	"net/http"
)

type Some struct {
	Login    string `xml:"login"`
	Password string `xml:"-"`
}

func main() {
	m := martini.New()
	route := martini.NewRouter()

	// map xml encoder
	m.Use(func(c martini.Context, w http.ResponseWriter) {
		c.MapTo(encoder.XMLEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	})

	route.Get("/test", func(enc encoder.Encoder) (int, []byte) {
		result := &Some{"awesome", "hidden"}
		return http.StatusOK, encoder.Must(enc.Encode(result))
	})

	m.Action(route.Handle)

	log.Println("Waiting for connections...")

	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}
```