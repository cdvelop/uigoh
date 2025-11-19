package formPlatform

type permitted struct {
	Letters         bool
	Tilde           bool
	Numbers         bool
	BreakLine       bool     // line breaks allowed
	WhiteSpaces     bool     // white spaces allowed
	Tabulation      bool     // tabulation allowed
	TextNotAllowed  []string // text not allowed eg: "hola" not allowed
	Characters      []rune   // other special characters eg: '\','/','@'
	Minimum         int      // min characters eg 2 "lo" ok default 0 no defined
	Maximum         int      // max characters eg 1 "l" ok default 0 no defined}
	ExtraValidation func(string) error
	StartWith       *permitted // characters allowed at the beginning
}
