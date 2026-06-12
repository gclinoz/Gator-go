package main

import (
	"context"
	"fmt"
	"time"
	"strconv"

	"github.com/google/uuid"
	"github.com/gclinoz/Gator-go/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	_, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("can't login to an account that doesn't exist!")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("login", cmd.args[0], "success!")
	return nil
}

func handlerRegis(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	ctx := context.Background()
	userInfo := database.CreateUserParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],
	}
	_, err := s.db.CreateUser(ctx, userInfo)
	if err != nil {
		return fmt.Errorf("username already exists")
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println(cmd.args[0], "created!")
	data, err := s.db.GetUser(ctx, cmd.args[0])
	if err != nil{
		return err
	}
	fmt.Println(data)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("additional arguments will be ignored")
	}

	ctx := context.Background()
	err := s.db.DeleteAllUser(ctx)
	if err != nil{
		return err
	}
	return nil
}

func handlerListUser (s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("additional arguments will be ignored")
	}

	ctx := context.Background()
	users, err := s.db.GetAllUser(ctx)
	if err != nil {
		return err
	}

	for _, val := range users {
		if val.Name == s.cfg.Username {
			fmt.Println("*", val.Name, "(current)")
		} else {
			fmt.Println("*", val.Name)
		}
	}
	return nil
}

func handlerAgg (s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time btw req>", cmd.name)
	}

	timeBetweenReq, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("parsing time fail: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReq)

	ticker := time.NewTicker(timeBetweenReq)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
	return nil
}

func handlerAddFeed (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.name)
	}
	ctx := context.Background()

	feedInfo := database.CreateFeedParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		Name:		cmd.args[0],
		Url:		cmd.args[1],
		UserID:		user.ID,
	}

	f, err := s.db.CreateFeed(ctx, feedInfo)
	if err != nil {
		return fmt.Errorf("couldn't insert feed: %w", err)
	}
	fmt.Println(f)

	followParam := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		f.ID,
	}

	_, err = s.db.CreateFeedFollow(ctx, followParam)
	if err != nil {
		return fmt.Errorf("couldn't create new feed-follow record: %w", err)
	}

	return nil
}

func handlerListFeed (s *state, cmd command) error {
	if len(cmd.args) != 0 {
		fmt.Println("additional arguments will be ignored")
	}

	ctx := context.Background()
	data, err := s.db.GetAllFeed(ctx)
	if err != nil {
		fmt.Errorf("couldn't get feeds information: %w", err)
	}

	for _, val := range data {
		fmt.Println("Feed name:", val.Name)
		fmt.Println("URL:", val.Url)
		fmt.Println("Created by:", val.User)
		fmt.Println("---------------------")
	}
	return nil
}

func handlerAddFollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
	ctx := context.Background()

	feedInfo, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed information: %w", err)
	}

	followParam := database.CreateFeedFollowParams{
		ID:			uuid.New(),
		CreatedAt:	time.Now(),
		UpdatedAt:	time.Now(),
		UserID:		user.ID,
		FeedID:		feedInfo.ID,
	}

	f, err := s.db.CreateFeedFollow(ctx, followParam)
	if err != nil {
		return fmt.Errorf("couldn't create new feed-follow record: %w", err)
	}
	fmt.Println(f.UserName, "follows", f.FeedName)
	return nil
}

func handlerListFollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		fmt.Println("additional arguments will be ignored")
	}

	ctx := context.Background()
	data, err := s.db.GetFeedFollowForUser(ctx, user.Name)
	if err != nil {
		fmt.Errorf("couldn't get following information: %w", err)
	}

	for _, val := range data {
		fmt.Println(val.FeedName)
	}
	return nil
}

func handlerDelFollow (s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}
	ctx := context.Background()

	feedInfo, err := s.db.GetFeed(ctx, cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed information: %w", err)
	}

	delParam := database.DeleteFeedForUserParams{
		UserID:	user.ID,
		FeedID: feedInfo.ID,
	}
	err = s.db.DeleteFeedForUser(ctx, delParam)
	if err != nil {
		return fmt.Errorf("couldn't remove feed follow: %w", err)
	}
	return nil
}

func handlerPost (s *state, cmd command, user database.User) error {
	limit := 2

	if len(cmd.args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = parsedLimit
	}

	ctx := context.Background()
	getParam := database.GetPostForUserParams{
		Name:	user.Name,
		Limit:	int32(limit),
	}
	p, err := s.db.GetPostForUser(ctx, getParam)
	if err != nil {
		return fmt.Errorf("couldn't get posts: %w", err)
	}

	for _, val := range p {
		fmt.Println(val)
	}
	return nil
}
