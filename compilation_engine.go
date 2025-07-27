package main

import (
	"bytes"
	"fmt"
)

type CompilationEngine struct {
	buffer       *bytes.Buffer
	tokenizer    *Tokenizer
	currentToken Token
}

func NewCompilationEngine(tokenizer *Tokenizer, buffer *bytes.Buffer) *CompilationEngine {
	ce := &CompilationEngine{tokenizer: tokenizer, buffer: buffer}
	tokenizer.Reset()
	ce.advance()
	return ce
}

func (ce *CompilationEngine) print(s string) error {
	_, err := ce.buffer.WriteString(s)
	return err
}
func (ce *CompilationEngine) printOpenTag(s string) error  { return ce.print(fmt.Sprintf("<%s>", s)) }
func (ce *CompilationEngine) printCloseTag(s string) error { return ce.print(fmt.Sprintf("</%s>", s)) }

func (ce *CompilationEngine) advance() (Token, error) {
	token, err := ce.tokenizer.Advance()
	if err != nil {
		return Token{}, err
	}
	ce.currentToken = token
	return token, nil
}

func (ce *CompilationEngine) process(tok TokenType, val string) error {
	ct := ce.currentToken
	if ct.tokenType != tok || (val != "" && ct.tokenValue != val) {
		return NewTokenErr(ct, "expected %s %s , got %s %s", tok, val, ct.tokenType, ct.tokenValue)
	}
	ce.advance()
	return ce.print(ct.Tag())
}

func (ce *CompilationEngine) ProcessClass() error {
	ce.printOpenTag("class")
	// print class keyword
	if err := ce.process(KEYWORD, KwCLASS); err != nil {
		return err
	}
	// print class name
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	// print {
	if err := ce.process(SYMBOL, SymLBRACE); err != nil {
		return err
	}
	// process all class variables
	for ce.currentToken.IsMulti(KEYWORD, KwSTATIC, KwFIELD) {
		if err := ce.processClassVar(); err != nil {
			return err
		}
	}
	// process all subroutines
	for ce.currentToken.IsMulti(KEYWORD, KwCONSTRUCTOR, KwFUNCTION, KwMETHOD) {
		if err := ce.processSubroutine(); err != nil {
			return err
		}
	}
	// print }
	if err := ce.process(SYMBOL, SymRBRACE); err != nil {
		return err
	}
	ce.printCloseTag("class")
	return nil
}

func (ce *CompilationEngine) processClassVar() error {
	ce.printOpenTag("classVarDec")
	// print field or static
	if err := ce.process(KEYWORD, ""); err != nil {
		return err
	}
	// print type
	if err := ce.processType(); err != nil {
		return err
	}
	// print varName
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	// process multiple varName
	for ce.currentToken.Is(SYMBOL, SymCOMMA) {
		if err := ce.process(SYMBOL, SymCOMMA); err != nil {
			return err
		}
		if err := ce.process(IDENTIFIER, ""); err != nil {
			return err
		}
	}
	if err := ce.process(SYMBOL, SymSEMICOLON); err != nil {
		return err
	}
	ce.printCloseTag("classVarDec")
	return nil
}

func (ce *CompilationEngine) processType() error {
	ct := ce.currentToken
	isType := ct.tokenType == KEYWORD &&
		(ct.tokenValue == KwINT || ct.tokenValue == KwCHAR || ct.tokenValue == KwBOOLEAN || ct.tokenValue == KwVOID)
	if isType {
		return ce.process(KEYWORD, "")
	}
	return ce.process(IDENTIFIER, "")
}

func (ce *CompilationEngine) processSubroutine() error {
	ce.printOpenTag("subroutineDec")
	// print function keyword (method, function, constructor)
	if err := ce.process(KEYWORD, ""); err != nil {
		return err
	}
	// print function type
	if err := ce.processType(); err != nil {
		return err
	}
	// print function name
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	// print (
	if err := ce.process(SYMBOL, SymLPAREN); err != nil {
		return err
	}
	// print parameterList
	if err := ce.processParameterList(); err != nil {
		return err
	}
	// print )
	if err := ce.process(SYMBOL, SymRPAREN); err != nil {
		return err
	}
	// process statements
	if err := ce.processSubroutineBody(); err != nil {
		return err
	}

	ce.printCloseTag("subroutineDec")
	return nil
}

