package wordmaker

import (
	// "fmt"
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
		MakeChoices(header.val, items, 0.7)
	})

	It("can handles multichar choices from the lexer", func() {
		input := "C:x/yyy/z"
		_, items := Lex("p1", input)
		header := <-items
		cls := MakeChoices(header.val, items, 0.7)
		// fmt.Printf("cls %v", cls)
		Expect(len(cls.choices)).To(Equal(3))
	})

	It("can make a pattern from lexer input", func() {
		input := "r:C(X/XX/)T"
		_, items := Lex("p2", input)
		<-items
		cls := MakePattern(items, 0.7)
		Expect(cls.steps[0].Choose()).To(Equal("C"))
		Expect(cls.steps[1].Choose()).To(MatchRegexp("X|"))
		Expect(cls.steps[2].Choose()).To(Equal("T"))
	})

	It("makes the correct pattern for multichar choices", func() {
		input := "r:C(X/XX/)T"
		_, items := Lex("p2", input)
		<-items
		cls := MakePattern(items, 0.7)
		//fmt.Printf("cls %v", cls)
		item := cls.steps[1]
		//fmt.Printf("item %q", item)
		Expect(len(item.Items())).To(Equal(3))
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

	Describe("a pattern", func() {
		input := "r:C(N/-X/-XN/)T"
		_, items := Lex("p2", input)
		<-items
		pat := MakePattern(items, 0.7)

		It("returns a channel of strings from Run", func() {
			steps := pat.Run()
			for c := range steps {
				Expect(len(c)).To(Equal(1))
			}
		})

	})
})
