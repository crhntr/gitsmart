package gitsmart

import (
	"net/http"
	"net/http/httptest"
	"testing"

	Ω "github.com/onsi/gomega"
)

func TestHandle(t *testing.T) {
	t.Run("when the service is unknown", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://example.com/repo-name/info/refs?service=not-a-service", nil)

		Handle(res, req, nil)

		please := Ω.NewGomegaWithT(t)

		please.Expect(res.Code).To(Ω.Equal(http.StatusNotFound))
		please.Expect(res.Body).To(Ω.ContainSubstring("not found"))
	})
}
