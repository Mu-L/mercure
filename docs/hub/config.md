# Configuration

The Mercure.rocks hub is a custom build of the [Caddy web server](https://caddyserver.com/) including the Mercure.rocks module.

Read [Caddy web server's getting started guide](https://caddyserver.com/docs/getting-started) to learn the basics.

While all supported ways to configure Caddy are also supported by the Mercure.rocks Hub, the easiest one is [to use a `Caddyfile`](https://caddyserver.com/docs/quick-starts/caddyfile).
A default `Caddyfile` is provided in [the archive containing the Mercure.rocks Hub](install.md).

A minimal `Caddyfile` for the Mercure hub looks like this:

```Caddyfile
# The address of your server
localhost {
	mercure {
		# Publisher JWT key
		publisher_jwt !ChangeThisMercureHubJWTSecretKey!
		# Subscriber JWT key
		subscriber_jwt !ChangeThisMercureHubJWTSecretKey!
	}

	respond "Not Found" 404
}
```

Caddy will automatically generate a Let's Encrypt TLS certificate for you! So you can use HTTPS.
To disable HTTPS, prefix the name of the server by `http://`:

```Caddyfile
http://my-domain.test:3000 {
    # ...
}
```

Note that HTTPS is automatically disabled if you set the server port to 80.

## Directives

The following Mercure-specific directives are available:

| Directive                                                  | Description                                                                                                                                                                                                                                     | Default                |
|------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------|
| `publisher_jwt <key> [<algorithm>]`                        | the JWT key and algorithm to use for publishers, can contain [placeholders](https://caddyserver.com/docs/conventions#placeholders)                                                                                                              |                        |
| `subscriber_jwt <key> [<algorithm>]`                       | the JWT key and algorithm to use for subscribers, can contain [placeholders](https://caddyserver.com/docs/conventions#placeholders)                                                                                                             |                        |
| `publisher_jwks_url`                                       | the URL of the JSON Web Key Set (JWK Set) URL (provided by identity providers such as Keycloak or AWS Cognito) to use for validating publishers JWT (take precedence over `publisher_jwt`)                                                      |                        |
| `subscriber_jwks_url`                                      | the URL of the JSON Web Key Set (JWK Set) URL to use for validating subscribers JWT (take precedence over `subscriber_jwt`)                                                                                                                     |                        |
| `anonymous`                                                | allow subscribers with no valid JWT to connect                                                                                                                                                                                                  | `false`                |
| `publish_origins <origins...>`                             | a list of origins allowed publishing, can be `*` for all (only applicable when using cookie-based auth)                                                                                                                                         |                        |
| `cors_origins <origin...>`                                 | a list of allowed CORS origins, ([troubleshoot CORS issues](troubleshooting.md#cors-issues))                                                                                                                                                    |                        |
| `cookie_name <name>`                                       | the name of the cookie to use for the authorization mechanism                                                                                                                                                                                   | `mercureAuthorization` |
| `subscriptions`                                            | expose the subscription web API and dispatch private updates when a subscription between the Hub and a subscriber is established or closed. The topic follows the template `/.well-known/mercure/subscriptions/{topicSelector}/{subscriberID}`  |                        |
| `heartbeat`                                                | interval between heartbeats (useful with some proxies, and old browsers), set to `0s` disable                                                                                                                                                   | `40s`                  |
| `transport <name> [{ <options...> }]`                      | The transport to use. Options are transport-specific. See also [the cluster mode](cluster.md)                                                                                                                                                   | `bolt://mercure.db`    |
| `dispatch_timeout <duration>`                              | maximum duration of the dispatch of a single update, set to `0s` disable                                                                                                                                                                        | `5s`                   |
| `write_timeout <duration>`                                 | maximum duration before closing the connection, set to `0s` disable                                                                                                                                                                             | `600s`                 |
| `protocol_version_compatibility`                           | version of the protocol to be backward compatible with (only version 7 is supported)                                                                                                                                                            | disabled               |
| `demo`                                                     | enable the UI and expose demo endpoints                                                                                                                                                                                                         |                        |
| `ui`                                                       | enable the UI but do not expose demo endpoints                                                                                                                                                                                                  |                        |
| `topic_selector_cache <maxEntriesPerShard> [<shardCount>]` | Topic Selector cache configuration (see [golang-lru docs](https://github.com/hashicorp/golang-lru)) (pass `0` to disable the cache)                                                                                                             | `10000` `256`          |
| `transport_url <url>`                                      | **Deprecated: use `transport` instead.** URL representation of the transport to use. Use `local://local` to disable the history, (example `bolt:///var/run/mercure.db?size=100&cleanup_frequency=0.4`), see also [the cluster mode](cluster.md) | `bolt://mercure.db`    |

See also [the list of built-in Caddyfile directives](https://caddyserver.com/docs/caddyfile/directives).

## Environment Variables

The provided `Caddyfile` and the Docker image provide convenient environment variables:

| Environment variable            | Description                                                                                                                                                                                                          | Default value |
|---------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------|
| `GLOBAL_OPTIONS`                | the [global options block](https://caddyserver.com/docs/caddyfile/options#global-options) to inject in the `Caddyfile`, one per line                                                                                 |               |
| `CADDY_EXTRA_CONFIG`            | the [snippet](https://caddyserver.com/docs/caddyfile/concepts#snippets) or the [named-routes](https://caddyserver.com/docs/caddyfile/concepts#named-routes) options block to inject in the `Caddyfile`, one per line |               |
| `CADDY_SERVER_EXTRA_DIRECTIVES` | [`Caddyfile` directives](https://caddyserver.com/docs/caddyfile/concepts#directives)                                                                                                                                 |               |
| `SERVER_NAME`                   | the server name or address                                                                                                                                                                                           | `localhost`   |
| `MERCURE_PUBLISHER_JWT_KEY`     | the JWT key to use for publishers                                                                                                                                                                                    |               |
| `MERCURE_PUBLISHER_JWT_ALG`     | the JWT algorithm to use for publishers                                                                                                                                                                              | `HS256`       |
| `MERCURE_SUBSCRIBER_JWT_KEY`    | the JWT key to use for subscribers                                                                                                                                                                                   |               |
| `MERCURE_SUBSCRIBER_JWT_ALG`    | the JWT algorithm to use for subscribers                                                                                                                                                                             | `HS256`       |
| `MERCURE_EXTRA_DIRECTIVES`      | a list of extra [Mercure directives](#directives) inject in the Caddy file, one per line                                                                                                                             |               |
| `MERCURE_LICENSE`               | the license to use ([only applicable for the HA version](cluster.md))                                                                                                                                                |               |

## HealthCheck

The Mercure.rocks Hub provides a `/healthz` endpoint that returns a `200 OK` status code if the server is healthy.

Here is an example of how to use the health check in a Docker Compose file:

```yaml
# compose.yaml
services:
  mercure:
    # ...
    healthcheck:
      test: ["CMD", "wget", "-O-", "https://localhost/healthz"]
      timeout: 5s
      retries: 5
      start_period: 60s
```

## JWT Verification

JWT can be validated using HMAC and RSA algorithms.
In addition, it's possible to use JSON Web Key Sets (JWK Sets) (usually provided by OAuth and OIDC providers such as Keycloak or Amazon Cognito) to validate the keys.

When using RSA public keys for verification make sure the key is properly formatted and make sure to set the correct algorithm as second parameter of the `publisher_jwt` or `subscriber_jwt` directives (for example `RS256`).

Here is an example of how to use environments variables with an RSA public key.

Generate keys (if you don't already have them):

```console
ssh-keygen -t rsa -b 4096 -m PEM -f publisher.key
openssl rsa -in publisher.key -pubout -outform PEM -out publisher.key.pub

ssh-keygen -t rsa -b 4096 -m PEM -f subscriber.key
openssl rsa -in subscriber.key -pubout -outform PEM -out subscriber.key.pub
```

Start the hub:

```console
MERCURE_PUBLISHER_JWT_KEY=$(cat publisher.key.pub) \
MERCURE_PUBLISHER_JWT_ALG=RS256 \
MERCURE_SUBSCRIBER_JWT_KEY=$(cat subscriber.key.pub) \
MERCURE_SUBSCRIBER_JWT_ALG=RS256 \
./mercure run
```

## Bolt Adapter

| Option              | Description                                                                                                                                                        |
|---------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `path`              | path of the database file (default: `mercure.db`)                                                                                                                  |
| `bucket_name`       | name of the bolt bucket to store events (default: `updates`)                                                                                                       |
| `cleanup_frequency` | chances to trigger history cleanup when an update occurs, must be a number between `0` (never cleanup) and `1` (cleanup after every publication, default to `0.3`) |
| `size`              | size of the history (to retrieve lost messages using the `Last-Event-ID` header), set to `0` to never remove old events (default)                                  |

You can visualize and edit the content of the database using [boltdbweb](https://github.com/evnix/boltdbweb).

### Legacy URL

**This feature is deprecated: use the new `transport` directive instead**.

The [Data Source Name (DSN)](https://en.wikipedia.org/wiki/Data_source_name) specifies the path to the [bolt](https://github.com/etcd-io/bbolt) database as well as options. All options available as directive except `path` can be passed.

Below are common examples of valid DSNs showing a combination of available values:

```Caddyfile
# absolute path to `updates.db`
transport_url bolt:///var/run/database.db

# path to `updates.db` in the current directory
transport_url bolt://database.db

# custom options
transport_url bolt://database.db?bucket_name=demo&size=1000&cleanup_frequency=0.5
```

## Legacy Server

**The legacy server is deprecated and will be removed in the next version. Consider upgrading to the Caddy-based build.**

The legacy Mercure.rocks Hub is configurable using [environment variables](https://en.wikipedia.org/wiki/Environment_variable) (recommended in production, [twelve-factor app methodology](https://12factor.net/)), command-line flags and configuration files (JSON, TOML, YAML, HCL, envfile and Java properties files are supported).

Environment variables must be the name of the configuration parameter in uppercase.
Run `./mercure -h` to see all available command-line flags.

Configuration files must be named `mercure.<format>` (ex: `mercure.yaml`) and stored in one of the following directories:

- the current directory (`$PWD`)
- `~/.config/mercure/` (or any other XDG configuration directory set with the `XDG_CONFIG_HOME` environment variable)
- `/etc/mercure`

Most configuration parameters are hot reloaded: changes made to environment variables or configuration files are immediately taken into account, without having to restart the hub.

When using environment variables, list must be space separated. As flags parameters, they must be comma separated.

| Parameter                  | Description                                                                                                                                                                                                                                                                                                                                                                                  | Default                                                  |
|----------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------|
| `acme_cert_dir`            | the directory where to store Let's Encrypt certificates                                                                                                                                                                                                                                                                                                                                      |                                                          |
| `acme_hosts`               | a list of hosts for which Let's Encrypt certificates must be issued                                                                                                                                                                                                                                                                                                                          |                                                          |
| `acme_http01_addr`         | the address used by the acme server to listen on (example: `0.0.0.0:8080`)                                                                                                                                                                                                                                                                                                                   | `:http`                                                  |
| `addr`                     | the address to listen on (example: `127.0.0.1:3000`. Note that Let's Encrypt only supports the default port: to use Let's Encrypt, **do not set this parameter**.                                                                                                                                                                                                                            | `:http` or `:https` depending if HTTPS is enabled or not |
| `allow_anonymous`          | allow subscribers with no valid JWT to connect                                                                                                                                                                                                                                                                                                                                               | `false`                                                  |
| `cert_file`                | a cert file (to use a custom certificate)                                                                                                                                                                                                                                                                                                                                                    |                                                          |
| `key_file`                 | a key file (to use a custom certificate)                                                                                                                                                                                                                                                                                                                                                     |                                                          |
| `compress`                 | Use HTTP compression                                                                                                                                                                                                                                                                                                                                                                         | `false`                                                  |
| `cors_allowed_origins`     | a space-separated list of allowed CORS origins, can be `*` for all                                                                                                                                                                                                                                                                                                                           |                                                          |
| `debug`                    | debug mode, **dangerous, don't enable in production** (logs updates' content, why an update is not send to a specific subscriber and recovery stack traces)                                                                                                                                                                                                                                  | `false`                                                  |
| `demo`                     | demo mode (automatically enabled when `debug` is `true`) and enables ui at `https://example.com/.well-known/mercure/ui/`                                                                                                                                                                                                                                                                     | `false`                                                  |
| `dispatch_timeout`         | maximum duration of the dispatch of a single update, set to `0s` to disable                                                                                                                                                                                                                                                                                                                  | `5s`                                                     |
| `subscriptions`            | expose the subscription web API and dispatch private updates when a subscription between the Hub and a subscriber is established or closed. The topic follows the template `/.well-known/mercure/subscriptions/{topicSelector}/{subscriberID}`                                                                                                                                               | `false`                                                  |
| `heartbeat_interval`       | interval between heartbeats (useful with some proxies, and old browsers), set to `0s` to disable                                                                                                                                                                                                                                                                                             | `40s`                                                    |
| `jwt_key`                  | the JWT key to use for both publishers and subscribers                                                                                                                                                                                                                                                                                                                                       |                                                          |
| `jwt_algorithm`            | the JWT verification algorithm to use for both publishers and subscribers, e.g. `HS256` or `RS512`                                                                                                                                                                                                                                                                                           | `HS256`                                                  |
| `metrics_enabled`          | Enable the `/metrics` HTTP endpoint. Provide metrics for Hub monitoring in the OpenMetrics (Prometheus) format                                                                                                                                                                                                                                                                               | `false`                                                  |
| `metrics_addr`             | the address to listen on                                                                                                                                                                                                                                                                                                                                                                     | `127.0.0.1:9764`                                         |
| `publish_allowed_origins`  | a list of origins allowed to publish (only applicable when using cookie-based auth)                                                                                                                                                                                                                                                                                                          |                                                          |
| `publisher_jwt_key`        | must contain the secret key to valid publishers' JWT, can be omitted if `jwt_key` is set                                                                                                                                                                                                                                                                                                     |                                                          |
| `publisher_jwt_algorithm`  | the JWT verification algorithm to use for publishers, e.g. `HS256` or `RS512`                                                                                                                                                                                                                                                                                                                | `HS256`                                                  |
| `read_timeout`             | maximum duration for reading the entire request, including the body, set to `0s` to disable                                                                                                                                                                                                                                                                                                  | `5s`                                                     |
| `read__header_timeout`     | the amount of time allowed to read request headers, set to `0s` to disable                                                                                                                                                                                                                                                                                                                   | `5s`                                                     |
| `subscriber_jwt_key`       | must contain the secret key to valid subscribers' JWT, can be omitted if `jwt_key` is set                                                                                                                                                                                                                                                                                                    |                                                          |
| `subscriber_jwt_algorithm` | the JWT verification algorithm to use for subscribers, e.g. `HS256` or `RS512`                                                                                                                                                                                                                                                                                                               | `HS256`                                                  |
| `transport_url`            | URL representation of the history database. Provided database are `null` to disable history, `bolt` to use [bbolt](https://github.com/etcd-io/bbolt) (example `bolt:///var/run/mercure.db?size=100&cleanup_frequency=0.4`)                                                                                                                                                                   | `bolt://updates.db`                                      |
| `use_forwarded_headers`    | use the `X-Forwarded-For`, and `X-Real-IP` for the remote (client) IP address, `X-Forwarded-Proto` or `X-Forwarded-Scheme` for the scheme (`http` or `https`), `X-Forwarded-Host` for the host and the RFC 7239 `Forwarded` header, which may include both client IPs and schemes. If this option is enabled, the reverse proxy must override or remove these headers or you will be at risk | `false`                                                  |
| `write_timeout`            | maximum duration before closing the connection, set to `0s` to disable                                                                                                                                                                                                                                                                                                                       | `600s`                                                   |

If `acme_hosts` or both `cert_file` and `key_file` are provided, an HTTPS server supporting HTTP/2 connection will be started.
If not, an HTTP server will be started (**not secure**).
