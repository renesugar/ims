package ims

import (
	"fmt"
	"image"
	"io"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/disintegration/imaging"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
)

const (

	// timeout is the cache timeout used to add to requests to prevent the
	// browser from re-requesting the image.
	timeout = 15 * time.Minute
)

// ProcessImage uses the github.com/disintegration/imaging lib to perform the
// image transformations.
func ProcessImage(input io.Reader, w http.ResponseWriter, r *http.Request) error {
	srcImage, format, err := image.Decode(input)
	if err != nil {
		return errors.Wrap(err, "can't decode the image")
	}

	width, err := strconv.Atoi(r.URL.Query().Get("w"))
	if err == nil {
		srcImage = imaging.Resize(srcImage, width, 0, imaging.Linear)
	}

	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int64(timeout.Seconds())))
	w.Header().Set("Last-Modified", time.Now().String())

	encoder := GetEncoder(format, r)
	if err := encoder.Encode(srcImage, w); err != nil {
		return errors.Wrap(err, "can't encode the image")
	}

	return nil
}

// GetFilename fetches the filename from the request path.
func GetFilename(r *http.Request) (string, error) {

	// We expect that the router sends us requests in the form `/resize/:filename`
	// so we check to see if the path contains the image url that we want to
	// parse. In this case, we check to see that the path is at least 9 characters
	// long, which will ensure that the filename has at least 1 character.
	if len(r.URL.Path) < 9 {
		return "", errors.New("filename too short")
	}

	return r.URL.Path[8:], nil
}

// HandleFileSystemResize performs the actual resizing by loading the image
// from the filesystem.
func HandleFileSystemResize(dir http.Dir) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Extract the filename from the request.
		filename, err := GetFilename(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Try to open the image from the virtual filesystem.
		f, err := dir.Open(filename)
		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// If an error occurred during the image processing, return with an internal
		// server error.
		if err := ProcessImage(f, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// HandleOriginResize performs the actual resizing by loading the image
// from the origin.
func HandleOriginResize(origin string) (http.HandlerFunc, error) {
	originURL, err := url.Parse(origin)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse the origin url")
	}

	return func(w http.ResponseWriter, r *http.Request) {

		// Extract the filename from the request.
		filename, err := GetFilename(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Parse the incomming url.
		filenameURL, err := url.Parse(filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Resolve it relative to the origin url.
		fileURL := originURL.ResolveReference(filenameURL)

		// Perform the GET to the origin server.
		req, err := http.NewRequest("GET", fileURL.String(), nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer res.Body.Close()

		// If an error occurred during the image processing, return with an internal
		// server error.
		if err := ProcessImage(res.Body, w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}, nil
}

// Serve creates and starts a new server to provide image resizing services.
func Serve(addr string, debug bool, directory, origin string) error {
	mux := http.NewServeMux()

	// By default, we'll try to use the directory resize, otherwise, if the origin
	// url is provided, use it.
	if origin == "" {
		logrus.WithField("directory", directory).Debug("serving from the filesystem")
		mux.HandleFunc("/resize/", HandleFileSystemResize(http.Dir(directory)))
	} else {
		logrus.WithField("origin", origin).Debug("serving from the origin")

		handler, err := HandleOriginResize(origin)
		if err != nil {
			return errors.Wrap(err, "can't create origin resize handler")
		}

		mux.HandleFunc("/resize/", handler)
	}

	// When debug mode is enabled, mount the debug handlers on this router.
	if debug {
		mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	}

	n := negroni.Classic() // Includes some default middlewares

	n.UseHandler(mux)

	logrus.Debugf("Now listening on %s", addr)
	return http.ListenAndServe(addr, n)
}
