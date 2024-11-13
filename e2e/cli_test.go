// e2e/cli_test.go
package e2e_test

import (
	"bytes"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI E2E Suite")
}

var _ = Describe("CLI Tool", func() {
	Context("when running in debug mode", func() {
		It("should show debug information", func() {
			cmd := exec.Command("../cli", "mfwmysql", "-h", "127.0.0.1:3306", "--debug")

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("Debug information: starting MySQL connection"))
		})
	})

	Context("when running in trace mode", func() {
		It("should show tracing information", func() {
			cmd := exec.Command("../cli", "mfwmysql", "-h", "127.0.0.1:3306", "--trace")

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("MySQL Address: 127.0.0.1:3306"))
		})
	})

	Context("when providing MySQL credentials", func() {
		It("should show connection details", func() {
			cmd := exec.Command("../cli", "mfwmysql", "-h", "127.0.0.1:3306", "-u", "root", "-p", "mypassword")

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			Expect(err).NotTo(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("MySQL Address: 127.0.0.1:3306"))
			Expect(out.String()).To(ContainSubstring("Username: root"))
		})
	})

	Context("when running an invalid command", func() {
		It("should return an error", func() {
			cmd := exec.Command("../cli", "unknowncmd")

			var out bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &out

			err := cmd.Run()
			Expect(err).To(HaveOccurred())
			Expect(out.String()).To(ContainSubstring("Error: unknown command"))
		})
	})
})
