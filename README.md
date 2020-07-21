# GitSmart (very early WIP)

An http.Handler type thing that works well with go-git.
The name is based on the name of the more robust/modern HTTP protocal used in the Git plumbing/tranfer stuff.

## Notes

- Introductory notes on git Plumbing https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain
  - Intro to transfer protocols: https://git-scm.com/book/en/v2/Git-Internals-Transfer-Protocols
- Somewhat deeper documentation of the git "Smart HTTP" protocol https://github.com/git/git/blob/master/Documentation/technical/http-protocol.txt
  - Some of the symbols in the specification are documented here: https://github.com/git/git/blob/master/Documentation/technical/protocol-common.txt

## Design 

I would like the API to be something like this.

```go
package app

import (
	"net/http"

	"github.com/crhntr/gitsmart"
)

var db *Database

func MyGitHandler(res http.ResponseWriter, req *http.Request) {
	s := db.Session()
	s.StartTransaction()

	userAuth, err := loadUserAuthorizationFromWebToken(req)
	if err != nil {
		http.Error(res, "repository not found", http.StatusNotAuthorized)
		return
	}

	// note any store that implements github.com/go-git/go-git/v5/plumbing/storer would be permitted.
	store, err := loadGitStore(db, req.URL.Path)
	if err != nil {
		http.Error(res, "repository not found", http.StatusNotFound)
		return
	}

	err = gitsmart.Handle(res, req, store,
		// Handler functionality can be modified with functional options
		// see: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
		//
		// For example, we can implement our own access control for individual branches...
		// or maybe even more fine grained control over git objects
		gitsmart.BeforeBranchUpdate(CheckIfUserCanUpdateBranch(&userAuth)),
		//
		// Another example, because we know our store can handle transactional updates, the
		// Handler can be notified to broadcast that functionality to clients.
		gitsmart.BroadcastTransactionalCapability,
	)

	if err != nil {
		s.AbortTransaction()
		return
	}

	s.CommitTransaction()
}

func CheckIfUserCanUpdateBranch(user *UserAuth) gitsmart.BeforeBranchUpdateFunc {
	return func(ref *plumbing.Reference) error {
		return nil /* check if user allowed to edit reference otherwise return an error */
	}
}
```
