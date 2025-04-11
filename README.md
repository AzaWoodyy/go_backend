# Go Project Starter with Docker & MySQL

A basic Go project starter featuring:

* A simple "Hello World" HTTP server.
* Dockerized setup using `docker compose`.
* MySQL database service.
* Live reloading during development via `docker compose watch`.
* Configuration managed through a `.env` file.

## Prerequisites

* [Docker](https://docs.docker.com/get-docker/) and Docker Compose (usually included with Docker Desktop).
* [Go](https://go.dev/doc/install) (only needed if you want to run/test locally outside Docker, or for `go mod` commands).

## Setup

1. **Clone the repository (or create the files):**

    ```bash
    # If cloning:
    # git clone <your-repo-url>
    # cd go_backend

    # If creating manually, ensure all files from above are in place.
    ```

2. **Initialize Go Modules (if not cloned):**
    Run these commands in the project's root directory:

    ```bash
    go mod init github.com/yourusername/yourrepo # Or your preferred module path
    go mod tidy
    ```

3. **Configure Environment Variables:**
    * Copy the example environment file:

        ```bash
        cp .env.example .env
        ```

        *(If you didn't create a `.env.example`, just create `.env` based on the content provided above)*
    * **Edit the `.env` file:**
        * **Crucially, change `MYSQL_PASSWORD` and `MYSQL_ROOT_PASSWORD` to strong, unique passwords.**
        * Adjust `APP_PORT` and `MYSQL_PORT` if needed (e.g., if those ports are already in use on your host machine).

## Running the Application (Development with Live Reload)

1. **Build and Start:**
    The first time, or after changing `Dockerfile` or `compose.yml`, build the images:

    ```bash
    docker compose up --build
    ```

    *(You can add `-d` to run in detached mode)*

2. **Enable Live Reload:**
    To automatically rebuild and restart the `app` service when your Go source files (`*.go`, `go.mod`, `go.sum`) change, run:

    ```bash
    docker compose watch
    ```

    Keep this terminal window open. When you save changes to your Go files, Docker Compose will detect them, rebuild the `app` image, and restart the container.

    *Alternatively, some newer Docker Compose versions might support starting directly with watch:*

    ```bash
    # Try this first - might combine build, up, and watch
    docker compose up --build --watch
    ```

## Accessing the Services

* **Web Application:** Open your browser and navigate to `http://localhost:PORT`, where `PORT` is the value of `APP_PORT` in your `.env` file (default is `http://localhost:8080`). You should see "Hello World!".
* **Database:** You can connect to the MySQL database from your host machine using a database client (like DBeaver, TablePlus, MySQL Workbench) with the following details:
  * Host: `127.0.0.1` or `localhost`
  * Port: Value of `MYSQL_PORT` in `.env` (default: `3306`)
  * Database: Value of `MYSQL_DATABASE` in `.env`
  * User: Value of `MYSQL_USER` in `.env`
  * Password: Value of `MYSQL_PASSWORD` in `.env`

## Stopping the Application

1. To stop and remove the containers, network, and volumes defined in `compose.yml`, run:

    ```bash
    docker compose down
    ```

    To stop without removing the database volume (preserving data):

    ```bash
    docker compose down
    ```

    To stop *and* remove the database volume (deleting data):

    ```bash
    docker compose down -v
    ```
