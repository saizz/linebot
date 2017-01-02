# linebot

## required

- go get -u github.com/line/line-bot-sdk-go/linebot
- go get -u github.com/joho/godotenv
- go get -u google.golang.org/appengine
- go get -u golang.org/x/oauth2
- go get -u cloud.google.com/go/storage

## get start

1. make src/line.env file

```:line.env
LINE_BOT_CHANNEL_SECRET=<CHANNEL_SECRET>
LINE_BOT_CHANNEL_TOKEN=<CHANNEL_TOKEN>
```

1. set `application` in app.yaml

1. deploy

```sh
cd src
goapp deploy .
```