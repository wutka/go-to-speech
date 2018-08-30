package pkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ShutUp bool

func SpeakGoFile(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		speak("I can't find the file named " + speakableFilename(filename))
		fmt.Printf("File %s does not exist\n", filename)
		return
	}

	fset := token.NewFileSet() // positions are relative to fset

	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	speak("package " + f.Name.String())
	speak("imports")
	for _, imp := range f.Imports {
		symSpeech := symbolToSpeech(imp.Path.Value)
		if imp.Name != nil {
			symSpeech = symSpeech + " as " + symbolToSpeech(imp.Name.String())
		}
		speak(symSpeech)
	}

	speak("declarations")
	for _, d := range f.Decls {
		speakDecl(d)
	}

}

func speakableFilename(filename string) string {
	if strings.HasSuffix(filename, ".go") {
		filename = filename[:len(filename)-3] + " dot go"
	}
	return filename
}

var symbolTranslations = map[string]string{
	"os":      "oh ess",
	"github":  "git hub",
	"fmt":     "fumt",
	"printf":  "printf f",
	"sprintf": "s printf f",
	"fprintf": "f printf f",
	".":       "dot",
	"/":       "slash",
	"utf":     "you tee f",
	"ast":     "a s t",
	"strconv": "s t r conv",
}

func symbolToSpeech(sym string) string {
	splits := splitSymbol(sym)
	trans := translateSymbols(splits)
	return strings.Join(trans, " ")
}

func splitSymbol(symbol string) []string {
	symbols := []string{}
	currSymbol := []byte{}
	runeBuff := make([]byte, 4)
	for _, ch := range symbol {
		if unicode.IsLetter(ch) {
			n := utf8.EncodeRune(runeBuff, ch)
			currSymbol = append(currSymbol, runeBuff[:n]...)
		} else if len(currSymbol) > 0 {
			symbols = append(symbols, string(currSymbol))
			currSymbol = []byte{}
			n := utf8.EncodeRune(runeBuff, ch)
			symbols = append(symbols, string(runeBuff[:n]))
		} else {
			n := utf8.EncodeRune(runeBuff, ch)
			symbols = append(symbols, string(runeBuff[:n]))
		}

	}
	if len(currSymbol) > 0 {
		symbols = append(symbols, string(currSymbol))
	}

	return symbols
}

func translateSymbols(symbols []string) []string {
	newSyms := []string{}
	for _, sym := range symbols {
		newSym, ok := symbolTranslations[sym]
		if ok {
			sym = newSym
		}
		newSyms = append(newSyms, sym)
	}
	return newSyms
}

func speak(speech string) {
	if ShutUp {
		return
	}
	log.Printf("Saying: %s\n", speech)
	cmd := exec.Command("/usr/bin/say", speech)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Unable to run say: %+v\n", err)
		return
	}
}

func speakDecl(decl ast.Decl) {
	switch v := decl.(type) {
	case *ast.FuncDecl:
		fmt.Printf("function decl:\n%+v\n", v)
		speak("function " + symbolToSpeech(v.Name.String()))
		fmt.Printf("function name: %s\n", v.Name.String())
		speakFieldList(v.Type.Params, "taking ", "parameter")
		speakFieldList(v.Type.Results, "and returning ", "value")
		speakBlockStmt(v.Body, "function")
	case
	}
}

func speakFieldList(fields *ast.FieldList, takeOrRec string, fieldType string) {
	if fields == nil {
		speak(takeOrRec + " no " + fieldType + "s")
		return
	}
	if fields.NumFields() == 0 {
		speak(takeOrRec + " no " + fieldType + "s")
	} else if fields.NumFields() == 1 {
		speak(takeOrRec + strconv.Itoa(fields.NumFields()) + " " + fieldType)
	} else {
		speak(takeOrRec + strconv.Itoa(fields.NumFields()) + " " + fieldType + "s")
	}
	if fields.List != nil {
		for _, field := range fields.List {
			speakField(field)
		}
	}
}

func speakField(field *ast.Field) {
	as := "as "
	if len(field.Names) > 1 {
		as = "all as"
	}
	fmt.Printf("There are %d names in this field\n", len(field.Names))
	for _, fn := range field.Names {
		fmt.Printf("Field name = %s\n", fn.String())
		speak(symbolToSpeech(fn.String()))
	}
	speak(as)
	speakExpr(field.Type)
}

