# Reverser

This directory contains a Service Weaver web application that reverses text.
The application has a main component and a `Reverser` component.

```mermaid
%%{init: {"flowchart": {"defaultRenderer": "elk"}} }%%
graph TD
    %% Nodes.
    github.com/brunobls/weft/Main(weft.Main)
    github.com/brunobls/weft/examples/reverser/Reverser(reverser.Reverser)

    %% Edges.
    github.com/brunobls/weft/Main --> github.com/brunobls/weft/examples/reverser/Reverser
```

## How to run?

To run this application locally in a single process, run `go run .` To run the
application locally across multiple processes, use `weaver multi deploy`.

```console
$ go run .                        # Run in a single process.
$ weaver multi deploy weft.toml # Run in multiple processes.
```

## How to interact with the application?

If running locally, open `localhost:9000` in a browser. You can also curl the
`/reverse` endpoint directly.

```console
$ curl localhost:9000/reverse?s=foo
oof
```

TODO(mwhittaker): Run on GKE.
