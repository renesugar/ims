
# jpeg
    import "github.com/wyattjoh/ims/internal/image/encoder/jpeg"







## type Encoder
``` go
type Encoder struct {
    Quality int
}
```
Encoder allows the encoding of JPEG's to a http.ResponseWriter.









### func NewEncoder
``` go
func NewEncoder(r *http.Request) Encoder
```
NewEncoder creates a new Encoder based on the input request, this
parses the `q` query variable to check to see if it needs to change the
default quality format.




### func (Encoder) Encode
``` go
func (e Encoder) Encode(i image.Image, w http.ResponseWriter) error
```
Encode writes the encoded image data out to the http.ResponseWriter.









- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)