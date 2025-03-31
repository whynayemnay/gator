# Gator - Installation & Usage Guide

## Prerequisites
Before installing Gator, ensure you have the following installed:  
- **PostgreSQL** (for database storage)  
- **Go** (for building and running the program)  

---

## Setup Configuration  
Create a configuration file at `~/.gatorconfig.json` with the following content:  

```json
{
    "db_url": "YOURDBCONNECTION?sslmode=disable",
    "current_user_name": ""
}
```

Replace `YOURDBCONNECTION` with your actual PostgreSQL connection string.

---

## Installation  
Navigate to the repository’s root folder and run:  

```bash
go install
```

This will compile and install the program.

---

## Usage  
Once installed, you can use Gator by running:  

```bash
gator <COMMAND>
```

### Available Commands  
| Command   | Description                          |
|-----------|--------------------------------------|
| `login`   | Log in to your account              |
| `register` | Register a new user                |
| `reset`   | Reset user data                     |
| `users`   | List all users                      |
| `agg`     | Fetch posts from followed feeds     |
| `addfeed` | Add an RSS feed to the database     |
| `feeds`   | List available feeds                |
| `follow`  | Follow an existing feed             |
| `following` | List the feeds you’re following  |
| `unfollow` | Unfollow a feed                    |
| `browse`  | Browse posts from feeds you follow  |
---

## How It Works  

1. **Register a user** using the `register` command.  
2. **Add an RSS feed** to the database with `addfeed`.  
   - You will automatically follow the feed after adding it.  
3. **Follow existing feeds** added by other users using `follow`.  
4. **Run `agg`** in a separate terminal to fetch posts from your followed feeds into the database.  
5. **Browse your feeds** with the `browse` command.  