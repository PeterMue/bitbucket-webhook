Bitbucket Webhook
=================

Simple webhook listener for Bitbucket webhooks that executes configurable shell commands when a webhook is triggered.

# Usage

```sh
    bitbucket-webhook [-config "<config-file>"] [-listen "<[host]:port>"]] [-secret "<secret>"]
```

> **Note**: Command line arguments take higher precedance over configuration file options.


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
BODY='{ "actor" : { "name" : "Peter Müller" }}' 
SIG="$(echo -n $BODY | openssl dgst -sha256 -hmac secret)" 
curl -XPOST -H "X-Request-Id: $(uuidgen)" -H "X-Event-Key: pr:merged" -H "X-Hub-Signature: $SIG" -v --data "$BODY" http://localhost:3000/webhook
```

*or as one-liner*

```bash
BODY='{ "actor" : { "name" : "Peter Müller" }}' SIG="$(echo -n $BODY | openssl dgst -sha256 -hmac secret)" bash -c 'curl -XPOST -H "X-Request-Id: $(uuidgen)" -H "X-Event-Key: pr:merged" -H "X-Hub-Signature: $SIG" -v --data "$BODY" http://localhost:3000/webhook'
```
