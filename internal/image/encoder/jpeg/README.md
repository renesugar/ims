

# jpeg
`import "github.com/wyattjoh/ims/internal/image/encoder/jpeg"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>



## <a name="pkg-index">Index</a>
* [type Encoder](#Encoder)
  * [func NewEncoder(r *http.Request) Encoder](#NewEncoder)
  * [func (e Encoder) Encode(i image.Image, w http.ResponseWriter) error](#Encoder.Encode)


#### <a name="pkg-files">Package files</a>
[jpeg.go](/src/github.com/wyattjoh/ims/internal/image/encoder/jpeg/jpeg.go) 






## <a name="Encoder">type</a> [Encoder](/src/target/jpeg.go?s=667:703#L21)
``` go
type Encoder struct {
    Quality int
}
```
Encoder allows the encoding of JPEG's to a http.ResponseWriter.







### <a name="NewEncoder">func</a> [NewEncoder](/src/target/jpeg.go?s=390:430#L9)
``` go
func NewEncoder(r *http.Request) Encoder
```
NewEncoder creates a new Encoder based on the input request, this
parses the `q` query variable to check to see if it needs to change the
default quality format.





### <a name="Encoder.Encode">func</a> (Encoder) [Encode](/src/target/jpeg.go?s=777:844#L26)
``` go
func (e Encoder) Encode(i image.Image, w http.ResponseWriter) error
```
Encode writes the encoded image data out to the http.ResponseWriter.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
