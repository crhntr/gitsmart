package gitsmart

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/protocol/packp/capability"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

var (
	errServiceNotFound = errors.New("service not found")
)

func Handle(res http.ResponseWriter, req *http.Request, store storer.Storer) error {
	switch req.URL.Query().Get("service") {
	default:
		http.Error(res, errServiceNotFound.Error(), http.StatusNotFound)

		return errServiceNotFound
	case transport.UploadPackServiceName:
		return handleUploadPackService(res, req, store)

	case transport.ReceivePackServiceName:
		return handleReceivePackService(res, req, store)
	}
}

func handleUploadPackService(res http.ResponseWriter, req *http.Request, store storer.Storer) error {
	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")

	refIter, _ := store.IterReferences()

	advRefs := packp.NewAdvRefs()

	advRefs.Capabilities = capability.NewList()

	if err := refIter.ForEach(func(reference *plumbing.Reference) error {
		return advRefs.AddReference(reference)
	}); err != nil {
		return err
	}

	res.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(res, "001e# service=git-upload-pack\n%s\n", pktline.FlushPkt)
	return advRefs.Encode(res)
}

func handleReceivePackService(res http.ResponseWriter, req *http.Request, store storer.Storer) error {
	return nil
}
