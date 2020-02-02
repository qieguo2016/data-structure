package list

type LinkListNode struct {
	value int
	next  *LinkListNode
}

func (node *LinkListNode) IsEqual(target *LinkListNode) bool {
	return node.value == target.value && node.next == target.next
}

type LinkList struct {
	head   *LinkListNode
	length int
}

func CreateLinkList() *LinkList {
	return &LinkList{head: nil, length: 0}
}

// Insert 插入到头部
func (ll *LinkList) Insert(value int) *LinkListNode {
	node := LinkListNode{value: value}
	if ll.head != nil {
		node.next = ll.head
	}
	ll.head = &node
	ll.length++
	return &node
}

func (ll *LinkList) InsertAfter(value int, pos *LinkListNode) {
	node := LinkListNode{value: value}
	if pos.next != nil {
		node.next = pos.next
	}
	pos.next = &node
	ll.length++
}

func (ll *LinkList) Find(value int) *LinkListNode {
	for cur := ll.head; cur != nil; cur = cur.next {
		if value == cur.value {
			return cur
		}
	}
	return nil
}

// Delete 从头部删除
func (ll *LinkList) Delete() *LinkListNode {
	if ll.head == nil {
		return nil
	}
	node := ll.head
	ll.head = node.next
	ll.length--
	return node
}

// 根据取值删除
func (ll *LinkList) DeleteByValue(value int) *LinkListNode {
	cur := ll.head
	pre := ll.head
	for ; cur != nil; cur = cur.next {
		if value == cur.value {
			ll.length--
			if pre == ll.head {
				ll.head = nil
				return pre
			} else {
				pre.next = cur.next
				return cur
			}
		}
		pre = cur
	}
	return nil
}

// Visit 遍历高阶函数
func (ll *LinkList) Visit(fn func(node *LinkListNode)) {
	for cur := ll.head; cur != nil; cur = cur.next {
		fn(cur)
	}
}
