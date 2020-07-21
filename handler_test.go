package gitsmart

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	Ω "github.com/onsi/gomega"
)

func TestHandle(t *testing.T) {
	t.Run("when the service is unknown", func(t *testing.T) {
		please := Ω.NewWithT(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/repo-name/info/refs?service=NOT_A_SERVICE", nil)

		err := Handle(res, req, nil)

		please.Expect(err).To(Ω.HaveOccurred())
		please.Expect(err).To(Ω.MatchError(Ω.ContainSubstring("service not found")))

		please.Expect(res.Code).To(Ω.Equal(http.StatusNotFound), "the HTTP status code should be status not found")
		please.Expect(res.Body).To(Ω.ContainSubstring("service not found"), "the response body should indicate something about the service not being known")
	})

	t.Run("when the the repo is empty", func(t *testing.T) {
		please := Ω.NewWithT(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/repo-name/info/refs?service=git-upload-pack", nil)

		repo, _ := git.Init(memory.NewStorage(), memfs.New())
		err := Handle(res, req, repo.Storer)

		please.Expect(err).NotTo(Ω.HaveOccurred())

		please.Expect(res.Code).To(Ω.Equal(http.StatusOK))

		buf, _ := ioutil.ReadAll(res.Body)
		lines := bytes.Split(buf, []byte("\n"))

		please.Expect(lines).To(Ω.HaveLen(4))
		please.Expect(lines[0]).To(Ω.Equal([]byte("001e# service=git-upload-pack")))
		please.Expect(lines[1]).To(Ω.Equal([]byte("0000")))
		please.Expect(lines[2][4:]).To(Ω.Equal([]byte("0000000000000000000000000000000000000000 capabilities^{}\x00")))
		please.Expect(lines[3]).To(Ω.Equal([]byte("0000")))
	})

	t.Run("when one branch exists", func(t *testing.T) {
		please := Ω.NewWithT(t)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/repo-name/info/refs?service=git-upload-pack", nil)

		repo, _ := git.Init(memory.NewStorage(), memfs.New())

		initialCommitHash := helloWorldReadMe(repo)

		err := Handle(res, req, repo.Storer)

		please.Expect(err).NotTo(Ω.HaveOccurred())

		please.Expect(res.Code).To(Ω.Equal(http.StatusOK))

		buf, _ := ioutil.ReadAll(res.Body)
		lines := bytes.Split(buf, []byte("\n"))

		please.Expect(lines).To(Ω.HaveLen(4))
		please.Expect(lines[0]).To(Ω.Equal([]byte("001e# service=git-upload-pack")))
		please.Expect(lines[1]).To(Ω.Equal([]byte("0000")))
		please.Expect(lines[2][4:]).To(Ω.Equal([]byte(fmt.Sprintf("%s refs/heads/master^{}\x00some-capability\n", initialCommitHash))))
		please.Expect(lines[3]).To(Ω.Equal([]byte("0000")))
	})
}


func helloWorldReadMe(repo *git.Repository) plumbing.Hash {
	wt, _ := repo.Worktree()
	md, _ := wt.Filesystem.Create("README.md")
	_, _ = md.Write([]byte("Hello, world!\n"))
	_ = md.Close()
	initialCommit, _ := wt.Commit("initial", &git.CommitOptions{
		All: true,
		Author: &object.Signature{Name: "person", Email: "person@example.com", When: time.Unix(1595300154, 0)},
		Committer: &object.Signature{Name: "person", Email: "person@example.com", When: time.Unix(1595300154, 0)},
	})
	return initialCommit
}



/* Github's capabilities
multi_ack
thin-pack
side-band
side-band-64k
ofs-delta
shallow
deepen-since
deepen-not
deepen-relative
no-progress
include-tag
multi_ack_detailed
allow-tip-sha1-in-want
allow-reachable-sha1-in-want
no-done
symref=HEAD:refs/heads/master
filter
agent=git/github
*/
