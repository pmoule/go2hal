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
Generated JSON
```
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    }
}
```
To add some resource state, you can add some properties.
These properties must be valid JSON

(Currently, this is not checked)
```go
data := root.Data()
data["doctorCount"] = 12
```
Generated JSON
```
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "doctorCount": 12
}
```
Todo: add embedded resource
```go
actorResources []hal.Resource
// ...
// Skipped populating actorResources.
// ...

doctors, _ := hal.NewResourceRelation("doctor")
doctors.SetResources(actorResources)
```
Generated JSON
```
{
    "_links": {
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "_embedded": {
        "doctors": [
            {
                "_links": {
                    "self": {
                        "href": "/docwhoapi/doctors/1"
                    }
                },
                "name": "William Hartnell"
            },
            {
                "_links": {
                    "self": {
                        "href": "/docwhoapi/doctors/2"
                    }
                },
                "name": "Patrick Troughton"
            },
        ]
    },
    "doctorCount": 12
}
```
Todo: add CURIEs
```go
curieLink, _ := hal.NewCurieLink("doc", "http://example.com/docs/relations/{rel}")
curieLinks := []*hal.LinkObject {curieLink}
root.AddCurieLinks(curieLinks)

doctors.SetCurieLink(curieLink)
```
Generated JSON
```
{
    "_links": {
        "curies": [
            {
                "href": "http://example.com/docs/relations/{rel}",
                "templated": true,
                "name": "doc"
            }
        ]
        "self": {
            "href": "/docwhoapi/doctors"
        }
    },
    "_embedded": {
            "doc:doctors": [
                {
                    "_links": {
                        "self": {
                            "href": "/docwhoapi/doctors/1"
                        }
                    },
                    "name": "William Hartnell"
                },
                {
                    "_links": {
                        "self": {
                            "href": "/docwhoapi/doctors/2"
                        }
                    },
                    "name": "Patrick Troughton"
                },
            ]
        },
    "doctorCount": 12
}
```

##Documentation
See package documentation:

[![GoDoc](https://godoc.org/github.com/pmoule/go2hal/hal?status.svg)](https://godoc.org/github.com/pmoule/go2hal/hal)

## Todo
- provide better examples for usage
- howto: download