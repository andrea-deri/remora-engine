package conversion

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// AsBoolean attempts to convert a generic input value to a boolean.
// Supported types:
//   - bool: returned as-is
//   - string: supports "true", "false", "yes", "no", "y", "n" (case-insensitive)
//   - int: 0 is false, non-zero is true
//   - float64: 0.0 is false, non-zero is true
//
// Returns the converted boolean and a flag indicating whether the conversion succeeded.
func AsBoolean(value any) (bool, bool) {

	switch checkedValue := value.(type) {

	case bool:
		return checkedValue, true

	case string:
		switch strings.ToLower(checkedValue) {

		case "true", "yes", "y":
			return true, true
		case "false", "no", "n":
			return false, true
		}

	case int:
		return checkedValue != 0, true

	case float64:
		return checkedValue != 0.0, true
	}

	return false, false
}

// AsNumber attempts to convert a generic input value into a float64 number.
// Supported types include:
//   - float64, float32: converted directly
//   - signed integers: int, int32, int64
//   - unsigned integers: uint, uint32, uint64
//   - string: parsed using strconv.ParseFloat (after trimming spaces)
//
// Returns the float64 value and a flag indicating whether the conversion succeeded.
func AsNumber(value any) (float64, bool) {

	// Fast-path type switching: avoids reflection and minimizes overhead.
	switch checkedValue := value.(type) {

	case float64:
		return checkedValue, true
	case float32:
		return float64(checkedValue), true
	case int:
		return float64(checkedValue), true
	case int64:
		return float64(checkedValue), true
	case int32:
		return float64(checkedValue), true
	case uint:
		return float64(checkedValue), true
	case uint64:
		return float64(checkedValue), true
	case uint32:
		return float64(checkedValue), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(checkedValue), 64)
		return parsed, err == nil
	}

	return 0.0, false
}

