package token

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT   TokenType = "IDENT" // add, foobar, x, y, ...
	INT     TokenType = "INT"   // 1343456
	FLOAT   TokenType = "FLOAT" // 1.23456
	STRING  TokenType = "STRING"
	RUNE    TokenType = "RUNE"
	COMMENT TokenType = "COMMENT"

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
	AND      TokenType = "&&"
	OR       TokenType = "||"

	LT TokenType = "<"
	GT TokenType = ">"

	SLASH TokenType = "/"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	DQUOTE    TokenType = "\""
	SQUOTE    TokenType = "'"
	COLON     TokenType = ":"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET      TokenType = "LET"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	RETURN   TokenType = "RETURN"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	FOR      TokenType = "FOR"

	// Macros
	MAGIC TokenType = "MAGIC"
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
	"magic":  MAGIC,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
