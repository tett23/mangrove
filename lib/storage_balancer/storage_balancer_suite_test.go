package storage_balancer_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStorageBalancer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StorageBalancer Suite")
}
