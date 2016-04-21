# go2hal
A **Go** implementation of **Hypertext Application Language (HAL)**.
It provides essential data structures and features a JSON generator
to produce **JSON** output as proposed in [JSON Hypertext Application Language](https://tools.ietf.org/html/draft-kelly-json-hal).

##Features
- HAL API
    - Create root **Resource Object**.
    - Supports "_links" property

        Define **Link Relations** and assign **Link Object** value(s).

    - Supports "_embedded" property:

        Define **Link Relations** and assign **Resource Object** value(s).
    - Supports "curies"

        Define CURIE **Link Objects** and assign to defined **Link Relations**.
- JSON generator to produce HAL Document


##Usage
Import the `hal` package to get started.
```go
import "github.com/pmoule/go2hal/hal"
```
First create a Resource Object as your HAL document's root element.
```go
root := hal.NewResourceObject()
```
This is all you need, to create a valid HAL document.
```go
encoder := new(hal.Encoder)
bytes, error := encoder.ToJSON(root)
```
The generated JSON is
```
{}
```
There's potential for more :smile:
So let's add a `self` Link Relation.
```go
link := &hal.LinkObject{ Href: "/docwhoapi/doctors"}

self, _ := hal.NewLinkRelation("self") //skipped error handling
self.SetLink(link)

root.AddLink(self)
```
This is the generated JSON
```
{
    "_links": {
        "self": { "href": "/docwhoapi/doctors" }
    }
}
```

##Documentation
See package documentation:

[![GoDoc](https://godoc.org/github.com/pmoule/go2hal/hal?status.svg)](https://godoc.org/github.com/pmoule/go2hal/hal)

## Todo
- provide better examples for usage
- howto: download