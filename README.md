Bitbucket Webhook
=================

Simple webhook listener for Bitbucket webhooks that executes configurable shell commands when a webhook is triggered.

# Usage

```sh
    bitbucket-webhook [-config "<config-file>"] [-listen "<[host]:port>"]] [-secret "<secret>"]
```

Note: Command line arguments take higher precedance over configuration file options.


# Configuration File

```yaml
listen: ":3000"
secret: "my$ecret"
hooks:
- event: "event:key"
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