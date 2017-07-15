# go2hal
A **Go** implementation of **Hypertext Application Language (HAL)**.
It provides essential data structures and features a JSON generator
to produce **JSON** output as proposed in [JSON Hypertext Application Language](https://tools.ietf.org/html/draft-kelly-json-hal).

## Features
- HAL API
    - Create root **Resource Object**.
    - Supports "_links" property.

        Define **Link Relations** and assign **Link Object** value(s).

    - Supports "_embedded" property.

        Define **Link Relations** and assign **Resource Object** value(s).
    - Supports "curies".

        Define CURIE **Link Objects** and assign to defined **Link Relations**.
- JSON generator to produce HAL Document


## Usage
### Preliminary stuff
Download and install `go2hal` into your GOPATH.
```
go get github.com/pmoule/go2hal
```
Import the `hal` package to get started.
```go
import "github.com/pmoule/go2hal/hal"
```
### Create a Resource and JSON generation
First create a Resource Object as your HAL document's root element.
```go
root := hal.NewResourceObject()
```
This is all you need, to create a valid HAL document.

Next, create an `Encoder` and call it's `ToJSON` function to generate valid JSON.
```go
encoder := hal.NewEncoder()
bytes, error := encoder.ToJSON(root)
```
The generated JSON is
```
{}
```
There's potential for more :smile:
### Add a Link Relation
So let's add a `self` Link Relation. Additionally we attach a single Link Object.
```go
link := &hal.LinkObject{ Href: "/docwhoapi/doctors"}

self, _ := hal.NewLinkRelation("self") //skipped error handling
self.SetLink(link)

root.AddLink(self)
```
Since `self` is a IANA registered link relation name, **go2hal** provides a shortcut
```go
self := hal.NewSelfLinkRelation()
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
### Resource state
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
### Embedding Resources
Now, let's embed some resources.
```go
// a simple struct for actors of the doctor
type Actor struct {
    Id int
    Name string
}

actors := []Actor {
    Actor{1, "William Hartnell"},
    Actor{2, "Patrick Troughton"},
}

// convert the actors to resources
var embeddedActors []hal.Resource

for _, actor := range actors {
    href := fmt.Sprintf("/docwhoapi/doctors/%d", actor.Id)
    selfLink, _ := hal.NewLinkObject(href)

    self, _ := hal.NewLinkRelation("self")
    self.SetLink(selfLink)

    embeddedActor := hal.NewResourceObject()
    embeddedActor.AddLink(self)
    embeddedActor.Data()["name"] = actor.Name
    embeddedActors = append(embeddedActors, embeddedActor)
}

doctors, _ := hal.NewResourceRelation("doctors")
doctors.SetResources(embeddedActors)

root.AddResource(doctors)
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
            }
        ]
    },
    "doctorCount": 12
}
```
### CURIEs
A Resource Object can have a set of CURIE links. Same for used Link Relations, that are capable of setting
a CURIE link.
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
                }
            ]
        },
    "doctorCount": 12
}
```
### Relations and the array vs single value discussion
I'm aware of existing discussions regarding Relations and the type of assigned values.
I simply deal with this topic by leaving the decision to the developer whether to
assign a single value or an array value.

**go2hal** provides a `LinkRelation` and a `ResourceRelation`. Both are capable of holding
a single value or an array value by providing special setter functions.

**Example**
```go
link := hal.NewLinkObject("myHref")
linkRelation := hal.NewLinkRelation("linkRel")

resource := hal.NewResourceObject()
resource.Data()["value"] = "myValue"
resourceRelation := hal.NewResourceRelation("resourceRel")
```
Assign single values
```go
linkRelation.SetLink(link)
resourceRelation.SetResource(resource)
```
Generated JSON
```
{
    "_links": {
        "linkRel": {
            "href": "myHref"
        }
    },
    "_embedded": {
            "resourceRel": {
                "value": "myValue"
            }
        },
}
```
Assign array values
```go
linkRelation.SetLinks([]*LinkObject{link})
resourceRelation.SetResources([]hal.Resource{resource})
```
Generated JSON
```
{
    "_links": {
        "linkRel": [
            {
                "href": "myHref"
            }
        ]
    },
    "_embedded": {
            "resourceRel": [
                {
                    "value": "myValue"
                }
            ]
        },
}
```
**go2hal** does not evaluate the assigned values. The developer is fully responsible of this.

CURIEs are an exception to this. As stated in the [specification](https://tools.ietf.org/html/draft-kelly-json-hal#section-8.2)
CURIEs are always an array of Link Objects.

## Documentation
See package documentation:

[![GoDoc](https://godoc.org/github.com/pmoule/go2hal/hal?status.svg)](https://godoc.org/github.com/pmoule/go2hal/hal)

## License
`go2hal` is released under MIT license. See [LICENSE](LICENSE.txt).