[![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ioki-mobility/summaraizer-slack/blob/main/LICENSE)

# summaraizer Slack

## What?

Summaraizer summarizes your Slack thread discussions.

## How?

There are a few steps required to add this integration to your Slack workspace:

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
5. Deploy this web app to `vercel` by running `vercel`
6. Copy the (stable) deployment URL and go back to the Slack integration and click `Event subscriptions`
7. Enable it and paste the vercel URL in it (for validation a small `Verifired ✅` should appear above it)
8. In the section `Subscribe to bot events`, enable `app_mention`
9. Reinstall the Slack Bot to your workspace
10. Go back to your `vercel` deployment -> settings and set the following secrets:
  * `SLACK_BOT_TOKEN` (can be found in the Slack integration `OAuth & Permissions`)
  * `OPENAI_API_TOKEN`
11. Redeploy the version instance (to add the new secrets to that deployment)

After this is done, you are able to invite your bot in your channels, group messages, etc.
by typing `/invite @BotName`.

Finally, you can summaraize your threads by typing `@BotName summarize please`.

## Why?

Ever run into a situation to got a bit overhelmed with the amound of comments
inside a Slack thread?
Or you just came back from a nice 3 weeks vacation trip and now has to read
a 75 comments long thread?

Such problems are from the past now!
Just run the summaraizer over it to get a summarization of all comments
to get on track faster than ever before!
