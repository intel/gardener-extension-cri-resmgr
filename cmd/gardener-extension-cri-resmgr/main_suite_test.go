package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCharts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CRI-resource-manager extension test suite")
}
