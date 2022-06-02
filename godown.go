package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var level int = 0
var spaces string = "                    " // 10 levels

func _indent(lv int) string {
	index := lv * 2
	return spaces[:index]
}

func indent() string {
	return _indent(level)
}

func addIndent(msg string) string {
	log.Printf("%s%s\n", indent(), msg)
	level += 1
	return _indent(level)
}

func minusIndent() string {
	level += -1
	return _indent(level)
}

func getNodeLevel(node *blackfriday.Node) int {
	level := 0
	parent := node.Parent
	for parent != nil {
		level += 1
		parent = parent.Parent
	}
	return level
}

func isHeadingWithGo(node *blackfriday.Node) (bool, string) {
	addIndent(fmt.Sprintf("isHeadingWithGo: %s", node))
	defer minusIndent()

	if node.Type == blackfriday.Heading {
		// Use only the first word with ".go"
		firstWordNode := getFirstLeaf(node)
		if firstWordNode != nil {
			firstWord := string(firstWordNode.Literal)
			if strings.HasSuffix(firstWord, ".go") {
				log.Printf("%sgo file spec: %q", indent(), firstWord)
				return true, firstWord
			}
		}
	}
	return false, ""
}

func isHeadingWithFunc(node *blackfriday.Node) (bool, string) {
	addIndent(fmt.Sprintf("isHeadingWithFunc: %s", node))
	defer minusIndent()

	if node.Type == blackfriday.Heading {
		firstWordNode := getFirstLeaf(node)
		// log.Printf("%sfirst leaf: %q", indent(), firstWordNode.Literal)
		if firstWordNode != nil {
			firstWord := string(firstWordNode.Literal)
			if firstWord == "func " {
				// Use whole words as the func spec
				funcspec := concatenateLeafNodes(node) // concatenate the whole words
				log.Printf("%sfunction spec: %q", indent(), funcspec)
				return true, funcspec
			}
		}
	}
	return false, ""
}

func isListItemWith(node *blackfriday.Node, headerName string) bool {
	addIndent(fmt.Sprintf("isListItemWith %q: %s", headerName, node))
	defer minusIndent()

	if node.Type == blackfriday.Item {
		// Concatenate the literals of leaves to trim decorations
		log.Printf("%sfirst child: %q", indent(), node.FirstChild.Type)
		nodeLiteralOfFirstChild := string(getFirstLeaf(node.FirstChild).Literal)
		log.Printf("%slist item: %q", indent(), nodeLiteralOfFirstChild)
		if nodeLiteralOfFirstChild == headerName {
			log.Printf("%slist item found with %q: %s", indent(), headerName, nodeLiteralOfFirstChild)
			return true
		}
	}
	log.Printf("%slist item not found with %q: %s", indent(), headerName, node)
	return false
}

func concatenateLeafNodes(node *blackfriday.Node) string {
	var buff strings.Builder
	var err error

	// concatenates leaf literals
	node.Walk(leafWalker(&buff, &err))
	if err != nil {
		err = fmt.Errorf("error during concatenating leaf nodes: %s", err)
		fmt.Fprintln(os.Stderr, err)
		log.Fatal(err)
	}
	funcspec := buff.String()
	return funcspec
}

func findParentWithType(node *blackfriday.Node, nodeType blackfriday.NodeType) (bool, *blackfriday.Node) {
	// Find parent with the given node type
	var parent *blackfriday.Node = node.Parent // finding upward
	for parent != nil {
		if parent.Type == nodeType {
			return true, parent
		}
		parent = parent.Parent
	}
	return false, parent
}

func findSiblingWithType(node *blackfriday.Node, findBackword bool, nodeType blackfriday.NodeType) (bool, *blackfriday.Node) {
	// Find sibling with the given node type
	var sibling *blackfriday.Node
	if findBackword {
		sibling = node.Prev // finding backward
	} else {
		sibling = node.Next // finding forward
	}
	for sibling != nil {
		if sibling.Type == nodeType {
			return true, sibling
		}

		if findBackword {
			sibling = sibling.Prev // finding backward
		} else {
			sibling = sibling.Next // finding forward
		}
	}
	return false, sibling
}

func findSiblingWithTypeBefore(node *blackfriday.Node, nodeType blackfriday.NodeType, findBackword bool, beforeType blackfriday.NodeType) (bool, *blackfriday.Node) {
	return findSiblingWithTypeAndLiteralBefore(node, nodeType, false, "", findBackword, beforeType)
}

