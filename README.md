Bitbucket Webhook
=================

Simple webhook listener for Bitbucket webhooks that executes configurable shell commands when a webhook is triggered.

# Usage

```sh
    bitbucket-webhook [-config "<config-file>"] [-listen "<[host]:port>"]] [-secret "<secret>"]
```

> **Note**: Command line arguments take higher precedance over configuration file options.

## Docker

```sh
docker run -it --volume ./config.yaml:/config.yaml -p 3000:3000 ghcr.io/petermue/bitbucket-webhook:latest -config /config.yaml -secret "your-secret"
```

Or using environment variables with a configuration file:

```sh
docker run -it \
  --volume ./config.yaml:/config.yaml \
  -p 3000:3000 \
  ghcr.io/petermue/bitbucket-webhook:latest
```

The image is available on [GitHub Container Registry (GHCR)](https://github.com/PeterMue/bitbucket-webhook/pkgs/container/bitbucket-webhook).


# Configuration File

```yaml
listen: ":3000"
secret: "my$ecret"
hooks:
- event: "event:key"
  async: false
  command: /bin/bash
  args:
  - "-c"
  - echo "Hello World"
```

## Templating

The `args` support go template syntax and were populated with the entire Webhook payload object.

*Example*
```yaml
hooks:
- event: "pr:merged"
  command: /bin/bash
  args:
  - "-c"
  - echo "There's a merged PR from {{ .author.name }}"
```

See [Bitbucket: Event Payload](https://confluence.atlassian.com/bitbucketserver/event-payload-938025882.html) for more Details.

# Example Webhook Requests

> **Note**: The body is just a minimalistic - but correct - subset of the actual bitbucket `pr:merged` webhook payload.

```bash
SECRET='secret'
BODY='{ "actor" : { "name" : "Peter Müller" }}' 
SIG="$(echo -n $BODY | openssl dgst -sha256 -hmac $SECRET -binary | xxd -p -c 256)" 
curl -XPOST -H "X-Request-Id: $(uuidgen)" -H "X-Event-Key: pr:merged" -H "X-Hub-Signature: $SIG" -v --data "$BODY" http://localhost:3000/webhook
```

*or as one-liner*

```bash
SECRET='secret' BODY='{ "actor" : { "name" : "Peter Müller" }}' SIG="$(echo -n $BODY | openssl dgst -sha256 -hmac $SECRET -binary | xxd -p -c 256)" bash -c 'curl -XPOST -H "X-Request-Id: $(uuidgen)" -H "X-Event-Key: pr:merged" -H "X-Hub-Signature: $SIG" -v --data "$BODY" http://localhost:3000/webhook'
```