func (ce *CompilationEngine) processParameterList() error {
	ce.printOpenTag("parameterList")
	if !ce.currentToken.Is(SYMBOL, SymRPAREN) {
		// print type
		if err := ce.processType(); err != nil {
			return err
		}
		// print varName
		if err := ce.process(IDENTIFIER, ""); err != nil {
			return err
		}
		// process multiple varName
		for ce.currentToken.Is(SYMBOL, SymCOMMA) {
			// print ,
			if err := ce.process(SYMBOL, SymCOMMA); err != nil {
				return err
			}
			// print type
			if err := ce.processType(); err != nil {
				return err
			}
			// print varName
			if err := ce.process(IDENTIFIER, ""); err != nil {
				return err
			}
		}
	}
	ce.printCloseTag("parameterList")
	return nil
}

func (ce *CompilationEngine) processSubroutineBody() error {
	ce.printOpenTag("subroutineBody")
	// print {
	if err := ce.process(SYMBOL, SymLBRACE); err != nil {
		return err
	}
	// process multiple varDec
	for ce.currentToken.Is(KEYWORD, KwVAR) {
		if err := ce.processVarDec(); err != nil {
			return err
		}
	}
	// process statements
	if err := ce.processStatements(); err != nil {
		return err
	}
	// print }
	if err := ce.process(SYMBOL, SymRBRACE); err != nil {
		return err
	}
	ce.printCloseTag("subroutineBody")
	return nil
}

func (ce *CompilationEngine) processVarDec() error {
	ce.printOpenTag("varDec")
	// print var keyword
	if err := ce.process(KEYWORD, KwVAR); err != nil {
		return err
	}
	// print type
	if err := ce.processType(); err != nil {
		return err
	}
	// print varName
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	// process multiple varName
	for ce.currentToken.Is(SYMBOL, SymCOMMA) {
		if err := ce.process(SYMBOL, SymCOMMA); err != nil {
			return err
		}
		if err := ce.process(IDENTIFIER, ""); err != nil {
			return err
		}
	}
	// print ;
	if err := ce.process(SYMBOL, SymSEMICOLON); err != nil {
		return err
	}
	ce.printCloseTag("varDec")
	return nil
}

func (ce *CompilationEngine) processStatements() error {
	ce.printOpenTag("statements")
	for ce.currentToken.IsMulti(KEYWORD, KwLET, KwDO, KwIF, KwWHILE, KwRETURN) {
		var err error
		switch ce.currentToken.tokenValue {
		case KwLET:
			err = ce.processLetStm()
		case KwDO:
			err = ce.processDoStm()
		case KwRETURN:
			err = ce.processReturnStm()
		case KwIF:
			err = ce.processIfStm()
		case KwWHILE:
			err = ce.processWhileStm()
		default:
			return NewTokenErr(ce.currentToken, "unknown statement: %s", ce.currentToken.Tag())
		}
		if err != nil {
			return err
		}
	}
	ce.printCloseTag("statements")
	return nil
}

func (ce *CompilationEngine) processLetStm() error {
	ce.printOpenTag("letStatement")
	// print let keyword
	if err := ce.process(KEYWORD, KwLET); err != nil {
		return err
	}
	// print varName
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	// print [
	if ce.currentToken.Is(SYMBOL, SymLSQBR) {
		// print [
		if err := ce.process(SYMBOL, SymLSQBR); err != nil {
			return err
		}
		// print expression
		if err := ce.processExpression(); err != nil {
			return err
		}
		// print ]
		if err := ce.process(SYMBOL, SymRSQBR); err != nil {
			return err
		}
	}
	// print =
	if err := ce.process(SYMBOL, SymEQ); err != nil {
		return err
	}
	// print expression
	if err := ce.processExpression(); err != nil {
		return err
	}
	// print ;
	if err := ce.process(SYMBOL, SymSEMICOLON); err != nil {
		return err
	}
	ce.printCloseTag("letStatement")
	return nil
}

func (ce *CompilationEngine) processDoStm() error {
	ce.printOpenTag("doStatement")
	// print do keyword
	if err := ce.process(KEYWORD, KwDO); err != nil {
		return err
	}
	// print identifier || do game.run(); / do draw();
	if err := ce.process(IDENTIFIER, ""); err != nil {
		return err
	}
	if err := ce.processSubroutineCall(); err != nil {
		return err
	}
	// print ;
	if err := ce.process(SYMBOL, SymSEMICOLON); err != nil {
		return err
	}
	ce.printCloseTag("doStatement")
	return nil
}

