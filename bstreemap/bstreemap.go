package bstreemap

import (
	"container/list"
	"errors"
	"log"
	"runtime"
	"sync"
)

type node struct {
	key   string
	value string
	left  *node
	right *node
}

type BinarySearchTreeMap interface {
	Insert(string, string) error
	Get(string) (*string, error)
	HasKey(string) (bool, error)
	Iter() (map[string]string, error)
	LazyIter(int) <-chan map[string]string
}

type bstree struct {
	root     *node
	lock     sync.Mutex
	items    chan map[string]string
	iterSize int
}

func NewBstreeMap() BinarySearchTreeMap {
	return &bstree{
		root:     nil,
		lock:     sync.Mutex{},
		items:    nil,
		iterSize: 0,
	}
}

func (bst *bstree) Insert(key string, value string) error {
	bst.lock.Lock()
	defer bst.lock.Unlock()
	newNode := &node{key, value, nil, nil}
	var err error
	if bst.root == nil {
		bst.root = newNode
		err = nil
	} else {
		err = bst.insertNode(newNode, bst.root, key, value)
	}
	return err
}

func (bst *bstree) insertNode(nNode *node, root *node, key, value string) error {
	if key > root.key {
		if root.right == nil {
			root.right = nNode
			return nil
		}
		bst.insertNode(nNode, root.right, key, value)
	} else if key < root.key {
		if root.left == nil {
			root.left = nNode
			return nil
		}
		bst.insertNode(nNode, root.left, key, value)
	} else {

		root.value = value
	}
	return nil
}

func (bst *bstree) Get(key string) (*string, error) {
	if bst.root == nil {
		return nil, errors.New("map is not initialized")
	}
	return bst.getValue(key, bst.root)
}

func (bst *bstree) Iter() (map[string]string, error) {
	if bst.root == nil {
		return nil, errors.New("map is not initialized")
	}
	items := make(map[string]string)
	err := bst.iterTree(items, bst.root)
	log.Println("there are", len(items))
	return items, err
}

func (bst *bstree) iterTree(items map[string]string, root *node) error {
	if root == nil {
		return nil
	}
	items[root.key] = root.value
	if root.left != nil {
		bst.iterTree(items, root.left)
	}
	if root.right != nil {
		bst.iterTree(items, root.right)
	}
	return nil
}

func (bst *bstree) HasKey(key string) (bool, error) {
	if bst.root == nil {
		return false, errors.New("map is not initialized")
	}
	var err error
	if _, err = bst.getValue(key, bst.root); err == nil {
		return true, nil
	}
	return false, err
}

func (bst *bstree) getValue(key string, root *node) (*string, error) {
	if key == root.key {
		return &root.value, nil
	} else if key > root.key {
		if root.right == nil {
			return nil, errors.New("key not present")
		}
		return bst.getValue(key, root.right)
	} else {
		if root.left == nil {
			return nil, errors.New("key not present")
		}
		return bst.getValue(key, root.left)
	}
}

func (bst *bstree) LazyIter(size int) <-chan map[string]string {
	if bst.root == nil {
		return nil
	}
	bst.items = make(chan map[string]string)
	bst.iterSize = size
	tempItems := make(map[string]string)

	go bst.iterValues(bst.root, tempItems)

	return bst.items
}

func (bst *bstree) iterValues(root *node, tempItems map[string]string) {
	defer close(bst.items)

	nodes := list.New()
	nodes.PushBack(root)

	for nodes.Len() > 0 {
		elg := nodes.Front()
		nodes.Remove(elg)
		el := elg.Value.(*node)
		if len(tempItems) == bst.iterSize {
			bst.items <- tempItems
			runtime.Gosched()
			tempItems = make(map[string]string, bst.iterSize)
		}
		tempItems[el.key] = el.value
		if el.left != nil {
			nodes.PushBack(el.left)
		}
		if el.right != nil {
			nodes.PushBack(el.right)
		}
	}
	//final values
	bst.items <- tempItems
}
