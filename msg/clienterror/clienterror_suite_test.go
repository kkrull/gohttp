package clienterror_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMsgClientError(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "msg/clienterror")
}
