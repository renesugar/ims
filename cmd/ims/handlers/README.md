

# handlers
`import "github.com/wyattjoh/ims/cmd/ims/handlers"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package handlers provides the routes used by the `ims` binary.




## <a name="pkg-index">Index</a>
* [func Image(timeout time.Duration) http.HandlerFunc](#Image)


#### <a name="pkg-files">Package files</a>
[handlers.go](/src/github.com/wyattjoh/ims/cmd/ims/handlers/handlers.go) 





## <a name="Image">func</a> [Image](/src/target/handlers.go?s=1042:1092#L24)
``` go
func Image(timeout time.Duration) http.HandlerFunc
```
Image is the handler which loads the filename from the request, loads the
file via the provider, and processes the image to re-encode it with caching
headers.








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