func speakExpr(expr ast.Expr) {
	switch v := expr.(type) {
	case *ast.Ident:
		speak(symbolToSpeech(v.String()))
	case *ast.ArrayType:
		speakArraySize(v.Len)
		if v.Len == nil {
			speak("slice of")
		} else {
			speak("array of")
		}
		speakExpr(v.Elt)
	case *ast.StarExpr:
		speak("pointer to")
		speakExpr(v.X)
	case *ast.MapType:
		speak("map with ")
		speakExpr(v.Key)
		speak("key and ")
		speakExpr(v.Value)
		speak("value")
	case *ast.SelectorExpr:
		speakExpr(v.X)
		speak("dot")
		speakExpr(v.Sel)
	case *ast.BinaryExpr:
		speakExpr(v.X)
		speakBinaryOp(v.Op.String())
		speakExpr(v.Y)
	case *ast.ParenExpr:
		speak("left paren")
		speakExpr(v.X)
		speak("right paren")
	}
}

func speakArraySize(arrSize ast.Expr) {
	/*
		switch arrSize.(type) {
		case ast.Ellipsis:
			speak("ellipsis")
		case ast.
		}
	*/
}

var binaryOpSpeech = map[string]string{
	"||": "or",
	"&&": "and",
	"==": "equals",
	"!=": "does not equal",
	"<":  "is less than",
	"<=": "is less than or equal to",
	">":  "is greater than",
	">=": "is greater than or equal to",
	"+":  "plus",
	"-":  "minus",
	"|":  "bitwise or",
	"^":  "exclusive or",
	"*":  "times",
	"/":  "divided by",
	"%":  "modulo",
	"<<": "shifted left by",
	">>": "shifted right by",
	"&":  "bitwise and",
	"&^": "bitwise and not",
}

func speakBinaryOp(op string) {
	speechVal, ok := binaryOpSpeech[op]
	if ok {
		speak(speechVal)
	}
}

var unaryOpSpeech = map[string]string{
	"+":  "positive",
	"-":  "negative",
	"!":  "not",
	"^":  "bitwise not",
	"*":  "star",
	"&":  "ref",
	"<-": "receive from channel",
}

func speakUnaryOp(op string) {
	speechVal, ok := unaryOpSpeech[op]
	if ok {
		speak(speechVal)
	}
}

func speakBlockStmt(stmts *ast.BlockStmt, bodyType string) {
	speak(bodyType + " body")
	for _, bs := range stmts.List {
		speakStmt(bs)
	}
	speak("end " + bodyType)
}
func speakStmt(stmt ast.Stmt) {
	switch v := stmt.(type) {
	case *ast.BlockStmt:
		speak("begin block")
		for _, bs := range v.List {
			speakStmt(bs)
		}
		speak("end block")
	case *ast.IfStmt:
		speak("if")
		if v.Init != nil {
			speak("with initializer ")
			speakStmt(v.Init)
			speak("when")
		}
		if v.Cond != nil {
			speakExpr(v.Cond)
		}
		if v.Body != nil {
			speakBlockStmt(v.Body, "if")
		}
		if v.Else != nil {
			switch e := v.Else.(type) {
			case *ast.BlockStmt:
				speakBlockStmt(e, "else")
			default:
				speakStmt(e)
			}
		}
	case *ast.ForStmt:
		speakForLoop(v)
	case *ast.RangeStmt:
		speak("range over ")
		speakExpr(v.X)
		speak("with")
		if v.Key != nil {
			speak("key")
			speakExpr(v.Key)
			if v.Value != nil {
				speak("and")
			}
		}
		if v.Value != nil {
			speak("value")
			speakExpr(v.Value)
		}
		if v.Body != nil {
			speakBlockStmt(v.Body, "range")
		}
	case *ast.ReturnStmt:
		speak("return")
		first := true
		for _, e := range v.Results {
			if !first {
				speak("also")
			} else {
				first = false
			}
			speakExpr(e)
		}
	}
}

func speakForLoop(fl *ast.ForStmt) {

}
