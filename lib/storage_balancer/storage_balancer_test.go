package storage_balancer_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tett23/mangrove/lib/storage_balancer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorageBalancer", func() {
	var ss storage_balancer.Storages
	BeforeEach(func() {
		err := os.Setenv("ENV", "test")
		Expect(err).To(BeNil())

		err = os.RemoveAll("./tmp")
		if err != nil {
			panic(err)
		}

		if err != nil {
			panic(err)
		}

		ss, err = storage_balancer.LoadStorages()
		if err != nil {
			panic(err)
		}
	})

	Describe("ファイルを保存できる", func() {
		It("ファイルの保存", func() {
			path := "foo/bar/testfile"
			s, err := ss.Write(path, []byte{'a'})
			Expect(err).To(BeNil())

			absPath := filepath.Join(s.Path, path)
			Expect(absPath).Should(BeAnExistingFile())
		})

		It("ファイルを分散させて作れる", func() {
			testfile := make([]byte, 1024, 1024)

			for i := 0; i < 100; i++ {
				path := fmt.Sprintf("foo/bar/testfile%d", i)
				s, err := ss.Write(path, testfile)
				Expect(err).To(BeNil())

				absPath := filepath.Join(s.Path, path)
				Expect(absPath).Should(BeAnExistingFile())
			}
		})
	})

	Describe("ファイルの移動ができる", func() {
		It("ファイルの保存", func() {
			path := "foo/bar/testfile"
			s, err := ss.Write(path, []byte{'a'})
			Expect(err).To(BeNil())

			absPath := filepath.Join(s.Path, path)
			Expect(absPath).Should(BeAnExistingFile())

			err = ss[1].Move(s, path)
			Expect(err).To(BeNil())

			absPath = filepath.Join(ss[1].Path, path)
			Expect(absPath).Should(BeAnExistingFile())
		})

		It("存在しないファイルは移動できない", func() {
			err := ss[1].Move(&ss[0], "hogehoge")
			Expect(err).NotTo(BeNil())
		})
	})
})
