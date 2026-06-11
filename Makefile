.PHONY: all clean

all:
	go run . register kahya
	go run . addfeed "Hacker News" "https://news.ycombinator.com/rss"
	go run . addfeed "Bootdev Blog" "https://www.boot.dev/blog/index.xml"

clean:
	go run . reset