func findSiblingWithTypeAndLiteralBefore(node *blackfriday.Node, nodeType blackfriday.NodeType, literalCheck bool, nodeLiteral string, findBackword bool, beforeType blackfriday.NodeType) (bool, *blackfriday.Node) {
	addIndent(fmt.Sprintf("finding a sibling node with type %q and literal %q before %q from: %s", nodeType, nodeLiteral, beforeType, node))
	defer minusIndent()

	// Find sibling with the given node type before the given another node type
	var sibling *blackfriday.Node
	if findBackword {
		sibling = node.Prev // finding backward
	} else {
		sibling = node.Next // finding forward
	}

	for sibling != nil {
		log.Printf("%strying a sibling node: %s", indent(), sibling)

		if sibling.Type == beforeType {
			log.Printf("%ssibling NOT found with type %q before %q: %s", indent(), nodeType, beforeType, sibling)
			return false, sibling
		}

		if sibling.Type == nodeType {
			if literalCheck {
				siblingLiteral := concatenateLeafNodes(sibling)
				if siblingLiteral == nodeLiteral {
					log.Printf("%ssibling found with type %q and literal %q before %q: %s", indent(), nodeType, nodeLiteral, beforeType, sibling)
					return true, sibling
				}
			} else {
				log.Printf("%ssibling found with type %q before %q: %s", indent(), nodeType, beforeType, sibling)
				return true, sibling
			}
		}

		if findBackword {
			sibling = sibling.Prev // finding backward
		} else {
			sibling = sibling.Next // finding forward
		}
	}
	log.Printf("%ssibling NOT found with type %q before %q: %s", indent(), nodeType, beforeType, sibling)
	return false, sibling
}

func findSiblingHeadingWithFunc(node *blackfriday.Node) (bool, string) {
	// Find sibling of Heading with "func"
	found, heading := findSiblingWithType(node, true, blackfriday.Heading) // find backword
	if found {
		isHeadingWithColon, funcspec := isHeadingWithFunc(heading)
		if isHeadingWithColon {
			return true, funcspec
		}
	}
	// fallback. header not found or header without func is found.
	return false, ""
}

func findSiblingHeadingWithGo(node *blackfriday.Node) (bool, string) {
	// Find sibling of Heading with ".go"
	found, heading := findSiblingWithType(node, true, blackfriday.Heading) // find backword
	if found {
		isHeadingWithGo, gofilespec := isHeadingWithGo(heading)
		if isHeadingWithGo {
			return true, gofilespec
		}
	}
	// fallback. header not found or header without ".go" is found.
	return false, ""
}

func hasSiblingListItemUnderBeforeHeading(node *blackfriday.Node, findBackward bool, itemName string) bool {
	// Find sibling of Item under Item
	found, sibling := findSiblingWithTypeBefore(node, blackfriday.Item, findBackward, blackfriday.Heading) // find backword
	if found {
		isSiblingListItemUnder, _ := isListItemUnder(sibling, itemName)
		if isSiblingListItemUnder {
			return true
		}
	}
	return false
}

func hasSiblingListWithBeforeHeading(node *blackfriday.Node, findBackward bool, itemName string) bool {
	addIndent(fmt.Sprintf("finding a sibling with %q: %s", itemName, node))
	defer minusIndent()

	// Find sibling List
	found, listItem := findSiblingWithTypeAndLiteralBefore(node, blackfriday.List, false, "", findBackward, blackfriday.Heading)
	if found {
		//Find sibling of Item with itemName
		childNode := listItem.FirstChild
		if childNode != nil && childNode.Type == blackfriday.Item {
			isListItemWith := isListItemWith(childNode, itemName)
			if isListItemWith {
				return true
			}
		}
	}
	return false
}

func getFirstLeaf(node *blackfriday.Node) *blackfriday.Node {
	firstChild := node.FirstChild
	for firstChild != nil {
		if firstChild.IsLeaf() {
			return firstChild
		}
		firstChild = firstChild.FirstChild
	}
	return firstChild
}

func leafWalker(buff *strings.Builder, fmterr *error) blackfriday.NodeVisitor {
	return func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			if node.IsLeaf() {
				literal := string(node.Literal)
				fmt.Fprintf(buff, "%s", literal)
			}
		}
		return blackfriday.GoToNext
	}
}

