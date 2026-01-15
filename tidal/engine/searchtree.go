package engine

import (
	"maps"
	"sort"
	"strings"
)

// SearchTreeNodeType represents the type of a node in the search tree, primarily
// determined by the kind of token stored in the node.
// This is an alias for string to improve code readability and set predefined constants.
type SearchTreeNodeType string

const (

	// NODE_TOKEN_ROOT represents the enumerative value of the root node's type
	NODE_TOKEN_ROOT SearchTreeNodeType = "root"

	// NODE_TOKEN_PATH represents the enumerative value of the intermediate-path node's type
	NODE_TOKEN_PATH SearchTreeNodeType = "path"

	// NODE_TOKEN_WILDCARD represents the enumerative value of the wildcard node's type
	NODE_TOKEN_WILDCARD SearchTreeNodeType = "wildcard"

	// NODE_TOKEN_TERMINAL represents the enumerative value of the terminal leaf node's type
	NODE_TOKEN_TERMINAL SearchTreeNodeType = "terminal"
)

// SearchTree represents the tree structure used to index and efficiently retrieve Minnow references
// based on URL path, HTTP method, and optional special headers.
// It is constructed from MinnowReference objects and supports wildcard tokens for flexible matching.
type SearchTree struct {

	// Root defines the root search tree node where the search proces can starts.
	Root SearchTreeNode
}

// SearchTreeNode defines an instance of the node included in a Minnow search tree.
type SearchTreeNode struct {

	// Type represents the list of resources defined on the branch path
	Resources []Resource

	// Type represents the list of children related to this node
	Children []SearchTreeNode

	// Type represents the tipology of node regarding its content
	Type SearchTreeNodeType

	// Type represents the value of the token included in the node
	Token string
}

// Resource represents a single resource located at a specific path in the search tree,
// associated with a single Minnow profile.
type Resource struct {

	// ResourceId represents the unique identifier of the resource referenced by the branch
	ResourceId string

	// Headers represents the map of header on which the input must be compliant in order
	// to be associated as this resource.
	Headers map[string]string
}

// NewSearchTree initializes and returns an empty search tree.
func NewSearchTree() SearchTree {

	return SearchTree{
		Root: SearchTreeNode{
			Token:     "<root>",
			Type:      NODE_TOKEN_ROOT,
			Resources: nil,
			Children:  []SearchTreeNode{},
		},
	}
}

// AddSearchPath adds a new branch to the search tree for the specified URL path,
// HTTP method and optional special headers. The input path is tokenized and each
// token is used to traverse or create nodes along the branch.
//
// Existing nodes are reused whenever possible and new intermediate nodes are
// created only if they do not exist. At the end of the branch, a terminal node
// representing the HTTP method is created or reused and the Resource associated
// with the method and headers is attached.
//
// Child nodes are kept sorted by token to preserve efficient search performance
// and maintain proper wildcard handling.
func (searchTree *SearchTree) AddSearchPath(path string, method string, resourceId string, specialHeaders map[string]string) {

	currentNode := &searchTree.Root
	tokenizedPath := strings.SplitSeq(strings.Trim(path, "/"), "/")

	// Iterate through each token in the path to navigate or create nodes in the tree
	for token := range tokenizedPath {

		// Check if a child node with the current token already exists
		var existingNode *SearchTreeNode
		for index := range currentNode.Children {
			currentChild := &currentNode.Children[index]
			if currentChild.Token == token {
				existingNode = currentChild
				break
			}
		}

		// If no existing node is found, create a new node and insert it in sorted order
		if existingNode == nil {

			newChild := SearchTreeNode{
				Token:     token,
				Type:      extractTokenType(token),
				Children:  []SearchTreeNode{},
				Resources: []Resource{},
			}
			currentNode.Children = append(currentNode.Children, newChild)
			sortChildren(currentNode.Children)
			existingNode = currentNode.FindInChildren(token)
		}

		// Move to the next node in the path
		currentNode = existingNode
	}

	// Locate the terminal node corresponding to the HTTP method
	lowerCaseMethod := strings.ToLower(method)
	var terminalNode *SearchTreeNode
	for index := range currentNode.Children {
		currentChild := &currentNode.Children[index]
		if currentChild.Type == NODE_TOKEN_TERMINAL && currentChild.Token == lowerCaseMethod {
			terminalNode = currentChild
			break
		}
	}

	// If no terminal node is found, create a new node with HTTP method and
	// insert it in sorted order
	if terminalNode == nil {

		newTerminalNode := SearchTreeNode{
			Token:     lowerCaseMethod,
			Type:      NODE_TOKEN_TERMINAL,
			Children:  []SearchTreeNode{},
			Resources: []Resource{},
		}
		currentNode.Children = append(currentNode.Children, newTerminalNode)
		sortChildren(currentNode.Children)
		terminalNode = currentNode.FindInChildren(lowerCaseMethod)
	}

	// Attach the resource to the terminal node with its headers
	resource := Resource{
		ResourceId: resourceId,
		Headers:    extractHeaders(specialHeaders),
	}
	terminalNode.Resources = append(terminalNode.Resources, resource)
}

