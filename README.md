# [thespad/traefik-crowdsec-bouncer](https://github.com/thespad/traefik-crowdsec-bouncer)

A http service to verify requests and bounce them according to decisions made by CrowdSec. Fork of [https://github.com/fbonalair/traefik-crowdsec-bouncer](https://github.com/fbonalair/traefik-crowdsec-bouncer) with extra features.

[![GitHub Release](https://img.shields.io/github/release/thespad/traefik-crowdsec-bouncer.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github&include_prereleases)](https://github.com/thespad/traefik-crowdsec-bouncer/releases)
![Commits](https://img.shields.io/github/commits-since/thespad/traefik-crowdsec-bouncer/latest?color=26689A&include_prereleases&logo=github&style=for-the-badge)
[![GitHub Stars](https://img.shields.io/github/stars/thespad/traefik-crowdsec-bouncer.svg?color=26689A&labelColor=555555&logoColor=ffffff&style=for-the-badge&logo=github)](https://github.com/thespad/traefik-crowdsec-bouncer)

## Prerequisites

* [Docker](https://docs.docker.com/get-docker/) and [Docker-compose](https://docs.docker.com/compose/install/) installed.
* Traefik v2.x
* CrowdSec running natively or in a container and configured to read logs from Traefik

## Application Setup

1. Get a bouncer API key from CrowdSec with command `docker exec crowdsec cscli bouncers add bouncer-traefik`
2. Pull the [docker image](https://github.com/thespad/docker-traefik-crowdsec-bouncer) for the bouncer: `docker pull ghcr.io/thespad/traefik-crowdsec-bouncer`
3. Copy the API key printed. You **_WON'T_** be able the get it again.
4. Paste this API key as the value for bouncer environment variable `CROWDSEC_BOUNCER_API_KEY`, or use an `.env` file.
5. Set the other environment variables as required (see below for details).
6. Start bouncer.
7. Visit a site proxied by Traefik and confirm you can access it.
8. In another console, ban your IP with command `docker exec crowdsec cscli decisions add --ip <your ip> -R "Test Ban"`, modify the IP with your address.
9. Visit the site again, in your browser you will see "Forbidden" since this time since you've been banned.
10. Unban yourself with `docker exec crowdsec cscli decisions delete --ip <your IP>`
11. Visit the site one last time, you will have access to the site again.

### Traefik Setup

Create a Forward Auth middleware, i.e.

```yml
    middleware-crowdsec-bouncer:
      forwardauth:
        address: http://crowdsec-bouncer-traefik:8080/api/v1/forwardAuth
        trustForwardHeader: true
```

Then apply it either to individual containers you wish to protect or as a default middlware on the Traefik listener.

## Configuration Environment Variables

| Parameter | Function |
| - | - |
| CROWDSEC_BOUNCER_API_KEY | CrowdSec bouncer API key (required). |
| CROWDSEC_AGENT_HOST | Host and port of CrowdSec LAPI agent, i.e. crowdsec-agent:8080 (required). |
| CROWDSEC_BOUNCER_SCHEME | Scheme to query CrowdSec agent. Allowed values: `http`, `https`. Default is `http`. |
| TRUSTED_PROXIES | IP addresses of upstream proxies. Can accept a list of IP addresses in CIDR format, delimited by ','. Default is `0.0.0.0/0`. |
| PORT | Change listening port of web server. Default is `8080`. |
| CROWDSEC_BOUNCER_LOG_LEVEL | Minimum log level for bouncer. Allowed values: [zerolog levels](https://pkg.go.dev/github.com/rs/zerolog#readme-leveled-logging). Default is `1`. |
| GIN_MODE | Operational mode for Gin framework. Set to `debug` for noisy log output. Default is `release`. |
| CROWDSEC_BOUNCER_SKIPRFC1918 | Don't send RCF1918 (Private) IP addresses to the LAPI to check ban status. Allowed values: `true`, `false` . Default is `true`. |
| CROWDSEC_BOUNCER_REDIRECT | Optionally redirect instead of giving 403 Forbidden. Accepts relative or absolute URLs but must not be protected by the bouncer or you'll get a redirect loop. Default is `null`. |
| CROWDSEC_BOUNCER_CLOUDFLARE | Use the `CF-Connecting-IP` header instead of `X-Forwarded-For`. This is useful if you're using Cloudflare proxying as `CF-Connecting-IP` will contain the real source address rather than the Cloudflare address. Allowed values: `true`, `false` . Default is `false`. |

## Exposed routes

The webservice exposes some routes:

* GET `/api/v1/forwardAuth`             - Main route to be used by Traefik: query CrowdSec agent with the header `X-Real-Ip` as client IP`
* GET `/api/v1/ping`                    - Simple health route that respond pong with http 200`
* GET `/api/v1/healthz`                 - Another health route that query CrowdSec agent with localhost (127.0.0.1)`
* GET `/api/v1/metrics`                 - Prometheus route to scrap metrics

## Versions

* **17.05.23:** - Add CROWDSEC_BOUNCER_CLOUDFLARE option.
* **01.05.23:** - Move docker image to its own repo.
* **01.05.23:** - Update deps.
* **01.05.23:** - Restructure repo.
* **26.04.23:** - Support CF forwarded IP headers.
* **15.02.22:** - Initial Release.
