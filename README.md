# Gator CLI

`gator` is a command-line tool designed to simplify working with PostgreSQL databases. This project is written in Go and can be run as a standalone binary for production use.

## Prerequisites

Before using the `gator` CLI, ensure the following are installed:

1. **PostgreSQL**
   - [Download and install PostgreSQL](https://www.postgresql.org/download/) if you don't have it already.
   - Ensure the `psql` command-line tool is available on your system.

2. **Go**
   - Install Go from [https://go.dev/dl/](https://go.dev/dl/).
   - Verify the installation by running:
     ```bash
     go version
     ```

## Installation

To install the `gator` CLI, use the `go install` command:

```bash
go install github.com/ThienDuc3112/gator@latest
```

This will download, compile, and place the `gator` binary in your `$GOPATH/bin` directory. Ensure that `$GOPATH/bin` is included in your system's `PATH` environment variable so you can run `gator` from anywhere.

## Configuration

Before running `gator`, set up a configuration file for connecting to your PostgreSQL database.

1. Create a `.gatorconfig.json` file in your home directory.
2. Populate it with your PostgreSQL connection details:

   ```json
   {
       "db_url": "postgres://your-username:your-password@localhost:your-port/your-database-name?sslmode=disable"
   }
   ```

3. Ensure the PostgreSQL server is running and accessible using the details you provided.

## Usage

Once the `gator` CLI is installed and configured, you can use it to interact with your database.

### Running the Program

Start the `gator` CLI with:

```bash
gator register <username>
```

This command will register your first user

## Commands

The following commands are available in the Gator CLI. Commands marked as **Logged-in Commands** require the user to be logged in.

---

### **Login**
Log in to an existing user account.

```bash
gator login <username>
```

- **Description**: Sets the current user to `<username>` if the user exists in the database.
- **Example**:
  ```bash
  gator login alice
  ```
- **Output**:
  ```
  Username alice setted
  ```

---

### **Register**
Register a new user in the database.

```bash
gator register <username>
```

- **Description**: Creates a new user with the given `<username>` and sets them as the current user.
- **Example**:
  ```bash
  gator register bob
  ```
- **Output**:
  ```
  bob have been registered
  ```

---

### **Reset**
Reset the database, removing all users, feeds, posts.

```bash
gator reset
```

- **Description**: Clears everything in the database.
- **Example**:
  ```bash
  gator reset
  ```
- **Output**:
  ```
  Reset all users
  ```

---

### **List Users**
List all users in the database.

```bash
gator users
```

- **Description**: Displays a list of all users in the database. The current logged-in user is highlighted.
- **Example**:
  ```bash
  gator users
  ```
- **Output**:
  ```
  * alice
  * bob (current)
  ```

### **Aggregate Feeds**
Start scraping feeds at regular intervals.

```bash
gator agg <interval>
```

- **Description**: Periodically scrapes least recent RSS feeds, based on the specified `<interval>` (e.g., `10m`, `1h`). Should be run on a separated terminal to update update posts in the background
- **Example**:
  ```bash
  gator agg 30m
  ```
- **Output**:
  ```
  Scrape <feed_name>, added <number> posts, duplicate <number> posts
  ```

---

### **Add Feed** (Logged-in Command)
Add a new RSS feed and follow it.

```bash
gator addfeed <name> <url>
```

- **Description**: Adds a new RSS feed with the specified `<name>` and `<url>`, and automatically follows the feed for the currently logged-in user.
- **Example**:
  ```bash
  gator add-feed "Example Feed" https://example.com/rss
  ```
- **Output**:
  ```
  Successfully added feed
  ```

---

### **List Feeds**
List all available RSS feeds.

```bash
gator feeds
```

- **Description**: Displays a list of all RSS feeds available in the database.
- **Example**:
  ```bash
  gator feeds
  ```
- **Output**:
  ```
  - id: <feed_id>
    created at: <timestamp>
    updated at: <timestamp>
    name: <feed_name>
    url: <feed_url>
    username: <feed_owner>
  ```

---

### **Follow Feed** (Logged-in Command)
Follow an existing RSS feed by URL.

```bash
gator follow <url>
```

- **Description**: Starts following an RSS feed with the specified `<url>` for the currently logged-in user.
- **Example**:
  ```bash
  gator follow https://example.com/rss
  ```
- **Output**:
  ```
  User <username> followed <feed_name> successfully!
  ```

---

### **List Following Feeds** (Logged-in Command)
List all RSS feeds that the user is currently following.

```bash
gator following
```

- **Description**: Displays all RSS feeds being followed by the current user.
- **Example**:
  ```bash
  gator following
  ```
- **Output**:
  ```
  Here are all the feeds followed by <username>:
  - id: <feed_id>
    created at: <timestamp>
    updated at: <timestamp>
    name: <feed_name>
    url: <feed_url>
    created by: <feed_creator>
  ```

---

### **Unfollow Feed** (Logged-in Command)
Stop following an RSS feed by URL.

```bash
gator unfollow <url>
```

- **Description**: Removes the specified feed from the list of followed feeds for the currently logged-in user.
- **Example**:
  ```bash
  gator unfollow https://example.com/rss
  ```
- **Output**:
  ```
  Unfollowed successfully
  ```

---

### **Browse Posts** (Logged-in Command)
View posts from followed feeds.

```bash
gator browse [limit]
```

- **Description**: Displays posts from feeds followed by the logged-in user. You can specify an optional `<limit>` to restrict the number of posts shown per feed (default is 2).
- **Example**:
  ```bash
  gator browse 5
  ```
- **Output**:
  ```
  - Title: <post_title>
    Description: <post_description>
    Publication date: <post_pub_date>
    Link: <post_link>
  ```

## Example Workflow

1. Install `gator`:
   ```bash
   go install github.com/your-username/gator@latest
   ```

2. Set up the `config.json` file with PostgreSQL connection details.

3. Run the `gator` program:
   ```bash
   gator run
   ```

4. Use commands like `add`, `list`, and `delete` to manage database records.

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Feel free to fork the repository, create a branch, and submit a pull request.