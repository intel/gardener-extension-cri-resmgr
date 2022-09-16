package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

func TestCharts(t *testing.T) {
	// because we output very large charts
	format.MaxLength = 0
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRI-resource-manager extension test suite")
}
