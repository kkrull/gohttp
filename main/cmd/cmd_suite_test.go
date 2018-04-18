package cmd_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

func waitForStart() {
	time.Sleep(100 * time.Millisecond)
}

func scheduleShutdown(quit chan bool) {
	waitForStart()
	quit <- true
}
