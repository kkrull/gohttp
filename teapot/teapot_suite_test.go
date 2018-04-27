package teapot_test

import (
	"io"
	"testing"

	"github.com/kkrull/gohttp/http"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTeapot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "teapot")
}

type TeapotMock struct {
	RespondsToTarget string
	getTarget        string
}

func (mock *TeapotMock) Name() string {
	return "teapot mock"
}

func (mock *TeapotMock) RespondsTo(target string) bool {
	return mock.RespondsToTarget == target
}

func (mock *TeapotMock) Get(client io.Writer, req http.RequestMessage) {
	mock.getTarget = req.Target()
}

func (mock *TeapotMock) GetShouldHaveReceived(target string) {
	ExpectWithOffset(1, mock.getTarget).To(Equal(target))
}
