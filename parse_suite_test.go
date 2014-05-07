package wordmaker

import (
	// "fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var _ = Describe("the parser", func() {

	It("can make choices", func() {
		cl := Choices([]interface{}{"hi", "yo"}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("hi|yo"))
	})

	It("can make nested choices", func() {
		cl := Choices([]interface{}{[]interface{}{"joe", "bob"}}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("joe|bob"))
	})

	It("can handle mixed choices", func() {
		cl := Choices([]interface{}{"hi", []interface{}{"joe", "bob"}}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("hi|joe|bob"))
	})

})

func TestWordmakerParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wordmaker Parser Suite")
}
