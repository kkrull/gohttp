package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGohttp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gohttp Suite")
}
