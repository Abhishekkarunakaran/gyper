package gyper
type node struct {
	path         string
	getMethod    HandleFunc
	postMethod   HandleFunc
	putMethod    HandleFunc
	patchMethod  HandleFunc
	deleteMethod HandleFunc
	children     map[string]*node
	isLeaf       bool
}

