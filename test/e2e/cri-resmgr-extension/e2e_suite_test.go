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
	// Note: flags usaage were not tested! Meant to be used with integration setup.
	framework.RegisterShootCreationFrameworkFlags() // shot name/perfix , cloudProfile, seed name, allowPrivilegedContainers ... calls Garden
	flag.Parse()
	os.Exit(m.Run())
}

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Suite")
}
