package sol

func Run(commands []string, values [][]int) [][]int {
	result := [][]int{{}}
	twitter := Constructor()
	cLen := len(commands)
	for idx := 1; idx < cLen; idx++ {
		command := commands[idx]
		switch command {
		case "postTweet":
			twitter.PostTweet(values[idx][0], values[idx][1])
			result = append(result, []int{})
		case "getNewsFeed":
			result = append(result, twitter.GetNewsFeed(values[idx][0]))
		case "follow":
			twitter.Follow(values[idx][0], values[idx][1])
			result = append(result, []int{})
		case "unfollow":
			twitter.Unfollow(values[idx][0], values[idx][1])
			result = append(result, []int{})
		}
	}
	return result
}
