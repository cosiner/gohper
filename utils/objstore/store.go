package objstore

import "sync"

type Object struct {
	Value string
}

type Store struct {
	indexes map[string]int
	objects []Object
	maxSize uint

	lock     sync.RWMutex
	freeList []int
}

func New(initialSize, maxSize uint) *Store {
	if initialSize > maxSize {
		initialSize = maxSize
	}

	return &Store{
		indexes:  make(map[string]int, initialSize),
		objects:  make([]Object, 0, initialSize),
		freeList: make([]int, 0, initialSize/4),
		maxSize:  maxSize,
	}
}

func (s *Store) Size() int {
	s.lock.RLock()
	size := len(s.indexes)
	s.lock.RUnlock()
	return size
}

func (s *Store) Put(key string, object Object) (o *Object) {
	s.lock.Lock()
	currSize := len(s.objects)
	if uint(currSize) == s.maxSize && s.maxSize != 0 {
		s.lock.Unlock()
		return
	}

	var index int
	if lastFree := len(s.freeList) - 1; lastFree == -1 {
		index = currSize
		s.objects = append(s.objects, object)
	} else {
		index = s.freeList[lastFree]
		s.freeList = s.freeList[:lastFree]
		s.objects[index] = object
	}
	s.indexes[key] = index
	o = &s.objects[index]
	s.lock.Unlock()
	return
}

func (s *Store) Remove(key string) {
	s.lock.Lock()
	index, has := s.indexes[key]
	if has {
		delete(s.indexes, key)
		s.freeList = append(s.freeList, index)
	}
	s.lock.Unlock()
}

func (s *Store) Get(key string) (o *Object) {
	s.lock.RLock()
	index, has := s.indexes[key]
	if has {
		o = &s.objects[index]
	}
	s.lock.RUnlock()

	return o
}
