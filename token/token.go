package token

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT" // add, foobar, x, y, ...
	INT    TokenType = "INT"   // 1343456
	FLOAT  TokenType = "FLOAT" // 1.23456
	STRING TokenType = "STRING"
	RUNE   TokenType = "RUNE"

	// Operators
	ASSIGN   TokenType = "="
	EQ       TokenType = "=="
	NOT_EQ   TokenType = "!="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	BANG     TokenType = "!"
	LTE      TokenType = "<="
	GTE      TokenType = ">="
	INC      TokenType = "++"
	DEC      TokenType = "--"

	LT TokenType = "<"
	GT TokenType = ">"

	SLASH TokenType = "/"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	DQUOTE    TokenType = "\""
	SQUOTE    TokenType = "'"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	RETURN   TokenType = "RETURN"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	FOR      TokenType = "FOR"
)

type TokenType string

type TokenMetadata struct {
	Line   int
	Column int
}

type Token struct {
	Type     TokenType
	Literal  string
	Metadata TokenMetadata
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}