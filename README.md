# Sabet

Sabet (Static Automation Bot Event Tool; just kidding sabet mean fixed or static) is a tool that helps you to automate your events without the need to run it on a server for ever. You can just schedule a github action to run it periodically and it will do the rest. Or you can just run it manually on your PC whenever you want to update it.

## Installation

You could either install it using `go install`:

```bash
go install github.com/sabet/sabet
```

or you could just [download the binary](https://github.com/mehdy/sabet/releases).

## Usage

You need to create some yaml files to configure the jobs you need.

For example, let's create some jobs to deliver RSS feeds to a Telegram channel:

```yaml
# news.yaml
type: RSS
metadata:
  name: news-reader
  labels:
    app: news-reader
  run:
    when: always
spec:
  sources:
    - url: https://mehdy.me/en/index.xml
    - url: https://mehdy.me/fa/index.xml
store:
  type: fs
  path: news-reader
```

```yaml
# telegram.yaml
type: Telegram
metadata:
  name: news-reader-notifier
  labels:
    app: news-reader-notifier
  run:
    selector:
      app: news-reader
spec:
  tokenEnv: TELEGRAM_NEWS_READER_BOT_TOKEN
  channel: "@expfield"
  template:
    parseMode: HTML
    text: |
      <a href="{{ .link }}">{{ .link }}</a>
```

And you need to provide an environment variable named TELEGRAM_NEWS_READER_BOT_TOKEN containing the bot token.

and finally you need to run `sabet` inside the directory where you have the yaml files:

```bash
sabet
```