func findParentListItemWith(node *blackfriday.Node, itemName string) (bool, *blackfriday.Node) {
	addIndent(fmt.Sprintf("findParentListItemWith %q: %s", itemName, node))
	defer minusIndent()

	// Find parent Item with itemName (parameters, returns)
	found, parentItem := findParentWithType(node, blackfriday.Item) // find upword
	if found {
		if isListItemWith(parentItem, itemName) {
			return true, parentItem
		}
	}
	return false, nil
}

func findChildWithType(node *blackfriday.Node, nodeType blackfriday.NodeType) (bool, *blackfriday.Node) {
	// Find child with the given node type
	var child *blackfriday.Node = node.FirstChild // finding downward
	for child != nil {
		if child.Type == nodeType {
			return true, child
		}
		nextChild := child.FirstChild
		if nextChild == nil {
			nextChild = child.Next
		}
		child = nextChild
	}
	return false, child
}

func isListItemUnder(node *blackfriday.Node, itemName string) (bool, string) {
	addIndent(fmt.Sprintf("isListItemUnder %q: %s", itemName, node))
	defer minusIndent()

	if node.Type == blackfriday.Item {
		// Find Item with itemName (parameters, returns) upward
		itemFound, _ := findParentListItemWith(node, itemName)
		if itemFound {
			addIndent(fmt.Sprintf("an Item under a parent Item with %q is found.", itemName))
			defer minusIndent()

			// log.Printf("%sThe item is %q -> %q -> %q -> %q", indent(), node, node.FirstChild, node.FirstChild.FirstChild, node.FirstChild.FirstChild.Next)
			log.Printf("%sfinding the first Code child..", indent())

			// Use only the first Code words as the parameter spec
			codeFound, codeNode := findChildWithType(node, blackfriday.Code)
			if codeFound {
				addIndent(fmt.Sprintf("a Code %q of an Item under an Item with %q.", codeNode.Literal, itemName))
				defer minusIndent()

				// Concatenate leaf literals to trim decorations
				paramspec := concatenateLeafNodes(codeNode)
				log.Printf("%sparameter spec: %q", indent(), paramspec)
				return true, paramspec
			}
		}
	}
	log.Printf("%sfirst Code child was not found.", indent())
	return false, ""
}

func isCodeBlockUnderHeadingWithFunc(node *blackfriday.Node) (bool, string, string) {
	addIndent(fmt.Sprintf("isCodeBlockUnderHeadingWithFunc: %s", node))
	defer minusIndent()

	if node.Type == blackfriday.CodeBlock {
		recipes := string(node.Literal)

		// Find sibling of Heading with func
		headingFound, funcspec := findSiblingHeadingWithFunc(node)
		if headingFound {

			return true, funcspec, recipes
		}
	}
	return false, "", ""
}

func isCodeBlockUnderHeadingWithGo(node *blackfriday.Node) (bool, string, string) {
	addIndent(fmt.Sprintf("isCodeBlockUnderHeadingWithGo: %s", node))
	defer minusIndent()

	if node.Type == blackfriday.CodeBlock {
		recipes := string(node.Literal)

		// Find sibling of Heading with *.go
		headingFound, gofilespec := findSiblingHeadingWithGo(node)
		if headingFound {
			return true, gofilespec, recipes
		}
	}
	return false, "", ""
}

