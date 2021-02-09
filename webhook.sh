#!/bin/bash

if [[ -z "$BOT_TOKEN" ]]; then
  echo "BOT_TOKEN variable is not set"
else
  curl -X POST \
    -H "Content-Type: application/json" \
    -d '{"url": "https://fortnit-elves-bot.ew.r.appspot.com"}' \
    "https://api.telegram.org/bot$BOT_TOKEN/setWebhook"
fi