func (ce *CompilationEngine) processSubroutineCall() error {
	if ce.currentToken.Is(SYMBOL, SymDOT) {
		// print .
		if err := ce.process(SYMBOL, SymDOT); err != nil {
			return err
		}
		// print identifier
		if err := ce.process(IDENTIFIER, ""); err != nil {
			return err
		}
	}
	// print (
	if err := ce.process(SYMBOL, SymLPAREN); err != nil {
		return err
	}
	// process expression list
	if err := ce.processExpressionList(); err != nil {
		return err
	}
	// print )
	if err := ce.process(SYMBOL, SymRPAREN); err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) processReturnStm() error {
	ce.printOpenTag("returnStatement")
	// print return keyword
	if err := ce.process(KEYWORD, KwRETURN); err != nil {
		return err
	}
	if !ce.currentToken.Is(SYMBOL, SymSEMICOLON) {
		// print expression
		if err := ce.processExpression(); err != nil {
			return err
		}
	}
	// print ;
	if err := ce.process(SYMBOL, SymSEMICOLON); err != nil {
		return err
	}
	ce.printCloseTag("returnStatement")
	return nil
}

func (ce *CompilationEngine) processIfStm() error {
	ce.printOpenTag("ifStatement")
	// print if keyword
	if err := ce.process(KEYWORD, KwIF); err != nil {
		return err
	}
	// print (
	if err := ce.process(SYMBOL, SymLPAREN); err != nil {
		return err
	}
	// print expression
	if err := ce.processExpression(); err != nil {
		return err
	}
	// print (
	if err := ce.process(SYMBOL, SymRPAREN); err != nil {
		return err
	}
	// print {
	if err := ce.process(SYMBOL, SymLBRACE); err != nil {
		return err
	}
	// print statements
	if err := ce.processStatements(); err != nil {
		return err
	}
	// print }
	if err := ce.process(SYMBOL, SymRBRACE); err != nil {
		return err
	}
	if ce.currentToken.Is(KEYWORD, KwELSE) {
		// print else keyword
		if err := ce.process(KEYWORD, KwELSE); err != nil {
			return err
		}
		// print {
		if err := ce.process(SYMBOL, SymLBRACE); err != nil {
			return err
		}
		// print statements
		if err := ce.processStatements(); err != nil {
			return err
		}
		// print }
		if err := ce.process(SYMBOL, SymRBRACE); err != nil {
			return err
		}
	}
	ce.printCloseTag("ifStatement")
	return nil
}

func (ce *CompilationEngine) processWhileStm() error {
	ce.printOpenTag("whileStatement")
	// print while keyword
	if err := ce.process(KEYWORD, KwWHILE); err != nil {
		return err
	}
	// print (
	if err := ce.process(SYMBOL, SymLPAREN); err != nil {
		return err
	}
	// print expression
	if err := ce.processExpression(); err != nil {
		return err
	}
	// print )
	if err := ce.process(SYMBOL, SymRPAREN); err != nil {
		return err
	}
	// print {
	if err := ce.process(SYMBOL, SymLBRACE); err != nil {
		return err
	}
	// print statements
	if err := ce.processStatements(); err != nil {
		return err
	}
	// print }
	if err := ce.process(SYMBOL, SymRBRACE); err != nil {
		return err
	}
	ce.printCloseTag("whileStatement")
	return nil
}

func (ce *CompilationEngine) processExpressionList() error {
	ce.printOpenTag("expressionList")
	if !ce.currentToken.Is(SYMBOL, SymRPAREN) {
		if err := ce.processExpression(); err != nil {
			return err
		}
		for ce.currentToken.Is(SYMBOL, SymCOMMA) {
			if err := ce.process(SYMBOL, SymCOMMA); err != nil {
				return err
			}
			if err := ce.processExpression(); err != nil {
				return err
			}
		}
	}
	ce.printCloseTag("expressionList")
	return nil
}

// TODO: setup process expression
func (ce *CompilationEngine) processExpression() error {
	ce.printOpenTag("expression")

	ce.printOpenTag("term")
	if ce.currentToken.Is(KEYWORD, KwTHIS) {
		if err := ce.process(KEYWORD, ""); err != nil {
			return err
		}
	} else {
		if err := ce.process(IDENTIFIER, ""); err != nil {
			return err
		}
	}
	// print ;
	ce.printCloseTag("term")

	ce.printCloseTag("expression")
	return nil
}
