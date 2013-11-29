## Background

An experimental suite of content-addressable storage servers.

Inspired by Git, Dropbox, and Trousseau, this is my attempt to understand the
internals of these systems by writing my own simple implementation.

The generic design of each implementation will expose a REST API over the
storage service. The backend will vary depending on the implementation
language, but ultimately the content will be stored on disk.

Future enhancements:

* test suite to compare implementations
* support for versioned data
* support for de-duplication

## Implementations

The first implementation is written in Python. Versions using Ruby, Go, and
Erlang will eventually follow.

### Python

To start the service:

    python src/python/cas.py [--port ...] [--dir ...]

### Ruby

Requirements:

    gem install sinatra thin

To start the service:

    ruby src/ruby/cas.rb [--port ...]

### Go

Requirements:

    go get github.com/gorilla/mux

To start the service:

    go src/go/cas.go [-port ...]
