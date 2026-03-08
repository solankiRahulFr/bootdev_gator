package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	m map[string]func(*state, command) error
}

func (c *commands) register(name string, handler func(*state, command) error) {
	if c.m == nil {
		c.m = make(map[string]func(*state, command) error)
	}
	c.m[name] = handler
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.m[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

// Replace your handleLogin with this:
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("name is required")
	}
	name := cmd.args[0]
	ctx := context.Background()

	_, err := s.db.GetUser(ctx, name)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Fprintln(os.Stderr, "user does not exist")
		os.Exit(1)
	}
	if err != nil {
		return err
	}

	s.cfg.CurrentUser = name
	if err := s.cfg.SetUserName(name); err != nil {
		return err
	}

	fmt.Printf("logged in as %q\n", name)
	return nil
}

// Replace your handleRegister with this:
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("name is required")
	}
	name := cmd.args[0]
	ctx := context.Background()

	// Check if user exists
	_, err := s.db.GetUser(ctx, name)
	if err == nil {
		fmt.Fprintln(os.Stderr, "user already exists")
		os.Exit(1)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Create user
	id := uuid.New()
	now := time.Now().UTC()

	user, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})
	if err != nil {
		return err
	}

	s.cfg.CurrentUser = name
	if err := s.cfg.SetUserName(name); err != nil {
		return err
	}

	fmt.Printf("user %q created\n", name)
	log.Printf("created user: %+v\n", user)
	return nil
}

// New handler for reset command which will delete all users from the database
func handlerReset(s *state, cmd command) error {
	ctx := context.Background()
	if err := s.db.DeleteUsers(ctx); err != nil {
		return err
	}
	fmt.Println("all users have been deleted")
	return nil
}

// New handler for users command which will list all users in the database with the (current) user highlighted
func handlerGetAllUsers(s *state, cmd command) error {
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("%s (current)\n", user.Name)
		} else {
			fmt.Printf("%s\n", user.Name)
		}
	}
	return nil
}

// func handlerGetAllAggregations(s *state, cmd command) error {
// 	feed, err := fetchFeed("https://www.wagslane.dev/index.xml")
// 	if err != nil {
// 		return err
// 	}
// 	// Print the entire struct
// 	fmt.Printf("%+v\n", feed)
// 	return nil
// }

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("name and url are required")
	}
	name := cmd.args[0]
	url := cmd.args[1]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	id := uuid.New()
	now := time.Now().UTC()

	feed, err := s.db.Createfeed(ctx, database.CreatefeedParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return err
	}

	_, err = s.db.Createfeedfollow(ctx, database.CreatefeedfollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("feed %q created\n", name)
	log.Printf("created feed: %+v\n", feed)
	return nil
}

// New handler for feeds command which will list all name, url and username of the feeds
func handlerGetAllFeeds(s *state, cmd command) error {
	ctx := context.Background()
	feeds, err := s.db.GetAllfeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		username, err := s.db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Printf("name: %s, url: %s, user: %s\n", feed.Name, feed.Url, username.Name)
	}
	return nil
}

// New handler for follow command which will create a new entry in the feed_followers table and print the feed name and the username of the feed owner
func handlerfollowFeed(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("feed url is required")
	}
	feedUrl := cmd.args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	feed, err := s.db.Getfeedbyurl(ctx, feedUrl)
	if err != nil {
		return err
	}

	id := uuid.New()
	now := time.Now().UTC()

	feed_follow, err := s.db.Createfeedfollow(ctx, database.CreatefeedfollowParams{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("user %q is now following feed %q\n", feed_follow.UserName, feed_follow.FeedName)
	log.Printf("user %q is now following feed %q\n", feed_follow.UserName, feed_follow.FeedName)
	return nil
}

func handlerfollowingFeeds(s *state, cmd command) error {
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	feed_follows, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}

	for _, feed_follow := range feed_follows {
		fmt.Printf("user %q is following feed %q\n", feed_follow.UserName, feed_follow.FeedName)
	}
	return nil
}

func handlerunfollowFeed(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("feed url is required")
	}
	feedUrl := cmd.args[0]
	ctx := context.Background()

	user, err := s.db.GetUser(ctx, s.cfg.CurrentUser)
	if err != nil {
		return err
	}

	feed, err := s.db.Getfeedbyurl(ctx, feedUrl)
	if err != nil {
		return err
	}

	if err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}); err != nil {
		return err
	}
	fmt.Printf("user %q has unfollowed feed %q\n", user.Name, feed.Name)
	log.Printf("user %q has unfollowed feed %q\n", user.Name, feed.Name)
	return nil
}

func handlerGetAllAggregations(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("time between reqs is required")
	}
	time_between_reqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time duration: %v", err)
	}

	fmt.Printf("Collecting feeds in every: %v\n", time_between_reqs)

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		scrapeFeed(s)
	}
}

// function get feed with a limit arg from the post table for the current user and print the feed name, url and description
func handleGetFeedFromPosts(s *state, cmd command) error {
	//get the limit arg as an integer if not provided default to 2
	limit := 2
	if len(cmd.args) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %v", err)
		}
	}

	fmt.Printf("Getting feed from posts with limit: %v\n", limit)
	ctx := context.Background()
	user, err := s.db.GetUser(ctx, s.cfg.CurrentUser)
	if err != nil {
		return err
	}
	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
		Offset: 0,
	})
	if err != nil {
		return err
	}
	for _, post := range posts {
		fmt.Printf("name: %s, url: %s, description: %v\n", post.Title, post.Url, post.Description)
	}
	return nil
}
