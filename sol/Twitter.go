package sol

import "container/heap"

type Message struct {
	Count, Id int
	Next      *Message
}

type MaxHeap []*Message

func (h *MaxHeap) Len() int {
	return len(*h)
}

func (h *MaxHeap) Less(i, j int) bool {
	return (*h)[i].Count > (*h)[j].Count
}

func (h *MaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MaxHeap) Push(val interface{}) {
	*h = append(*h, val.(*Message))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type Twitter struct {
	count     int
	twitters  map[int]*Message
	followers map[int](map[int]struct{})
}

func Constructor() Twitter {
	return Twitter{
		count:     0,
		twitters:  make(map[int]*Message),
		followers: make(map[int](map[int]struct{})),
	}
}

func (this *Twitter) PostTweet(userId int, tweetId int) {
	newNode := &Message{Count: this.count, Id: tweetId}
	this.count++
	if node, exists := this.twitters[userId]; exists {
		newNode.Next = node
	}
	this.twitters[userId] = newNode
}

func (this *Twitter) GetNewsFeed(userId int) []int {
	pq := &MaxHeap{}
	self, exists := this.twitters[userId]
	if exists {
		heap.Push(pq, self)
	}
	followees, ok := this.followers[userId]
	if ok {
		for followeeId := range followees {
			if val, find := this.twitters[followeeId]; find {
				heap.Push(pq, val)
			}
		}
	}
	result := []int{}
	for pq.Len() > 0 {
		if len(result) == 10 {
			return result
		}
		top := heap.Pop(pq).(*Message)
		result = append(result, top.Id)
		if top.Next != nil && len(result) != 10 {
			heap.Push(pq, top.Next)
		}
	}
	return result
}

func (this *Twitter) Follow(followerId int, followeeId int) {
	if _, exists := this.followers[followerId]; !exists {
		this.followers[followerId] = make(map[int]struct{})
	}
	this.followers[followerId][followeeId] = struct{}{}
}

func (this *Twitter) Unfollow(followerId int, followeeId int) {
	if _, exists := this.followers[followerId][followeeId]; exists {
		delete(this.followers[followerId], followeeId)
	}
}

/**
* Your Twitter object will be instantiated and called as such:
* obj := Constructor();
* obj.PostTweet(userId,tweetId);
* param_2 := obj.GetNewsFeed(userId);
* obj.Follow(followerId,followeeId);
* obj.Unfollow(followerId,followeeId);
 */
