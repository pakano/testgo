package lfu

import "container/heap"

type LFU struct {
	maxBytes  int
	usedBytes int
	queue     *queue
	cache     map[string]*entry
}

func NewLFUCache(maxBytes int) *LFU {
	queue := make(queue, 0)
	return &LFU{
		maxBytes: maxBytes,
		queue:    &queue,
		cache:    make(map[string]*entry),
	}
}

func (l *LFU) Cap() int {
	return l.maxBytes
}

func (l *LFU) Size() int {
	return l.usedBytes
}

func (l *LFU) Len() int {
	return l.queue.Len()
}

func (l *LFU) Set(key string, val Value) {
	if en, ok := l.cache[key]; ok {
		l.usedBytes = l.usedBytes - en.value.Len() + val.Len()
		l.queue.Update(en, val, en.weight+1)
	} else {
		en := &entry{
			key:    key,
			value:  val,
			weight: 0,
		}

		heap.Push(l.queue, en)   // 插入queue 并重新排序为堆
		l.cache[key] = en        // 插入 map
		l.usedBytes += val.Len() // 更新内存占用

		// 如果超出内存长度，则删除最 '无用' 的元素，0表示无内存限制
		for l.maxBytes > 0 && l.usedBytes > l.maxBytes {
			l.removeOldest()
		}
	}
}

// 获取指定元素,访问次数加1
func (l *LFU) Get(key string) (Value, bool) {
	if en, ok := l.cache[key]; ok {
		l.queue.Update(en, en.value, en.weight+1)
		return en.value, true
	}
	return nil, false
}

// 删除指定元素（删除queue和map中的val）
func (l *LFU) Del(key string) {
	if en, ok := l.cache[key]; ok {
		heap.Remove(l.queue, en.index)
		l.removeElement(en)
	}
}

func (l *LFU) removeOldest() {
	if l.queue.Len() == 0 {
		return
	}
	v := heap.Pop(l.queue)
	l.removeElement(v.(*entry))
}

func (l *LFU) removeElement(en *entry) {
	delete(l.cache, en.key)
	l.usedBytes -= en.value.Len()
}
