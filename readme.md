# ErniePJT Main API using Golang

## Installation

Make sure to have all required external deps. Look at Godeps config file to view them all.

**Preferred Method, Live Reloading (optional):**

Install Gin `go get github.com/codegangsta/gin`

Then run: `gin run main.go serve`

**Otherwise:**

To run using Go: `go run main.go serve`

## Configuration

Make sure to copy `.env.sample` to `.env` and update all fields.

Normal fields are saved in `.env` file.

After the first run of the server, new configuration file `config.json` will be created with default settings under `config` directory.

**IMPORTANT**

DO NOT forget to update SMTP settings and DB settings.
