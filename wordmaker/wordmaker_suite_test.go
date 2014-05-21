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
		var cls string
		_, items := Lex("test2", input)
		for i := range items {
			switch i.typ {
			case itemChoice:
				result = append(result, i.val)
			case itemClass:
				cls = i.val
			}
		}
		Expect(result).To(Equal(expected))
		Expect(cls).To(Equal("C"))
	})

	It("Can lex an extended class name and choices", func() {
		input := "C1:x/y/z"
		expected := []string{"x", "y", "z"}
		result := []string{}
		var cls string
		_, items := Lex("test2", input)
		for i := range items {
			switch i.typ {
			case itemChoice:
				result = append(result, i.val)
			case itemClass:
				cls = i.val
			}
		}
		Expect(result).To(Equal(expected))
		Expect(cls).To(Equal("C1"))
	})

	It("Can lex a pattern", func() {
		input := "r:CVN(X-Y(Z))"
		expected := []string{"r", ":", "CVN", "(",
			"X-Y", "(", "Z", ")", ")"}
		result := []string{}
		_, items := Lex("test3", input)
		for i := range items {
			result = append(result, i.val)
		}
		Expect(result).To(Equal(expected))
	})

	It("Can lex a blank-first subpattern", func() {
		input := "r:CVN(/X-Y)"
		expected := []string{"r", ":", "CVN", "(",
			"", "/", "X-Y", ")"}
		result := []string{}
		_, items := Lex("test3", input)
		for i := range items {
			result = append(result, i.val)
		}
		Expect(result).To(Equal(expected))
	})

	It("Can lex a pattern with nested choices XXX", func() {
		input := "r:CVN((Z/)X-Y(Z))"
		expected := []string{"r", ":", "CVN", "(",
			"(", "Z", "/", "", ")", "X-Y", "(", "Z", ")", ")"}
		result := []string{}
		_, items := Lex("test3", input)
		for i := range items {
			result = append(result, i.val)
		}
		Expect(result).To(Equal(expected))
	})

	It("Can lex a pattern with a choice followed by a step", func() {
		input := "r:C(X/Y/)T"
		expected := []string{"r", ":", "C", "(",
			"X", "/", "Y", "/", "", ")", "T"}
		result := []string{}
		_, items := Lex("test3", input)
		for i := range items {
			result = append(result, i.val)
		}
		Expect(result).To(Equal(expected))
	})

	It("Can handle non-ascii chars and blanks", func() {
		input := "C:/å/ü/î//p/"
		expected := []string{"å", "ü", "î", "", "p", ""}
		result := []string{}
		_, items := Lex("test4", input)
		for i := range items {
			if i.typ == itemChoice {
				result = append(result, i.val)
			}
		}
		Expect(result).To(Equal(expected))
	})

	It("Can lex a pattern with a multi-char choice", func() {
		input := "r:C(X/YY/)T"
		expected := []string{"r", ":", "C", "(",
			"X", "/", "YY", "/", "", ")", "T"}
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
