package sol

import (
	"reflect"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	commands := []string{"Twitter", "postTweet", "getNewsFeed", "follow", "postTweet", "getNewsFeed", "unfollow", "getNewsFeed"}
	values := [][]int{{}, {1, 5}, {1}, {1, 2}, {2, 6}, {1}, {1, 2}, {1}}
	for idx := 0; idx < b.N; idx++ {
		Run(commands, values)
	}
}
func TestRun(t *testing.T) {
	type args struct {
		commands []string
		values   [][]int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "Example1",
			args: args{
				commands: []string{"Twitter", "postTweet", "getNewsFeed", "follow", "postTweet", "getNewsFeed", "unfollow", "getNewsFeed"},
				values:   [][]int{{}, {1, 5}, {1}, {1, 2}, {2, 6}, {1}, {1, 2}, {1}},
			},
			want: [][]int{{}, {}, {5}, {}, {}, {6, 5}, {}, {5}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.args.commands, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