// DecodeFromXML parses XML data from an io.Reader and converts it into a nested map[string]any.
//
// Rules of conversion:
//   - Each XML element becomes a key in the resulting map.
//   - Nested elements are represented as nested maps under their parent element's key.
//   - If multiple sibling elements share the same tag name, they are stored as a slice of maps.
//   - Text content inside an element is stored as a string value.
//   - Empty elements are represented as empty maps.
//
// This function uses two stacks internally:
//   - xmlNodeStack: keeps track of the current map for each nested element.
//   - xmlTagStack: keeps track of the tag name for each element in xmlNodeStack.
//
// It processes the XML token by token, handling StartElement, EndElement and CharData.
func DecodeFromXML(r io.Reader) (map[string]any, error) {

	decoder := xml.NewDecoder(r)

	root := map[string]any{}

	var xmlNodeStack []map[string]any
	var xmlTagStack []string
	xmlNodeStack = append(xmlNodeStack, root)
	xmlTagStack = append(xmlTagStack, "")

	for {

		nextToken, err := decoder.Token()
		if err != nil {
			// Found EOF, return from loop and end decode
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch xmlTag := nextToken.(type) {

		// XML start tag: include new element as child of stack's last element
		case xml.StartElement:

			// On start of an element, create a new map for its children
			xmlTagParent := xmlNodeStack[len(xmlNodeStack)-1]
			xmlTagKey := xmlTag.Name.Local

			emptyXmlTagContent := map[string]any{}

			// Check if xmlTagKey is a child of xmlTagParent taken from stack
			xmlTagAsChild, isXmlTagFoundInItsParent := xmlTagParent[xmlTagKey]
			if isXmlTagFoundInItsParent {

				switch xmlTagValueAsNewKey := xmlTagAsChild.(type) {
				case []any:
					// xmlTagValue is slice: append empty map in value of xmlTagValue key
					xmlTagParent[xmlTagKey] = append(xmlTagValueAsNewKey, emptyXmlTagContent)
				default:
					// xmlTagValue is not slice: set empty map as value of xmlTagValue key
					xmlTagParent[xmlTagKey] = []any{xmlTagValueAsNewKey, emptyXmlTagContent}
				}

			} else {
				// xmlTagKey is not child of xmlTagParent: add it as its new child
				xmlTagParent[xmlTagKey] = emptyXmlTagContent
			}

			// Push the stack with new xmlTag node: it will be evaluated on next cycle
			xmlNodeStack = append(xmlNodeStack, emptyXmlTagContent)
			xmlTagStack = append(xmlTagStack, xmlTagKey)

		// XML end tag: pop the stack from elaborated node
		case xml.EndElement:
			xmlNodeStack = xmlNodeStack[:len(xmlNodeStack)-1]
			xmlTagStack = xmlTagStack[:len(xmlTagStack)-1]

		// Include content in the elaborating node
		case xml.CharData:

			xmlTagContent := strings.TrimSpace(string(xmlTag))
			if xmlTagContent == "" {
				continue
			}

			xmlTagParent := xmlNodeStack[len(xmlNodeStack)-1]
			xmlTagKey := xmlTagStack[len(xmlTagStack)-1]

			// If the elaborating node is empty, replace it with this extracted value
			if len(xmlTagParent) == 0 {
				xmlNodeStack[len(xmlNodeStack)-2][xmlTagKey] = xmlTagContent
			}
		}
	}

	return root, nil
}

// EncodeAsXml serializes a map[string]any into an XML byte slice.
//
// The function handles three scenarios:
//  1. Empty map: returns a simple <root></root> element.
//  2. Single-key map: uses the key as the root element.
//  3. Multi-key map: wraps all key/value pairs inside a <root> element.
//
// It returns the XML content as a byte slice.
func EncodeAsXml(contentData map[string]any) []byte {

	xmlContentBuffer := &bytes.Buffer{}

	xmlContentBuffer.WriteString(xml.Header)
	contentDataLength := len(contentData)

	// Case 1: empty map -> return <root></root>
	if contentDataLength == 0 {
		xmlContentBuffer.WriteString("<root></root>")
		return xmlContentBuffer.Bytes()
	}

	// Case 2: single key -> use the key as the root element
	if contentDataLength == 1 {
		for contentKey, contentValue := range contentData {
			xmlContentBuffer.WriteString("<" + contentKey + ">")
			writeXML(xmlContentBuffer, contentValue)
			xmlContentBuffer.WriteString("</" + contentKey + ">")
		}
		return xmlContentBuffer.Bytes()
	}

	// Case 3: multiple keys -> wrap all elements inside a <root> element
	xmlContentBuffer.WriteString("<root>")
	for contentKey, contentValue := range contentData {
		xmlContentBuffer.WriteString("<" + contentKey + ">")
		writeXML(xmlContentBuffer, contentValue)
		xmlContentBuffer.WriteString("</" + contentKey + ">")
	}
	xmlContentBuffer.WriteString("</root>")

	return xmlContentBuffer.Bytes()
}

// writeXML writes a generic value to the buffer as XML.
//
// It handles the following types:
//   - map[string]any: each key becomes an XML tag, and its value is processed recursively
//   - []any: each item is wrapped in <item> tags and processed recursively
//   - primitive types: written directly as string content
func writeXML(xmlContentBuffer *bytes.Buffer, contentData any) {

	switch contentValue := contentData.(type) {

	// Handle maps: each key becomes an XML tag, value is processed recursively
	case map[string]any:
		for key, value := range contentValue {
			xmlContentBuffer.WriteString("<" + key + ">")
			writeXML(xmlContentBuffer, value)
			xmlContentBuffer.WriteString("</" + key + ">")
		}

	// Handle slices: each item is wrapped in <item> tags
	case []any:
		for _, singleItem := range contentValue {
			xmlContentBuffer.WriteString("<item>")
			writeXML(xmlContentBuffer, singleItem)
			xmlContentBuffer.WriteString("</item>")
		}

	// Handle primitive types: write directly as string
	default:
		fmt.Fprintf(xmlContentBuffer, "%v", contentValue)
	}
}