// RemoveSearchPath removes a specific branch from the search tree based on the
// given URL path, HTTP method and optional special headers. The method tokenizes
// the input path and traverses the corresponding node hierarchy to locate the
// terminal node representing the HTTP method. It then removes the Resource that
// matches the provided headers.
//
// After removing the Resource, any orphaned nodes (nodes without children and
// without resources) are recursively pruned up to the root, keeping the tree
// minimal and optimized for search performance. Child nodes remain sorted by token
// to preserve efficient lookups.
//
// Algorithm overview:
//   - If the path or terminal node does not exist, the method does nothing.
//   - If multiple Resources exist under a terminal node, only the matching one
//     is removed.
//   - Intermediate nodes that become empty after removal are pruned recursively.
func (searchTree *SearchTree) RemoveSearchPath(path string, method string, specialHeaders map[string]string) {

	currentNode := &searchTree.Root
	tokenizedPath := strings.SplitSeq(strings.Trim(path, "/"), "/")

	// Keep a stack of visited nodes for executing the cleanup process.
	// This stack is used later to recursively prune empty nodes upward.
	cleanupStack := []*SearchTreeNode{currentNode}

	// Traverse the tree by tokens extracted from the passed path
	// Traverse the tree according to the path tokens.
	for token := range tokenizedPath {

		nextNode := currentNode.FindInChildren(token)

		// Path does not exist, nothing to remove.
		if nextNode == nil {
			return
		}
		currentNode = nextNode
		cleanupStack = append(cleanupStack, currentNode)
	}

	// Locate the terminal node corresponding to the HTTP method.
	lowerCaseMethod := strings.ToLower(method)
	var terminalNode *SearchTreeNode
	for index := range currentNode.Children {
		child := &currentNode.Children[index]
		if child.Type == NODE_TOKEN_TERMINAL && child.Token == lowerCaseMethod {
			terminalNode = child
			break
		}
	}

	if terminalNode == nil {
		// Terminal node not found, nothing to remove.
		return
	}

	// Normalize headers for consistent comparison.
	targetHeaders := extractHeaders(specialHeaders)

	// Filter out resources that match the specified headers.
	filteredResources := make([]Resource, 0, len(terminalNode.Resources))
	for _, resource := range terminalNode.Resources {

		match := true
		for key, value := range resource.Headers {
			if targetHeaders[key] != value {
				match = false
				break
			}
		}

		// If is mismatching, mark the resource as not matching
		// and retain it in the filtered list.
		if !match {
			filteredResources = append(filteredResources, resource)
		}
	}

	terminalNode.Resources = filteredResources

	// If the terminal node is empty, prune any orphaned nodes up the tree.
	if len(terminalNode.Resources) == 0 && len(terminalNode.Children) == 0 {
		pruneEmptyNodes(cleanupStack)
	}
}

// FindInChildren searches for a child node matching the given token among the current node's
// children. The search uses a binary (dichotomic) search on non-wildcard children, which are
// kept sorted by token value. If a matching node is found, it is returned.
// If no exact match is found and a wildcard child exists, the wildcard node is
// returned. If neither an exact match nor a wildcard is present, the method returns nil.
func (currentNode *SearchTreeNode) FindInChildren(token string) *SearchTreeNode {

	// Initializing binary search cursors
	children := currentNode.Children
	leftCursor := 0
	rightCursor := len(children) - 1

	// If there are no children, return nil immediately
	if rightCursor < 0 {
		return nil
	}

	// Check for presence of wildcard node (always last in slice)
	isWildcardPresent := false
	if children[rightCursor].Type == NODE_TOKEN_WILDCARD {
		isWildcardPresent = true
		// Exclude wildcard from binary search
		rightCursor--
	}

	// Binary search for exact token match among non-wildcard children
	for leftCursor <= rightCursor {

		middleIndex := (leftCursor + rightCursor) / 2
		if children[middleIndex].Token == token {
			return &children[middleIndex]
		} else if token < children[middleIndex].Token {
			rightCursor = middleIndex - 1
		} else {
			leftCursor = middleIndex + 1
		}
	}

	// If no token is returned and the wildcard token is present, return the node related to it!
	// Return wildcard node if no exact match was found
	if isWildcardPresent {
		return &children[len(children)-1]
	}

	// No match found
	return nil
}

