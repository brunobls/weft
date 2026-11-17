# Weft

**Write a modular monolith in Go, deploy it as microservices.**

Weft is a community-maintained continuation of [Service Weaver](https://github.com/ServiceWeaver/weaver), a programming framework for writing, deploying, and managing distributed applications in Go. The original project was archived by its maintainers; Weft picks it up to keep it alive, fix bugs, and keep evolving it.

You write your application as a single modular binary made of components, run and test it locally, and then deploy it as a set of distributed services without rewriting your code.

## Status

Weft is in its early days. It started from the last release of Service Weaver and is being renamed and modernized: the Go module is now `github.com/brunobls/weft`, and the toolchain and dependencies are being updated. Import paths and some configuration keys differ from Service Weaver, so Weft is not yet a drop-in replacement; see the migration notes once they are published.

## Getting started

Coming from Service Weaver? See [Migrating from Service Weaver](docs/migrating-from-service-weaver.md) for what changed.

Until the Weft documentation is published, the [Service Weaver docs](https://serviceweaver.dev/docs.html) still apply. The documentation source lives in the [`website/`](website/) folder of this repository.

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) before opening an issue or pull request.

## Acknowledgements

Weft is a fork of [Service Weaver](https://github.com/ServiceWeaver/weaver), created by Google and its community. Weft is not affiliated with or endorsed by Google.

## License

Weft is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.
