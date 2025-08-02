# Jack Syntax Analyzer

**Nand2Tetris Project 10: Compiler I - Syntax Analysis**

A complete syntax analyzer (parser) for the Jack programming language, implemented in Go. This project is part of the Nand2Tetris course and transforms Jack source code into structured XML representations following the Jack language grammar.

## Overview

The Jack Syntax Analyzer performs two main tasks:

1. **Tokenization**: Breaks down Jack source code into atomic tokens (keywords, symbols, identifiers, constants)
2. **Parsing**: Analyzes token sequences according to Jack grammar rules and generates a parse tree

## Features

- ✅ Complete Jack language tokenization
- ✅ Full syntax analysis following Jack grammar specifications
- ✅ Parallel processing of multiple files
- ✅ Comprehensive error reporting with line numbers and context
- ✅ XML output formatting for easy visualization
- ✅ Support for single files or entire directories

## Usage

### Command Line Options

```bash
go run . -s <source> [-c <compare_file>]
```

**Parameters:**

- `-s`: Source file (.jack) or directory containing Jack files
- `-c`: Compare file (.xml) for validation (optional)

### Examples

**Process a single Jack file:**

```bash
go run . -s Main.jack
```

**Process all Jack files in a directory:**

```bash
go run . -s ./Square/
```

**Process with comparison file:**

```bash
go run . -s Main.jack -c Main.xml
```

## Input/Output

### Input Format

Jack source files with `.jack` extension containing valid Jack language code:

```jack
class Main {
    function void main() {
        var Array a;
        var int length;
        // ... more code
    }
}
```

### Output Format

The analyzer generates two XML files for each input `.jack` file:

#### 1. Tokenizer Output (`*T.xml`)

Lists all tokens sequentially:

```xml
<tokens>
<keyword> class </keyword>
<identifier> Main </identifier>
<symbol> { </symbol>
<keyword> function </keyword>
<keyword> void </keyword>
<identifier> main </identifier>
<!-- ... more tokens -->
</tokens>
```

#### 2. Parser Output (`*.xml`)

Structured parse tree following Jack grammar:

```xml
<class>
  <keyword> class </keyword>
  <identifier> Main </identifier>
  <symbol> { </symbol>
  <subroutineDec>
    <keyword> function </keyword>
    <keyword> void </keyword>
    <identifier> main </identifier>
    <symbol> ( </symbol>
    <parameterList></parameterList>
    <symbol> ) </symbol>
    <!-- ... structured parse tree -->
  </subroutineDec>
</class>
```

## Architecture

### Core Components

- **`main.go`**: Entry point, file processing, and parallel execution
- **`tokenizer.go`**: Lexical analysis - converts source code into tokens
- **`compilation_engine.go`**: Syntax analysis - builds parse tree from tokens
- **`error.go`**: Error handling with detailed context and stack traces
- **`xmlfmt.go`**: XML formatting utilities for readable output

### Supported Jack Language Elements

#### Tokens

- **Keywords**: `class`, `function`, `method`, `constructor`, `var`, `static`, `field`, `int`, `char`, `boolean`, `void`, `true`, `false`, `null`, `this`, `let`, `do`, `if`, `else`, `while`, `return`
- **Symbols**: `{`, `}`, `(`, `)`, `[`, `]`, `.`, `,`, `;`, `+`, `-`, `*`, `/`, `&`, `|`, `<`, `>`, `=`, `~`
- **Constants**: Integer constants, string constants
- **Identifiers**: Variable, method, and class names

#### Grammar Elements

- Class declarations
- Variable declarations (class variables, local variables)
- Subroutine declarations (methods, functions, constructors)
- Statements (let, if, while, do, return)
- Expressions and terms
- Parameter lists and expression lists

## Test Cases

The repository includes three test directories from the Nand2Tetris curriculum:

### `ArrayTest/`

Demonstrates array operations and basic Jack syntax:

- Variable declarations
- Array manipulation
- Simple expressions

### `ExpressionLessSquare/`

Tests control structures without complex expressions:

- Method calls
- Conditional statements
- Basic object-oriented features

### `Square/`

Complete game implementation testing:

- Full class hierarchies
- Complex expressions
- All Jack language features

## Error Handling

The analyzer provides detailed error reporting:

- **Line numbers** where errors occur
- **Context** showing the problematic line
- **Stack traces** for debugging
- **Clear error messages** describing the expected vs. actual tokens

Example error output:

```
Error in file Main.jack:15 -> expected symbol }, got identifier main
    var int x
--------------------------------
main.processClass()
compilation_engine.go:75
```

## Building and Running

### Prerequisites

- Go 1.16 or later

### Build

```bash
go build -o jack-analyzer
```

### Run

```bash
./jack-analyzer -s <source_file_or_directory>
```

## Implementation Details

### Tokenizer

- Uses regex patterns to identify token types
- Handles comments and whitespace
- Escapes XML special characters in symbols
- Maintains line number information for error reporting

### Parser

- Implements recursive descent parsing
- Follows Jack grammar specification exactly
- Generates well-formed XML output
- Provides detailed error messages with context

### Performance

- Parallel processing of multiple files
- Efficient memory usage with buffered output
- Fast regex-based tokenization

## Compliance

This implementation fully complies with the Nand2Tetris Project 10 specification and generates XML output that matches the expected format for the course's test suite.

## License

This project is part of the Nand2Tetris educational curriculum. See the course materials for licensing information.