// GetResourceByMatchingHeader finds a Resource under the terminal node that matches
// the provided input headers. A Resource is considered a match if all its defined
// headers are present and equal to all headers in inputHeaders. Resources without
// headers are treated as matching by default.
//
// Returns the first matching Resource and true if found, otherwise returns an
// empty Resource and false.
func (terminalNode *SearchTreeNode) GetResourceByMatchingHeader(inputHeaders map[string]any) (Resource, bool) {

	for _, currentResource := range terminalNode.Resources {

		// Check if all headers defined in the Resource match the input headers
		areHeadersCompliant := true
		for mappedHeaderKey, mappedHeaderValue := range currentResource.Headers {

			if inputHeaders[mappedHeaderKey] != mappedHeaderValue {
				areHeadersCompliant = false
				break
			}
		}

		if areHeadersCompliant {
			return currentResource, true
		}
	}

	// No matching Resource found
	return Resource{}, false
}

// extractHeaders normalizes a map of headers by converting all keys to lowercase.
// This ensures consistent matching against input headers regardless of the original case.
func extractHeaders(specialHeaders map[string]string) map[string]string {

	headers := map[string]string{}
	for key, value := range specialHeaders {
		lowercaseKey := strings.ToLower(key)
		headers[lowercaseKey] = value
	}
	return headers
}

// extractTokenType returns the SearchTreeNodeType corresponding to the passed token.
// Tokens starting with ':' are treated as wildcards, all other tokens are considered standard path segments.
func extractTokenType(token string) SearchTreeNodeType {

	if token[0] == ':' {
		return NODE_TOKEN_WILDCARD
	}
	return NODE_TOKEN_PATH
}

// sortChildren sorts the children nodes of a SearchTreeNode in alphanumeric order by their tokens.
// Wildcard nodes (tokens equal to "*") are always placed at the end of the slice.
// The input slice is modified in place, so no new slice is returned.
func sortChildren(children []SearchTreeNode) {

	sort.SliceStable(children, func(i, j int) bool {

		// Place wildcard nodes after standard nodes
		if children[i].Token == "*" && children[j].Token != "*" {
			return false
		}
		if children[j].Token == "*" && children[i].Token != "*" {
			return true
		}

		// Sort standard tokens in lexicographical order
		return children[i].Token < children[j].Token
	})
}

// pruneEmptyNodes recursively removes empty nodes from a search tree, walking the node stack from
// leaf to root. A node is considered empty if it has no children and no resources. When an empty
// node is found, it is removed from its parent's children slice.
//
// The pruning stops once a non-empty node is encountered, ensuring that only orphaned nodes are removed.
func pruneEmptyNodes(pruneStack []*SearchTreeNode) {

	// Traverse the stack in reverse (from leaf to root)
	for i := len(pruneStack) - 1; i > 0; i-- {

		child := pruneStack[i]
		parent := pruneStack[i-1]

		// Remove the child if it has no children and no resources
		if len(child.Children) == 0 && len(child.Resources) == 0 {

			// Rebuild the parent's children slice excluding the child node
			newChildren := make([]SearchTreeNode, 0, len(parent.Children))
			for _, node := range parent.Children {
				if node.Token != child.Token {
					newChildren = append(newChildren, node)
				}
			}
			parent.Children = newChildren

		} else {
			// Stop pruning once a non-empty node is encountered
			return
		}
	}
}

// DeepCopy creates a deep copy of the SearchTree, including all nodes and resources.
// This ensures that modifications to the copy do not affect the original tree.
func (searchTree *SearchTree) DeepCopy() SearchTree {

	return SearchTree{
		Root: searchTree.Root.deepCopyNode(),
	}
}

// deepCopyNode creates a deep copy of the SearchTreeNode, including its resources and children.
// All nested nodes and resources are recursively duplicated, ensuring the copy is independent
// of the original.
func (node *SearchTreeNode) deepCopyNode() SearchTreeNode {

	// Copy the slice of Resources
	copyResources := make([]Resource, len(node.Resources))
	for i, res := range node.Resources {
		copyResources[i] = res.deepCopyResource()
	}

	// Recursively copy the slice of children
	copyChildren := make([]SearchTreeNode, len(node.Children))
	for i, child := range node.Children {
		copyChildren[i] = child.deepCopyNode()
	}

	return SearchTreeNode{
		Type:      node.Type,
		Token:     node.Token,
		Resources: copyResources,
		Children:  copyChildren,
	}
}

// deepCopyResource creates a deep copy of the Resource, including a full copy of the headers map.
// This ensures that modifications to the copy do not affect the original resource.
func (resource *Resource) deepCopyResource() Resource {

	copyHeaders := make(map[string]string, len(resource.Headers))
	maps.Copy(copyHeaders, resource.Headers)
	return Resource{
		ResourceId: resource.ResourceId,
		Headers:    copyHeaders,
	}
}
