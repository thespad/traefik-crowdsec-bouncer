# [thespad/traefik-crowdsec-bouncer](https://github.com/thespad/traefik-crowdsec-bouncer)

A http service to verify requests and bounce them according to decisions made by CrowdSec. Fork of [https://github.com/fbonalair/traefik-crowdsec-bouncer](https://github.com/fbonalair/traefik-crowdsec-bouncer)

[![GitHub Release](https://img.shields.io/github/release/thespad/traefik-crowdsec-bouncer.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github&include_prereleases)](https://github.com/thespad/traefik-crowdsec-bouncer/releases)
![Commits](https://img.shields.io/github/commits-since/thespad/traefik-crowdsec-bouncer/latest?color=26689A&include_prereleases&logo=github&style=for-the-badge)
![Image Size](https://img.shields.io/docker/image-size/thespad/whisparr/latest?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=Size)
[![Docker Pulls](https://img.shields.io/docker/pulls/thespad/whisparr.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=pulls&logo=docker)](https://hub.docker.com/r/thespad/whisparr)
[![GitHub Stars](https://img.shields.io/github/stars/thespad/traefik-crowdsec-bouncer.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github)](https://github.com/thespad/traefik-crowdsec-bouncer)
[![Docker Stars](https://img.shields.io/docker/stars/thespad/whisparr.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&label=stars&logo=docker)](https://hub.docker.com/r/thespad/whisparr)

[![ci](https://img.shields.io/github/actions/workflow/status/thespad/traefik-crowdsec-bouncer/call-build-image.yml?labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github&label=Build%20Image)](https://github.com/thespad/traefik-crowdsec-bouncer/actions/workflows/call-build-image.yml)

## Supported Architectures

Our images support multiple architectures such as `x86-64`, `arm64` and `armhf`.

Simply pulling `ghcr.io/thespad/traefik-crowdsec-bouncer` should retrieve the correct image for your arch.

The architectures supported by this image are:

| Architecture | Available | Tag |
| :----: | :----: | ---- |
| x86-64 | ✅ | latest |
| arm64 | ✅ | latest |
| armhf | ✅ | latest |

## Prerequisites

* [Docker](https://docs.docker.com/get-docker/) and [Docker-compose](https://docs.docker.com/compose/install/) installed.
* Traefik v2.x
* CrowdSec running natively or in a container and configured to read logs from Traefik

## Application Setup

1. Get a bouncer API key from CrowdSec with command `docker exec crowdsec cscli bouncers add bouncer-traefik`
2. Copy the API key printed. You **_WON'T_** be able the get it again.
3. Paste this API key as the value for bouncer environment variable `CROWDSEC_BOUNCER_API_KEY`, or use an `.env` file.
4. Set the other environment variables as required (see below for details).
5. Start bouncer.
6. Visit a site proxied by Traefik and confirm you can access it.
7. In another console, ban your IP with command `docker exec crowdsec cscli decisions add --ip <your ip> -R "Test Ban"`, modify the IP with your address.
8. Visit the site again, in your browser you will see "Forbidden" since this time since you've been banned.
9. Unban yourself with `docker exec crowdsec cscli decisions delete --ip <your IP>`
10. Visit the site one last time, you will have access to the site again.

### Traefik Setup

Create a Forward Auth middleware, i.e.

```yml
    middleware-crowdsec-bouncer:
      forwardauth:
        address: http://crowdsec-bouncer-traefik:8080/api/v1/forwardAuth
        trustForwardHeader: true
```

Then apply it either to individual containers you wish to protect or as a default middlware on the Traefik listener.

## Parameters

The webservice configuration is made via environment variables:

* `CROWDSEC_BOUNCER_API_KEY`            - CrowdSec bouncer API key required to be authorized to request local API (required)
* `CROWDSEC_AGENT_HOST`                 - Host and port of CrowdSec agent, i.e. crowdsec-agent:8080 (required)
* `CROWDSEC_BOUNCER_SCHEME`             - Scheme to query CrowdSec agent. Expected value: http, https. Default to http
* `CROWDSEC_BOUNCER_LOG_LEVEL`          - Minimum log level for bouncer. Expected value [zerolog levels](https://pkg.go.dev/github.com/rs/zerolog#readme-leveled-logging). Default to 1
* `CROWDSEC_BOUNCER_SKIPRFC1918`        - Don't send RCF1918 (Private) IP addresses to the LAPI to check ban status. Expected value: "true", "false" . Default to "true"
* `CROWDSEC_BOUNCER_REDIRECT`           - Optionally redirect instead of giving 403 Forbidden. Accepts relative or absolute URLs but must not be protected by the bouncer or you'll get a redirect loop. Default to null
* `PORT`                                - Change listening port of web server. Default listen on 8080
* `GIN_MODE`                            - By default, run app in "debug" mode. Set it to "release" in production
* `TRUSTED_PROXIES`                     - Can accept a list of IP addresses in CIDR format, delimited by ','. Default is 0.0.0.0/0

## Exposed routes

The webservice exposes some routes:

* GET `/api/v1/forwardAuth`             - Main route to be used by Traefik: query CrowdSec agent with the header `X-Real-Ip` as client IP`
* GET `/api/v1/ping`                    - Simple health route that respond pong with http 200`
* GET `/api/v1/healthz`                 - Another health route that query CrowdSec agent with localhost (127.0.0.1)`
* GET `/api/v1/metrics`                 - Prometheus route to scrap metrics

## Support Info

* Shell access whilst the container is running: `docker exec -it traefik-crowdsec-bouncer /bin/bash`
* To monitor the logs of the container in realtime: `docker logs -f traefik-crowdsec-bouncer`

## Image Update Notifications - Diun (Docker Image Update Notifier)

* We recommend [Diun](https://crazymax.dev/diun/) for update notifications. Other tools that automatically update containers unattended are not recommended or supported.

## Building locally

If you want to make local modifications to these images for development purposes or just to customize the logic:

```shell
git clone https://github.com/thespad/traefik-crowdsec-bouncer.git
cd traefik-crowdsec-bouncer
docker build \
  --no-cache \
  --pull \
  -t ghcr.io/thespad/traefik-crowdsec-bouncer:latest .
```

## Versions

* **15.02.22:** - Initial Release.
