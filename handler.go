package gitsmart

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/pktline"
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
	w := pktline.NewEncoder(res)

	res.Header().Set("Cache-Control", "no-cache")
	res.Header().Set("Content-Type", "application/x-git-upload-pack-advertisement")

	if err := w.EncodeString("# service=git-upload-pack\n"); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}
	_, _ = res.Write([]byte{'\n'})


	refIter, _ := store.IterReferences()

	var references []plumbing.Reference

	if err := refIter.ForEach(func(reference *plumbing.Reference) error {
		references = append(references, *reference)
		return nil
	}); err != nil {
		return err
	}

	capabilities := []string{"some-capability"}

	if len(references) == 0 {
		if err := w.Encodef("%s capabilities^{}\u0000%s\n", plumbing.ZeroHash, strings.Join(capabilities, " ")); err != nil {
			return err
		}

		if err := w.Flush(); err != nil {
			return err
		}

		return nil
	}

	for i, ref := range references {
		suffix := "\n"

		if i == len(references)-1 {
			suffix += "^{}"
		}

		if i == 0 {
			suffix = "\u0000" + strings.Join(capabilities, " ") + suffix
		}

		if err := w.Encodef("%s %s%s"+suffix, ref.Hash(), ref.Hash(), suffix); err != nil {
			return err
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}


	return nil
}

func handleReceivePackService(res http.ResponseWriter, req *http.Request, store storer.Storer) error {
	return nil
}
