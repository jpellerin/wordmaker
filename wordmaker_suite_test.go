package wordmaker

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var _ = Describe("lexer", func() {
	It("Can lex a class name without crashing", func() {
		input := "C:x"
		_, items := Lex("test1", input)
		for i := range items {
			fmt.Sprintf("item %q", i)
		}
	})

	It("Can lex a class name and choices", func() {
		input := "C:x/y/z"
		expected := []string{"x", "y", "z"}
		result := []string{}
		_, items := Lex("test2", input)
		for i := range items {
			if i.typ == itemChoice {
				result = append(result, i.val)
			}
		}
		Expect(result).To(Equal(expected))
	})

	It("Can lex a pattern", func() {
		input := "r:CVN(X-Y(Z))"
		expected := []string{"r", ":", "C", "V", "N", "(",
			"X", "-", "Y", "(", "Z", ")", ")"}
		result := []string{}
		_, items := Lex("test3", input)
		for i := range items {
			result = append(result, i.val)
		}
		Expect(result).To(Equal(expected))
	})

})

func TestWordmaker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wordmaker Suite")
}
