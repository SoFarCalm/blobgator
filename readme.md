# Blobgator

Blobgator is a Go-based command-line application for managing RSS and Atom feeds, following feeds, and browsing posts stored in a PostgreSQL database. It acts like a lightweight feed reader and ingestion tool, combining user management, feed discovery, and periodic feed scraping in one CLI.

## What the project does

The application allows a user to:

- Register and switch between users
- List existing users
- Add and view available feeds
- Follow and unfollow feeds
- Browse recent posts from the database
- Run a background aggregation loop that fetches feed content on a schedule and stores new posts

## Main features

- User management through CLI commands such as register, login, and users
- Feed management with addfeed, feeds, follow, following, and unfollow
- Post browsing with the browse command
- Scheduled feed scraping with the agg command
- Persistent configuration stored in a JSON file in the user's home directory

## Project structure

- main.go: entry point that initializes the database connection, config, and CLI command router
- handler_*.go: command handlers for users, feeds, posts, follows, and aggregation
- middleware.go: wraps handlers that require an authenticated/current user
- internal/config: loads and saves application configuration
- internal/database: generated database access code powered by sqlc
- sql/: SQL queries and schema used to define the database model

## Example commands

Run commands like:

- go run . register alice
- go run . login alice
- go run . addfeed my-feed https://example.com/feed.xml
- go run . follow https://example.com/feed.xml
- go run . browse 10
- go run . agg 1m

## Configuration

The app reads its configuration from a file named .gatorconfig.json in the home directory. The config stores the database URL and the currently selected user.

## Tech stack

- Go
- PostgreSQL
- sqlc
- lib/pq
- Google UUID
