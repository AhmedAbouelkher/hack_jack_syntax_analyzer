package main

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
)

type TokenType string

const (
	errorInt = -9999

	KEYWORD      TokenType = "keyword"
	SYMBOL       TokenType = "symbol"
	INT_CONST    TokenType = "integerConstant"
	STRING_CONST TokenType = "stringConstant"
	IDENTIFIER   TokenType = "identifier"

	KwCLASS       string = "class"
	KwCONSTRUCTOR string = "constructor"
	KwFUNCTION    string = "function"
	KwMETHOD      string = "method"
	KwFIELD       string = "field"
	KwSTATIC      string = "static"
	KwVAR         string = "var"
	KwCHAR        string = "char"
	KwBOOLEAN     string = "boolean"
	KwVOID        string = "void"
	KwTRUE        string = "true"
	KwFALSE       string = "false"
	KwNULL        string = "null"
	KwTHIS        string = "this"
	KwLET         string = "let"
	KwDO          string = "do"
	KwIF          string = "if"
	KwELSE        string = "else"
	KwWHILE       string = "while"
	KwRETURN      string = "return"
	KwINT         string = "int"

	SymLBRACE    string = "{"
	SymRBRACE    string = "}"
	SymLPAREN    string = "("
	SymRPAREN    string = ")"
	SymLSQBR     string = "["
	SymRSQBR     string = "]"
	SymDOT       string = "."
	SymCOMMA     string = ","
	SymSEMICOLON string = ";"
	SymPLUS      string = "+"
	SymMINUS     string = "-"
	SymSTAR      string = "*"
	SymSLASH     string = "/"
	SymAMPERSAND string = "&"
	SymPIPE      string = "|"
	SymLT        string = "<"
	SymGT        string = ">"
	SymEQ        string = "="
	SymTILDE     string = "~"
)

var (
	errNoMoreTokens = errors.New("no more tokens")

	keywords = []string{
		KwCLASS, KwCONSTRUCTOR, KwFUNCTION, KwMETHOD, KwFIELD, KwSTATIC,
		KwVAR, KwINT, KwCHAR, KwBOOLEAN, KwVOID, KwTRUE, KwFALSE, KwNULL,
		KwTHIS, KwLET, KwDO, KwIF, KwELSE, KwWHILE, KwRETURN,
	}
	symbols = []string{
		SymLBRACE, SymRBRACE, SymLPAREN, SymRPAREN, SymLSQBR, SymRSQBR, SymDOT, SymCOMMA, SymSEMICOLON, SymPLUS, SymMINUS, SymSTAR, SymSLASH, SymAMPERSAND, SymPIPE, SymLT, SymGT, SymEQ, SymTILDE,
	}

	keywordRgx = strings.Join(keywords, "|") // Regex pattern for matching keywords
	symRgx     = buildSymbolRegex()          // Regex pattern for matching symbols
	numRgx     = `\d+`                       // Regex pattern for matching integer constants
	strRgx     = `"[^"\n]*"`                 // Regex pattern for matching string constants (anything between quotes, no newlines)
	idRgx      = `[\w\-]+`                   // Regex pattern for matching identifiers (alphanumeric + underscore + hyphen)
)

func buildSymbolRegex() string {
	escapedSymbols := make([]string, len(symbols))
	for i, symbol := range symbols {
		escapedSymbols[i] = regexp.QuoteMeta(symbol)
	}
	return strings.Join(escapedSymbols, "|")
}

type Token struct {
	tokenType  TokenType
	tokenValue string
	lineNum    int
}

func (t Token) Tag() string {
	return fmt.Sprintf("<%s> %s </%s>", t.tokenType, t.tokenValue, t.tokenType)
}

func (t Token) Int() int {
	if t.tokenType != INT_CONST {
		return 0
	}
	i, err := strconv.Atoi(t.tokenValue)
	if err != nil {
		return errorInt
	}
	return i
}

func (t Token) Require(tok TokenType, val string) bool {
	isKeywordOrSymbol := t.tokenType == KEYWORD || t.tokenType == SYMBOL
	if val != "" {
		isKeywordOrSymbol = isKeywordOrSymbol && t.tokenValue == val
	}
	negative := t.tokenType != tok || isKeywordOrSymbol
	return !negative
}

