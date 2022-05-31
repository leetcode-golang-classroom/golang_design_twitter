# golang_design_twitter

Design a simplified version of Twitter where users can post tweets, follow/unfollow another user, and is able to see the `10` most recent tweets in the user's news feed.

Implement the `Twitter` class:

- `Twitter()` Initializes your twitter object.
- `void postTweet(int userId, int tweetId)` Composes a new tweet with ID `tweetId` by the user `userId`. Each call to this function will be made with a unique `tweetId`.
- `List<Integer> getNewsFeed(int userId)` Retrieves the `10` most recent tweet IDs in the user's news feed. Each item in the news feed must be posted by users who the user followed or by the user themself. Tweets must be **ordered from most recent to least recent**.
- `void follow(int followerId, int followeeId)` The user with ID `followerId` started following the user with ID `followeeId`.
- `void unfollow(int followerId, int followeeId)` The user with ID `followerId` started unfollowing the user with ID `followeeId`.

## Examples

**Example 1:**

```
Input
["Twitter", "postTweet", "getNewsFeed", "follow", "postTweet", "getNewsFeed", "unfollow", "getNewsFeed"]
[[], [1, 5], [1], [1, 2], [2, 6], [1], [1, 2], [1]]
Output
[null, null, [5], null, null, [6, 5], null, [5]]

Explanation
Twitter twitter = new Twitter();
twitter.postTweet(1, 5); // User 1 posts a new tweet (id = 5).
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 1 tweet id -> [5]. return [5]
twitter.follow(1, 2);    // User 1 follows user 2.
twitter.postTweet(2, 6); // User 2 posts a new tweet (id = 6).
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 2 tweet ids -> [6, 5]. Tweet id 6 should precede tweet id 5 because it is posted after tweet id 5.
twitter.unfollow(1, 2);  // User 1 unfollows user 2.
twitter.getNewsFeed(1);  // User 1's news feed should return a list with 1 tweet id -> [5], since user 1 is no longer following user 2.

```

**Constraints:**

- `1 <= userId, followerId, followeeId <= 500`
- `0 <= tweetId <= 10^4`
- All the tweets have **unique** IDs.
- At most $`3 * 10^4`$ calls will be made to `postTweet`, `getNewsFeed`, `follow`, and `unfollow`.

## 解析

題目要設計一個簡易版的 Twitter 功能，需要實作以下幾個 method

1. Constructor(): 用來初始化 Twitter 物件
2. void postTweet(userId, tweetId int): 用來紀錄 userId 的某一個 tweetId 留言ID
3. List<Integer> getNewsFeed(userId int): 用來找出 userId 發出或是 userId follow user 的最近 10 筆 tweetId
4. void follow(followerId, followeeId int): 用來設定 follower 所追蹤的 followee
5. void unfollow(followerId, followeeId int): 用來取消 follower 對某 followee的追蹤

關鍵點：

1. 需要找出與 userId 相關的最近 10筆 tweetId 所以需要對每次輸入的 tweetId 做輸入id 的全域計數

 2. 需要有一個結構來紀錄 follower —> followees 的對應

 3. 需要有一個結構來紀錄每個 userId —> tweetIds 的對應

for 1. 需要在 Twitter 結構加一個 count 用來紀錄目前累計到的順位

         還有需要設計一個結構 Message 用來紀錄 tweetId 以及其順位的對應

         如下

```go
type Message struct {
   Count, Id int
   Next *Message
}
```

for 2. 用來最佳化搜尋 userId 對應到 follower 可以使用 HashMap

         然而，因為 userId 對應到的 follower 不只一個

         所以最佳化的方式 就是使用 Nested HashMap

```go
// map userId -> map followee -> struct
followers := make(map[int](map[int]struct{}))
```

for 3. 最佳化紀錄每個 userId 對應到的 *Message 結構就是 HashMap

```go
// map userId -> Message List
twitters := make(map[int]*Message)
```

每次只要把最新的 tweetId 放在 Message List 的最前面

而 getNewsFeed 的實作相當於是在多個排序好的 Message Lists 取出最新的 10 個

這樣只要透過實作 MaxHeap 依序放入 10 個就可以完成這個功能

 

## 程式碼

```go
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

```

## 困難點

1. 需要理解到要找到最新的 tweet_id 需要自行設計一個 counter 來做紀錄
2. 這題除了結構上的設計，getNewsFeed 與 [23. Merge k Sorted Lists](https://www.notion.so/23-Merge-k-Sorted-Lists-7da988357350453b96b05e09a7b7d341) 作法是一樣的
3. 需要知道 golang Nested map 作法

## Solve Point

- [x]  Understand what problem need to solve
- [x]  Analysis Complexity