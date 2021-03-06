package wordmaker

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"regexp"
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

	It("can make a word from patterns with nested choices", func() {
		input := []string{
			"A:a/e/u",
			"N:n/m/g",
			"T:t/p/k",
			"r:A((A/T)N/NN/)T",
		}
		cfg, err := Parse("test", input, 1)
		hit := false
		Expect(err).To(BeNil())
		Debugf("%q", cfg.patterns[0])
		for i := 0; i < 20; i++ {
			word, err := cfg.Word()
			Debugf("%v ", word)
			Expect(err).To(BeNil())
			Expect(word).To(MatchRegexp("[aeu]([aeutpk]){0,1}([nmg]){0,2}[tpk]"))
			match, err := regexp.MatchString("^[aeu][aeutpk][nmg][tpk]$", word)
			if match {
				hit = true
			}
		}
		Expect(hit).To(BeTrue())
	})

	It("can make a word with normal before nested choices", func() {
		input := []string{
			"A:a",
			"N:n",
			"T:t",
			"Q:q",
			"r:A(N(A/T)/Q)T",
		}
		cfg, err := Parse("test", input, 1)
		hit := false
		Expect(err).To(BeNil())
		Debugf("%q", cfg.patterns[0])
		for i := 0; i < 20; i++ {
			word, err := cfg.Word()
			Debugf("%v ", word)
			Expect(err).To(BeNil())
			Expect(word).To(MatchRegexp("^a((n[at])|q)t$"))
			match, err := regexp.MatchString("^an[at]t$", word)
			if match {
				hit = true
			}
		}
		Expect(hit).To(BeTrue())
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

	It("can make a word from classes that ref other classes", func() {
		input := []string{
			"A:a/e/u",
			"N:n/m/g",
			"Q:A/N",
			"r:Q",
		}
		cfg, err := Parse("test", input, 1)
		Expect(err).To(BeNil())
		word, err := cfg.Word()
		Expect(err).To(BeNil())
		Expect(word).To(MatchRegexp("[aeu]|[nmg]"))
	})

})
