package gitsmart

import (
	"net/http"

	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func Handle(res http.ResponseWriter, req *http.Request, store storer.Storer) {
	service := req.URL.Query().Get("service")

	switch service {
	default:
		http.Error(res, "git service not found", http.StatusNotFound)
	case transport.UploadPackServiceName:
		handleUploadPackService(res, req, store)
	case transport.ReceivePackServiceName:
		handleReceivePackService(res, req, store)
	}
}

func handleUploadPackService(res http.ResponseWriter, req *http.Request, store storer.Storer) {

}

func handleReceivePackService(res http.ResponseWriter, req *http.Request, store storer.Storer) {

}


