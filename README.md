# butterclove

This is a web server that serves an XMLTV-formatted EPG (Electronic Program Guide) using data scraped from:

* [Night Flight Plus](https://nightflightplus.com) guide pages for its NFTV streams (e.g. [NFTV 1](https://nightflightplus.com/guide/nftv-1))
* [BUZZR](https://buzzrtv.com/) tv channel data for its online stream

## Using Docker Compose

1. Clone this repository and modify the `compose.yaml` to point to a local path for the `config` directory.

2. Build and deploy the container:

```
corncobble@debian:~/butterclove$ docker compose up -d
```

## Development

Use `make` to see all recipes.