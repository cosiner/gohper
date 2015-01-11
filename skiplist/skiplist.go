package skiplist
//
//type SkipList interface {
//	Insert(node Node)
//	Delete(node Node)
//	Index(index, offset int) interface{}
//	Length() int
//}
//
//type skipList struct {
//	head Node
//	tail Node
//}
//
//func NewSkipList() SkipList {
//	return new(skipList)
//}
//
//func (sl *skipList) Insert(rank int, data interface{}) {
//	if sl.head == nil {
//		sl.head = newNode(0, nil, nil, rank, data)
//		sl.tail = head
//	} else {
//		hd, tl := sl.head, sl.tail
//		if hd.Rank() < rank {
//			node := newNode(0, nil, nil, rank, data)
//		}
//	}
//}
//
//// node is a level node which has specified level and common data node
//// if node is non-null, level will also be
//func newNode(level int, next, nextLevel Node, rank int, data interface{}) Node {
//	return &node{
//		l: &levelNode{
//			level:     level,
//			next:      next,
//			nextLevel: nextLevel,
//		},
//		d: &dataNode{
//			rank: rank,
//			data: data,
//		},
//	}
//}
//
//func fullNewNode(level int, next, nextLevel Node, rank int, data interface{}) {
//
//}
//
//type node struct {
//	l *levelNode
//	d *dataNode
//}
//
//// levelNode is a levelNode node
//type levelNode struct {
//	level     int
//	next      Node
//	nextLevel *levelNode
//}
//
//// dataNode store Node's data
//type dataNode struct {
//	rank int
//	data interface{}
//}
//
//func (n *node) Data() interface{} {
//	return n.d.data
//}
//
//func (n *node) Rank() int {
//	return n.d.rank
//}
//
//func (n *node) Next() Node {
//	return n.l.next
//}
//
//func (n *node) NextLevel() Node {
//	if l := n.l.nextLevel; l != nil {
//		return &node{l: l, d: n.d}
//	}
//	return nil
//}
//func (n *node) Tail() Node {
//	l := n.l
//	for l.nextLevel != nil {
//		l = l.nextLevel
//	}
//	return &node{l: l, d: n.d}
//}
//
//func (n *node) SetNext(node Node) {
//	n.l.next = node
//}
//
//func (n *node) SetNextLevel(node Node) {
//	n.l.nextLevel = node
//}
//
//func (n *node) Split(fn func() []Node) []Node {
//	return fn()
//}
//
//func (n *node) Concat(nodes []Node, fn func([]Node)) {
//	fn(nodes)
//}
//
//func (n *node) Level() int {
//	return n.l.level
//}
//
//// UpLevel return a new node base on current, don't copy next node
//func (n *node) UpLevel() Node {
//	newNode := &node{
//		l: &levelNode{
//			level:     n.Level(),
//			next:      nil,
//			nextLevel: n,
//		},
//		d: n.d,
//	}
//}
//
//type Node interface {
//	UpLevel() Node
//	Level() int
//	Data() interface{}
//	Rank() int
//	Next() Node
//	NextLevel() Node
//	Tail() Node
//	SetNext(node Node)
//	SetNextLevel(node Node)
//	Split(func() []Node) []Node
//	Concat([]Node, func([]Node))
//}
