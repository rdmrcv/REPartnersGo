# Solver

This repository contains code to solve the packing task with minimal oversend and minimal containers to ship.

The application right now is implemented as three parts:

1. UI — with the simple React app with in-browser React which does not require build.
2. API that right now solves the task and responds with a solution.
3. API Docs in swagger format

## Endpoints

| URL                 | Description                                             |
|---------------------|---------------------------------------------------------|
| /ui/                | UI to interact with API in accessible mode              |
| /api/solve          | API entrypoint where solver works                       |
| /swagger/index.html | Swagger UI with API description and interactive sandbox |
| /swagger/doc.json   | OpenAPI JSON doc with API definitions                   |

## UI

UI is a minimalistic React application that allows you to solve a packaging task with your own config.

It allows you to save your package sizes preset and share it with someone — click "save" and shareable URL will be in your clipboard, and when you follow that link, your package sizes preset will be loaded.

## Config

The app can be configured with env vars. List of supported config:

| ENV         | Description                                                  |
|-------------|--------------------------------------------------------------|
| LISTEN_ADDR | Set address to the listener where will be exposed API and UI |

## Build

To build the app, you can use Makefile.

```shell
$ make
```

At this point you will have ./out/server binary that can be useful for the local debug process.

To run the app, you should execute it via your cli:
```shell
$ export LISTEN_ADDR=:8080
$ ./out/server
```

### Container

If you need to build the app packed in a container, use docker (or compatible cli like podman):

```shell
$ docker build -t solver:v0.0.1 . 
```

To run the container, then you need to expose the proper port (by default, it is `:8080` but you can override it with
`LISTEN_ADDR` env).

```shell
$ docker run -p 8080:8080 solver:v0.0.1
```

## Docs

The app contains OpenAPI docs. To generate them, you could use the `make`:
```shell
$ make docs
```

Also, in the [#Endpoints](#Endpoints) you can find URLs where docs UI and docs file in the OpenAPI format will be available after the app is started.