package tidy

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"xml"
)

// String returns a string form of the xml.Name object.
func String(name xml.Name) string {
	if len(name.Space) > 0 {
		return fmt.Sprintf("%s:%s", name.Space, name.Local)
	}
	return name.Local
}

// indentation is the string to use as a single indentation level
var indentation = "  "

// Tidy takes an HTML string and tidys it up.
func Tidy(str string) (html string, err os.Error) {
	parser := xml.NewParser(strings.NewReader(str))
	// read str, token by token, and spit it out
	//
	// xml.Parser does most of the work for us here - we do a small 
	// bit by indenting
	indent := 0 // the current indent level
	token, err := parser.Token()
	inpre := false // true if we are in a <pre> tag
	for err == nil {
		switch token.(type) {
		case xml.StartElement:
			elem := token.(xml.StartElement)
			for i := 0; i < indent; i++ {
				html += indentation
			}
			if String(elem.Name)=="pre" {
				inpre = true
			}
			html += "<" + String(elem.Name)
			for _, attr := range elem.Attr {
				html += fmt.Sprintf(" %s=\"%s\"",
					String(attr.Name),
					attr.Value)
			}
			html += ">\n"
			indent++
		case xml.EndElement:
			elem := token.(xml.EndElement)
			if String(elem.Name)=="pre" {
				inpre = false
			}
			indent--
			for i := 0; i < indent; i++ {
				html += indentation
			}
			html += fmt.Sprintf("</%s>\n", String(elem.Name))
		case xml.CharData:
			data := token.(xml.CharData)
			str := bytes.NewBuffer(data).String()
			str = strings.Trim(str, " \r\n\t")
			if len(str) > 0 {
				if !inpre {
					for i := 0; i < indent; i++ {
						html += indentation
					}
				}
				html += str + "\n"
			}
		case xml.Comment:
			// don't print comments
		case xml.ProcInst:
			// TODO handle these somehow (server-side xslt?)
		case xml.Directive:
			// just spit it out
			directive := token.(xml.Directive)
			html += "<!" + bytes.NewBuffer(directive).String() + ">"
		default:
			// yikes! Not much to do about this...
		}
		token, err = parser.Token()
	}
	if err != os.EOF {
		fmt.Fprint(os.Stderr, err)
		// return the original string
		return str, err
	}
	err = nil
	return
}
