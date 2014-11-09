[![Build Status](https://drone.io/github.com/sigma/docker-fake-run/status.png)](https://drone.io/github.com/sigma/docker-fake-run/latest)

Small util to replay docker run. This includes:
- stdout
- stderr
- exit code

The entry point is a HTTP API that's compatible with the Docker /logs and
/wait calls.

Example:
$ docker-fake-run http://docker_host:4243/containers/4a6122222c22

This is useful when, for example, containers are scheduled by a separate system
(like fleet), and one still wants to integrate the actual run as part of a more
general shell script.
