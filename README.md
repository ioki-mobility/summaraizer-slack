[![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ioki-mobility/summaraizer-slack/blob/main/LICENSE)

# summaraizer Slack

## What?

summaraizer-slack summarizes your Slack thread discussions.

## How?

There are three steps to get this up and running:
1. Set up the Slack integration
2. Deploy the code to a server
3. Invite the bot to your Slack channels (`/invite @BotName`)

If this was successfully done, you can mention the bot in a thread 
and ask it to summarize the thread with `@BotName summarize please`.

### Slack Integration

1. Go to `api.slack.com/apps`, and click on `Create New App`
2. Go to `OAuth & Permissions` and add the following Scopes to the **Bot Token**:
   * `app_mentions:read`
   * `chat:write`
   * `channels:history`
   * `groups:history`
   * `mpim:history`
   * `im:history`
3. Click on "Install to Workspace"
4. Accept everything (or wait until your Slack Admin approved it)
5. [Deploy the code to a server](#deployment)
6. Use the deployment URL, go back to the Slack integration and click `Event subscriptions`
7. Enable it and paste the deployment URL in it (for validation a small `Verifired ✅` should appear above it)
8. In the section `Subscribe to bot events`, enable `app_mention`
9. Reinstall the Slack Bot to your workspace

### Deployment

Before deploying the code, you have follow the steps 1-4 from the [Slack Integration](#slack-integration) section.
After you deployed the code, you can continue with the steps 6-9 steps.

The following sections shows how to deploy the code to Vercel, as a Docker Image or on a Webserver.

However, all options requires the `SLACK_BOT_TOKEN` 
as well as the `SLACK_SIGNING_SECRET` environment variable to be set.

Depending on which AI provider you want to use, 
you have to set the `OPENAI_API_TOKEN` **or** `OLLAMA_URL` environment variable as well.

The value for the `SLACK_BOT_TOKEN` can be found in the Slack integration `OAuth & Permissions`.
The value for the `SLACK_SIGNING_SECRET` can be found in the Slack integration `Basic Information`.

#### Vercel

To deploy the code to Vercel, you need to have the Vercel CLI installed.
Then you can simply deploy it with:

```bash
vercel
```

Note that you have to set the environment variables in the Vercel dashboard (`Settings` -> `Environment Variables`).
After that, you have to **redeploy**. Otherwise, the env. variables changes won't take effect.

#### Docker

To build the Docker image, you can run the following command:

```bash  
docker build -t summaraizer-slack .
```

To run that Docker image, you can use the following command:

```bash
docker run -p 8080:8080 -e SLACK_BOT_TOKEN=your-token -e SLACK-SIGNING_SECRET=your-secret -e OPENAI_API_TOKEN=your-token summaraizer-slack
```


#### Webserver

To run the code on a webserver, you need to have Go installed.
Then you can simply run the following command:

```bash
go run cmd/summaraizer-slack/main.go
```

Optional, you can specify the port with the `--port` flag.

```bash
go run cmd/summaraizer-slack/main.go --port 1234
```

## Why?

Ever run into a situation where you got a bit overwhelmed with the amount of comments
inside a Slack thread?
Or you just came back from a nice 3-week vacation trip and now have to read
a 75-comment long thread?

Such problems are a thing of the past now!
Just run the summaraizer over it to get a summary of all comments
and get back on track faster than ever before!

## Testing

When you run the server locally, you need first to verify the domain with Slack.
For this you can use `ngrok` or [`bore`](https://github.com/ekzhang/bore) to create a tunnel to your local server.

Both will give you a free URL that you can paste into the Slack integration.
Once confirmed, you can disable the tunnel again.

### Docker

When you want to run the Docker image together with a local Ollama instance, you can use the following command:

```bash
docker run -p 8080:8080 -e SLACK_BOT_TOKEN=your-token -e OLLAMA_URL=http://host.docker.internal:11434 summaraizer-slack
```

### Fake Slack Event

To fake a Slack event, you can use the following `curl` command:

```bash
curl -X POST http://localhost:8080 -H "Content-Type: application/json" -d '{
  "type": "event_callback",
  "event": {
    "type": "app_mention",
    "user": "U123456",
    "text": "Please summarize this",
    "thread_ts": "[Thread Timestamp]", 
    "channel": "[Channel ID]",
  }
}'
```

> [!IMPORTANT]
> When faking this, you might want to remove the signature verification from the code.
> You can change this inside the [`slack/slack.go`](slack/slack.go) file.

## Release

1. Navigate to the [Actions tab](../../actions) in your repository.
2. Select the ["Create Release" workflow](../../actions/workflows/publishing.yml).
3. Click on "Run workflow" and enter the new version number (e.g., `1.2.3`).

The new version tag will be created and pushed.
A draft GitHub release will be generated. 
You can view the releases on the [Releases page](../../releases/latest).
