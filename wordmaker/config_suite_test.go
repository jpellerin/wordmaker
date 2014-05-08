package wordmaker

import (
	// "fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the configured system", func() {
	It("can make a word", func() {
		input := []string{
			"A:a",
			"N:n",
			"T:t",
			"r:ANT",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(Equal("ant"))
	})

	It("can make a word including literal characters", func() {
		input := []string{
			"A:a",
			"N:n",
			"T:t",
			"r:AN-T",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(Equal("an-t"))
	})

	It("can make a word from classes with choices", func() {
		input := []string{
			"A:a/e/u",
			"N:n/m/g",
			"T:t/p/k",
			"r:ANT",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(MatchRegexp("[aeu][nmg][tpk]"))
	})

	It("can make a word from patterns with choices", func() {
		input := []string{
			"A:a/e/u",
			"N:n/m/g",
			"T:t/p/k",
			"r:A(N/NN/)T",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(MatchRegexp("[aeu]([nmg]){0,2}[tpk]"))
	})

	It("can make a word from patterns with choices and literals", func() {
		input := []string{
			"A:a/e/u",
			"N:n/m/g",
			"Q:q",
			"T:t/p/k",
			"r:A(Q/-N/-NN/)T",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(MatchRegexp("[aeu](?:q|(-[nmg]{0,2}))?[tpk]"))
	})

})
