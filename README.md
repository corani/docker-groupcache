# Docker groupcache

This is a sample project for using https://github.com/mailgun/groupcache in docker-compose. It spins up three containers containing the same app joining the same group.

Usage:

```bash
$ ./build.sh -d         # builds the docker image
$ ./build.sh -up        # starts the instances
$ ./build.sh -down      # removes the instances
```

Prerequisites:

- Docker
- Docker Compose

Note:
The instances expose an HTTP endpoint on `8080`, `8081` and `8082` respectively.
