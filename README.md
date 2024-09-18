# gochat

Simple chat server over SSH made in Go

## Usage

### Server-side

```shell
git clone https://github.com/h0lm0/gochat.git
cd gochat
```

#### Dev environment
This environment uses the **dotenv** file `app/.dev.env`.
```bash
# To start:
./gochat.sh -e dev -u
# To stop:
./gochat.sh -e dev -d
```

#### Prod environment
This environment uses the **dotenv** file `app/.prod.env`.
```bash
# To start:
./gochat.sh -e prod -u
# To stop:
./gochat.sh -e prod -d
```

### Client-side

**WIP**
<!-- To remake -->

## Demo

**WIP**
<!-- To remake -->

## Deploying notes

From this first fix of the refacto branch. The deploy method has slightly changed. The Dockerfile was moved to the app folder within which was added a Makefile.
Instead of starting the app in dev mode, we compile the app in a binary that is then used to start the app. The Dockerfile+Makefile get the different libraries and build the application to then start it gracefully.

If for some reason some testing needs to be done without Docker, you have to go in the app directory and do `go mod vendor` first and then `make build`.

Concerning the stack of gochat, here are containers:

- Postgresql database: used to stored messages, rooms and messages. All informations stored in database will be encrypted.
- Keydb cache: fork of Redis, this container is used to store & manage sessions
- Gochat server: SSH server which can handle multiple sessions & create terminal for each
