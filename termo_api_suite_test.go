package main_test

import (
	. "github.com/ebenoist/termo-api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestTermoApi(t *testing.T) {
	ENV = "TEST"
	RegisterFailHandler(Fail)
	RunSpecs(t, "TermoApi Suite")
}
