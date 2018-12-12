# mattermost-buddybot

A Mattermost bot to randomly pair up team members for a private chat, written in Go.


## Run

To run the code in a docker container run the following:

```bash
docker build -t mattermost-buddybot .
```

And run the following to run a container in the foreground, which will be removed with `CTRL-C`.

```bash
docker run -it --rm \
    -e SERVER_URL=https://my-mattermost.com \
    -e BOT_USERNAME=buddybot \
    -e BOT_PASSWORD=password \
    -e TEAMNAME=my-team \
    -e DEBUG_CHANNEL=debug-channel \
    mattermost-buddybot
```