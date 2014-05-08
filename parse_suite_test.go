package wordmaker

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("the parser", func() {

	It("can make choices", func() {
		cl := Choices("A", []interface{}{"hi", "yo"}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("hi|yo"))
	})

	It("can make nested choices", func() {
		cl := Choices("A", []interface{}{[]interface{}{"joe", "bob"}}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("joe|bob"))
	})

	It("can handle mixed choices", func() {
		cl := Choices("A", []interface{}{"hi", []interface{}{"joe", "bob"}}, 0.5)
		ch, err := cl.Choose()
		Expect(err).To(BeNil())
		Expect(ch).To(MatchRegexp("hi|joe|bob"))
	})

	It("can handle input from the lexer", func() {
		input := "C:x/y/z"
		_, items := Lex("p1", input)
		header := <-items
		cls := MakeChoices(header.val, items, 0.7)
		fmt.Sprintf("cls %q", cls)
	})

	It("can make a pattern from lexer input", func() {
		input := "r:C(X/XX/)T"
		_, items := Lex("p2", input)
		<-items
		cls := MakePattern(items, 0.7)
		Expect(cls.steps[0].Item.(chooser).Choose()).To(Equal("C"))
		Expect(cls.steps[1].Item.(chooser).Choose()).To(MatchRegexp("X|"))
		Expect(cls.steps[2].Item.(chooser).Choose()).To(Equal("T"))
	})

	It("can produce a full config from input", func() {
		input := []string{"C:a/e/i",
			"X:t/k/p",
			"L:l/r/m/n",
			"r:XCL"}
		cfg, err := Parse("test", input, 0.7)
		Expect(err).To(BeNil())
		Expect(cfg.Name).To(Equal("test"))
		Expect(cfg.Class("C").Items()).To(Equal([]string{"a", "e", "i"}))
	})
})