func (t Token) IsMulti(typ TokenType, vals ...string) bool {
	for _, val := range vals {
		if t.Is(typ, val) {
			return true
		}
	}
	return false
}

func (t Token) Is(typ TokenType, val string) bool {
	return t.tokenType == typ && t.tokenValue == val
}

type TokenizerLine struct {
	rawLine string
	lNumber int
}

type Tokenizer struct {
	source            string
	tokens            []Token
	currentTokenIndex int
}

func NewTokenizer(source string) (*Tokenizer, error) {
	t := &Tokenizer{source: source, currentTokenIndex: -1}

	fileLines := strings.Split(t.source, "\n")
	for i, line := range fileLines {
		lNumber := i + 1
		trimmedLine := t.removeComments(line)
		if trimmedLine == "" {
			continue
		}
		tokens, err := t.tokenizeLine(lNumber, trimmedLine)
		if err != nil {
			return nil, &AnalyzerError{
				Err:     err,
				Line:    trimmedLine,
				LineNum: lNumber,
			}
		}
		t.tokens = append(t.tokens, tokens...)
	}

	return t, nil
}

func (t *Tokenizer) removeComments(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
		return ""
	}
	segments := strings.Split(line, "//")
	if len(segments) == 1 {
		return strings.TrimSpace(line)
	}
	return strings.TrimSpace(segments[0])
}

func (t *Tokenizer) tokenizeLine(lNumber int, line string) ([]Token, error) {
	wordRgx := []string{keywordRgx, symRgx, numRgx, strRgx, idRgx}
	wordRegex, err := regexp.Compile(strings.Join(wordRgx, "|"))
	if err != nil {
		return nil, err
	}
	matches := wordRegex.FindAllString(line, -1)
	tokens := []Token{}
	for _, match := range matches {
		tokenType, err := t.getTokenType(match)
		if err != nil {
			return nil, err
		}
		switch tokenType {
		case SYMBOL:
			match = html.EscapeString(match)
		case STRING_CONST:
			match = strings.ReplaceAll(match, "\"", "")
		}
		tokens = append(tokens, Token{
			tokenType:  tokenType,
			tokenValue: match,
			lineNum:    lNumber,
		})
	}
	return tokens, nil
}

func (t *Tokenizer) getTokenType(token string) (TokenType, error) {
	// use regex to match the token with proper anchoring
	keywordMatch, err := t.compileAndMatchRgx(token, "^("+keywordRgx+")$")
	if err != nil {
		return "", err
	}
	if keywordMatch {
		return KEYWORD, nil
	}
	symMatch, err := t.compileAndMatchRgx(token, "^("+symRgx+")$")
	if err != nil {
		return "", err
	}
	if symMatch {
		return SYMBOL, nil
	}
	numMatch, err := t.compileAndMatchRgx(token, "^"+numRgx+"$")
	if err != nil {
		return "", err
	}
	if numMatch {
		return INT_CONST, nil
	}
	strMatch, err := t.compileAndMatchRgx(token, "^"+strRgx+"$")
	if err != nil {
		return "", err
	}
	if strMatch {
		return STRING_CONST, nil
	}
	idMatch, err := t.compileAndMatchRgx(token, "^"+idRgx+"$")
	if err != nil {
		return "", err
	}
	if idMatch {
		return IDENTIFIER, nil
	}
	return "", errors.New("invalid token: " + token)
}

func (t *Tokenizer) compileAndMatchRgx(token string, rgx string) (bool, error) {
	rgxRegex, err := regexp.Compile(rgx)
	if err != nil {
		return false, err
	}
	return rgxRegex.MatchString(token), nil
}

func (t *Tokenizer) Reset() { t.currentTokenIndex = -1 }

func (t *Tokenizer) Top() Token {
	if t.currentTokenIndex == -1 {
		return t.tokens[0]
	}
	return t.tokens[t.currentTokenIndex+1]
}

func (t *Tokenizer) Advance() (Token, error) {
	t.currentTokenIndex++
	if t.currentTokenIndex >= len(t.tokens) {
		return Token{}, errNoMoreTokens
	}
	return t.tokens[t.currentTokenIndex], nil
}
