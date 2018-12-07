package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCfAppShowAll(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CfAppShowAll Suite")
}
