package cri_resmgr_extension_test

import (
	"flag"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gardener/gardener/test/framework"
)

func TestMain(m *testing.M) {
	framework.RegisterGardenerFrameworkFlags()
	flag.Parse()
	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}
