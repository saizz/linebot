# linebot

## required

- go get -u github.com/line/line-bot-sdk-go/linebot
- go get -u github.com/joho/godotenv
- go get -u google.golang.org/appengine
- go get -u golang.org/x/oauth2
- go get -u cloud.google.com/go/storage
- go get -u google.golang.org/api... or go get -u google.golang.org/api/calendar/v3
- go get -u github.com/nfnt/resize

or 

```
cd src
dep ensure
```

## get start

### add file

make src/line.env file

``` :line.env
LINE_BOT_CHANNEL_SECRET=<CHANNEL_SECRET>
LINE_BOT_CHANNEL_TOKEN=<CHANNEL_TOKEN>
```

### modify file

set `application` in app.yaml

### deploy

set GOPATH to checkout root directory.

```sh
cd src
goapp deploy .
```

## great thanks

I am deeply grateful to this site [LINE Developer Trialを GAE/Go で始める](http://qiita.com/koki_cheese/items/66980888d7e8755d01ec).
