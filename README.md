# Go File Server

A simple, yet powerful file server written in Go, designed to facilitate easy file uploads and serving static files over HTTP. This server supports customizable settings including port, file storage directory, IP address binding, and maximum upload size, making it a versatile choice for both development and production environments.

## Features

- **Easy to Use**: Start your server with just a few command-line arguments.
- **Customizable**: Set your server port, file storage directory, IP address, and maximum upload file size.
- **Efficient File Handling**: Serve static files and handle file uploads with ease.
- **Logging**: Full request logging with customizable log output.

## Usage

`web-api-filehandler` is a command-line application that can be run with the following arguments:

- `-port`: The port on which the server will listen for incoming requests. Default is `8080`.
- `-dir`: The directory where uploaded files will be stored. Default is `./uploads`.
- `-ip`: The IP address on which the server will listen for incoming requests. Default is `localhost`.
- `-max-upload-size`: The maximum file size that can be uploaded to the server, in bytes. Default is `10485760` (10MB).

To start the server, run the following command:

```bash
web-api-filehandler -port 8080 -dir ./uploads -ip localhost -max-upload-size 10485760
```

## Installation

To install the `web-api-filehandler` server, you can use the following `go get` command:

```bash
go get github.com/neverlless/web-api-filehandler
```

This will download the source code, compile it, and install the binary in your `$GOPATH/bin` directory.

Or you can clone the repository and build the binary yourself:

```bash
git clone https://github.com/neverlless/web-api-filehandler.git
cd web-api-filehandler
go build
```

This will create a binary named `web-api-filehandler` in the current directory.

Third option is to release the binary from the [releases page](https://github.com/neverlless/web-api-filehandler/releases).

And you can run the Docker image from the [Docker Hub](https://hub.docker.com/r/neverlless/web-api-filehandler).
