# mattermost-buddybot

A Mattermost bot to randomly pair up team members for a private chat, written in Go.

## Configuration

| Env Variable | Description | Example |
| --- | --- | --- |
| SERVER_URL | The URL of the Mattermost server, including the protocol (i.e. https) | `https://my-mattermost.com` |
| BOT_USERNAME | The username (or email) the bot should use to login | `buddybot` |
| BOT_PASSWORD | The password the bot should use to login | `password123` |
| DEBUG_CHANNEL | (optional) The name of the channel where the bot can write debug messages. | `debug-channel` |
| DEBUG_CHANNEL_TEAM | (optional) The name of the team where the debug channel is located. | `debug-channel` |

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
    -e BOT_PASSWORD=password123 \
    -e DEBUG_CHANNEL=debug-channel \
    -e DEBUG_CHANNEL_TEAM=my-team \
    mattermost-buddybot
```

## Kubernetes

To deploy the Docker container to Kubernetes, you'll first need to create a secret with the bot's login credentials.

```bash
kubectl -n $NAMESPACE create secret generic mattermost-buddybot \
    --from-literal=server-url=https://my-mattermost.com \
    --from-literal=username=buddybot \
    --from-literal=password=password123 \
    --from-literal=debug-channel=debug-channel \
    --from-literal=debug-channel-team=my-team
```

Then deploy the container via the provided manifest:

```bash
kubectl -n $NAMESPACE -f kubernetes.yaml
```

To check the bot is running, mention the bot in any channel like: `@buddybot, are you running?`.