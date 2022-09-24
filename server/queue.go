package main

type queue struct {
	buf   [30]string
	rear  int
	front int
	size  int
}

func NewQ() *queue {
	//constructor => rear = 0 front = 0
	newQ := new(queue)
	newQ.rear = 0
	newQ.front = 0
	newQ.size = 0
	return newQ
}

func (q *queue) IsFull() bool {
	if q.size == 30 {
		return true
	} else {
		return false
	}
}

func (q *queue) IsEmpty() bool {
	if q.size == 0 {
		return true
	} else {
		return false
	}
}

func (q *queue) Push(msg string) {
	if q.IsFull() {
		return
	}

	q.buf[q.rear] = msg
	q.rear = (q.rear + 1) % 30
	q.size++
}

func (q *queue) Pop() string {
	if q.IsEmpty() {
		return ""
	}

	ret := q.buf[q.front]
	q.size--
	q.front = (q.front + 1) % 30
	return ret
}
