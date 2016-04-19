# go2hal
A **Go** implementation of **Hypertext Application Language (HAL)**.
It provides essential data structures and features a JSON generator
to produce **JSON** output as proposed in [JSON Hypertext Application Language](https://tools.ietf.org/html/draft-kelly-json-hal-07).

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
Will provide good examples, soon.
To get started, create a root `ResourceObject` and try a bit :)


```go
import "github.com/pmoule/go2hal/hal"
...
root := hal.NewResourceObject()
selfLinkRelation := hal.NewLinkRelation("self")
selfLinkObject := &hal.LinkObject{ Href: "http://example.com/api/items"}
halResource.AddLinkObject(selfLinkRelation, selfLinkObject)
...
encodedJSON, err := root.ToJson()

```


## Todo
- provide better examples for usage
- howto: download