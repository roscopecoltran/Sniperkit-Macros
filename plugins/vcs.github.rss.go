package plugins

import (
    "github.com/kkdai/githubrss"
    "fmt"
)

func GithubRss() {

	g := NewGithubRss("kkdai")

	//Get latest 5 starred RSS
	ret, err := g.GetStarred(5)

	if err != nil {
		fmt.Println("Not get starred response:", err)
	}

	fmt.Println("Starred RSS:", ret)


	//Get latest 5 follower RSS
	ret, err = g.GetFollower(5)

	if err != nil {
		fmt.Println("Not get follower response:", err)
	}

	fmt.Println("Follower RSS:", ret)


	//Get latest 5 following RSS
	ret, err = g.GetFollowing(5)

	if err != nil {
		fmt.Println("Not get following response:", err)
	}

	fmt.Println("Following RSS:", ret)

}