# Deployment
1. install gcloud sdk
2. ensure `app.yaml` exists and is configured properly (check `app.yaml.template`)
3. set env variable BOT_TOKEN="TELEGRAM_BOT_TOKEN"
3. run `gcloud init`
4. run `gcloud app deploy`
5. run `./webhook.sh`