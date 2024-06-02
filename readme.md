## PastePass

_Secure, one-time paste bin for sharing secrets._

PastePass lets you share secrets with others. The pasted content is encrypted with AES and can only be viewed once. After the paste is viewed, it is deleted from the server.

You can use this service to share passwords, API keys, or any other sensitive information that you don't want to store in your chat history or email.

PastePass is a single-binary, no dependency, fast and lightweight web service written in Go. It uses BoltDB for storing pastes after encrypting them with AES.

### How to use

Run with default options:
```bash
./pastepass
```

Open http://localhost:8080/ to access the web app. Use the `-h` option to find all options:

```bash
./pastepass -h
``` 

Usage:
```
Usage of ./pastepass:
  -app-name string
        The name of the application (e.g. ACME PastePass) (default "PastePass")
  -db-path string
        The path to the database file (default "pastes.boltdb")
  -reset-db
        Reset the database on startup
  -server-addr string
        The server address to listen on (default ":8080")
```

### Motivation

This project is inspired by [SnapPass](https://github.com/pinterest/snappass) by Pinterest. Think of it as an adaptation made for simplicity and ease of use.

It has a modern looking user interface. It is written in Go to make it easy to deploy and run on any platform. The server is a single binary with no dependencies.

See the [screenshots](./docs/screenshots.md) for a preview of the web app.

### Security

The pastes are encrypted with AES-256-GCM. The encryption key for each paste is generated randomly and stored in the database. The key is never stored on the server.

The server does not log any information about the pastes. The only information stored is the encrypted paste and its metadata (e.g. expiration time).

> [!CAUTION]
> The server does not enforce HTTPS for the endpoints, but it is absolutely necessary to use HTTPS for all requests when deploying this service in production.

### TODO

- [ ] Deployment instructions
- [ ] Release binaries
- [ ] Demo web app
- [ ] Unit tests

