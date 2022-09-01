package cri_resmgr_extension_test

import (
	"flag"
	"os"
	"testing"

	"github.com/gardener/gardener/test/framework"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMain(m *testing.M) {
	// Flags to be used against existing shoot in our dedicated infrastructure.
	framework.RegisterGardenerFrameworkFlags()
	flag.Parse()
	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}
