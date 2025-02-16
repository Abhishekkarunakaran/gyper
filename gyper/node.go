package gyper

type node struct {
	pathPoints map[string]*node
	methods    map[Method]HandleFunc
}