func godownWalker(buff *strings.Builder, funcspecs *[]string, fmterr *error) blackfriday.NodeVisitor {
	return func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			level = getNodeLevel(node)
			addIndent(fmt.Sprintf("walking node at level %d: %s", level, node))
			defer minusIndent()

			switch node.Type {
			case blackfriday.CodeBlock:
				// # main.go
				// ```
				// xxxx xxxx
				// ```
				isCodeBlockUnderHeadingWithGo, _, directives := isCodeBlockUnderHeadingWithGo(node)
				if isCodeBlockUnderHeadingWithGo {
					fmt.Fprintf(buff, "%s\n\n", directives)
				}

				// ## func main
				// ```
				// xxxx xxxx
				// ```
				isCodeBlockUnderHeadingWithFunc, _, codes := isCodeBlockUnderHeadingWithFunc(node)
				if isCodeBlockUnderHeadingWithFunc {
					// open curly braces
					foundCodeBlockBeforeHeading, _ := findSiblingWithTypeBefore(node, blackfriday.CodeBlock, true, blackfriday.Heading) // find backword
					if !foundCodeBlockBeforeHeading {

						if !hasSiblingListWithBeforeHeading(node, true, "parameters") && !hasSiblingListWithBeforeHeading(node, true, "returns") {
							// there is no parameter
							log.Printf("%sThere is no parameters and returns between Heading and CodeBlock. () is appended.", indent())
							fmt.Fprintf(buff, "()")
						}

						// This code block is the first code block under heading with ".go"
						fmt.Fprintf(buff, " {\n")
					}

					// func body
					fmt.Fprintf(buff, "%s", codes)

					// close curly braces
					foundCodeBlockAfter, _ := findSiblingWithTypeBefore(node, blackfriday.CodeBlock, false, blackfriday.Heading) // find forward
					if !foundCodeBlockAfter {
						// This code block is the last code block under heading with ".go"
						fmt.Fprintf(buff, "}\n\n")
					}
				}
			case blackfriday.Heading:
				// ## func `main`
				isHeadingWithFunc, funcspec := isHeadingWithFunc(node)
				if isHeadingWithFunc {
					fmt.Fprintf(buff, "%s", funcspec)
				}
			case blackfriday.Item:
				// - parameters
				isListItemWithParameters := isListItemWith(node, "parameters")
				if isListItemWithParameters {
					fmt.Fprintf(buff, "( ")
				}

				// - `message string` : xxxxx, xxx ..
				isListItemUnderParameters, paramspec := isListItemUnder(node, "parameters")
				if isListItemUnderParameters {
					if hasSiblingListItemUnderBeforeHeading(node, true, "parameters") {
						// the second or higher
						fmt.Fprint(buff, ", ")
					}

					// parameter spec
					fmt.Fprintf(buff, "%s", paramspec)
				}

				// - returns
				isListItemWithReturns := isListItemWith(node, "returns")
				if isListItemWithReturns {
					//
					foundSiblingItem, _ := findSiblingWithTypeBefore(node, blackfriday.Item, true, blackfriday.Heading) // find backward
					if !foundSiblingItem {
						// no parameter
						fmt.Fprintf(buff, "() ")
					}

					// opening parentheses for return types. gofmt removes unnecessary parenthesis later.
					fmt.Fprintf(buff, "( ")
				}

				// - `string` : xxxx, xxx ..
				isListItemUnderReturns, paramspec := isListItemUnder(node, "returns")
				if isListItemUnderReturns {
					if hasSiblingListItemUnderBeforeHeading(node, true, "returns") { // find backward
						fmt.Fprint(buff, ", ")
					}
					fmt.Fprintf(buff, "%s ", paramspec)
				}
			}
		} else {
			if node.Type == blackfriday.Item {
				// - parameters
				isListItemWithParameters := isListItemWith(node, "parameters")
				if isListItemWithParameters {
					fmt.Fprintf(buff, ") ")
				}
				// - returns
				isListItemWithReturns := isListItemWith(node, "returns")
				if isListItemWithReturns {
					// closing parentheses for return types. gofmt removes unnecessary parenthesis later.
					fmt.Fprintf(buff, ") ")
				}
			}
		}
		return blackfriday.GoToNext
	}
}

func GenerateGoSourceFromMarkdown(input_filename string, md []byte) ([]byte, []string, error) {
	addIndent(fmt.Sprintf("generating go source file for %q.", input_filename))
	defer minusIndent()

	var err error
	n := blackfriday.New(blackfriday.WithExtensions(blackfriday.FencedCode)).Parse(md)

	var buff strings.Builder
	var funcspecs []string

	n.Walk(godownWalker(&buff, &funcspecs, &err))
	if err != nil {
		return nil, funcspecs, fmt.Errorf("%w", err)
	}

	// Print the footer
	fmt.Fprintf(&buff, "\n// This file is generated from %q by godown.\n", input_filename)
	fmt.Fprintf(&buff, "// https://github.com/hirokistring/godown\n")

	bs := bytes.NewBufferString(buff.String())
	return bs.Bytes(), funcspecs, nil
}

// This file is generated from "godown.go.md" by godown.
// https://github.com/hirokistring/godown
