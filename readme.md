## PastePass

_Secure, one-time paste bin for sharing secrets._

PastePass lets you share secrets with others. The pasted content is encrypted with AES and can only be viewed once. After the paste is viewed, it is deleted from the server.

You can use this service to share passwords, API keys, or any other sensitive information that you don't want to store in your chat history or email.

PastePass is a single-binary, no dependency, fast and lightweight web service written in Go. It uses BoltDB for storing pastes after encrypting them with AES.

**Check it out live:** https://pastepass.v1k45.com/

### How to use

#### Docker

Easiest way to run PastePass is using Docker. Run the following command to start the service:

```bash
docker run -p 8008:8008 -v ./:/data ghcr.io/v1k45/pastepass
```

You can also use the image in a docker-compose file:

```yaml
version: '3'
services:
  pastepass:
    image: ghcr.io/v1k45/pastepass:latest
    ports:
      - "8008:8008"
    volumes:
      - ./data/:/data
```

You can customize the options by passing them as command line arguments to docker-compose:

```yaml
version: '3'
services:
  pastepass:
    image: ghcr.io/v1k45/pastepass:latest
    command: ["pastepass", "-app-name", "MyPastePass", "-server-addr", ":8080"]
    ports:
      - "8080:8080"
    volumes:
      - ./data/:/data
```

#### Download

Download the binary from the [releases](https://github.com/v1k45/pastepass/releases/latest) page or build it from source.

#### Download binary

Here is a shortcut to download the binary for your platform:

```bash
curl -L  "https://github.com/v1k45/pastepass/releases/latest/download/pastepass-$(uname | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/; s/i[3-6]86/386/; s/aarch64/arm64/; s/armv7l/arm/')" -o pastepass
chmod +x pastepass
```

#### Build from source

```bash
git clone https://github.com/v1k45/pastepass.git
cd pastepass
make setup
make build
```

The binary will be created in the `bin/` directory.  If you are downloading the binary, make sure to make it executable:

```bash
chmod +x pastepass-linux-amd64
```

#### Run

Run with default options:

```bash
./pastepass
```

Open http://localhost:8008/ to access the web app.

Use the `-h` option to find all options:

```bash
./pastepass -h
``` 

#### Options

| Option       | Description                                                                | Default        |
|--------------|----------------------------------------------------------------------------|----------------|
| -app-name    | The name to display in the nav to ensure you are on the right environment. | PastePass      |
| -db-path     | The path of the database file                                              | pastepass.db   |
| -reset-db    | Delete all pastes before starting the server                               | `false`        |
| -server-addr | The server address to listen to                                            | `:8008`        |
| -h           | Show help message                                                          |                |


### Motivation

This project is inspired by [SnapPass](https://github.com/pinterest/snappass) by Pinterest. Think of it as an adaptation made for simplicity and ease of use.

It has a modern looking user interface. It is written in Go to make it easy to deploy and run on any platform. The server is a single binary with no dependencies.

See the [screenshots](./docs/screenshots.md) for a preview of the web app.

### Security

The pastes are encrypted with AES-256-GCM. The encryption key for each paste is generated randomly and only the encrypted text is stored in the database. The key is never stored on the server.

The server does not log any information about the pastes. The only information stored is the encrypted paste and its metadata (e.g. expiration time).

PastePass is only intended to be used as a self-hosted service, not a public paste bin.

> [!CAUTION]
> The server does not enforce HTTPS for the endpoints, but it is absolutely necessary to use HTTPS for all requests when deploying this service in production.

### TODO

- [ ] Deployment instructions
