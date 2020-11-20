# REFERENCE
- https://github.com/serverless/serverless-golang
- https://github.com/sminamot/sheets-go-example/blob/master/main.go
- https://sminamot-dev.hatenablog.com/entry/2019/12/05/204403

## Quick Start

1. リポジトリを初期化
```
serverless create -u https://github.com/KotaroYamazaki/slack-notifer -p slack-notifer
```
2. Secret Manager に環境変数をセット
WEBHOOK_URL, SHEET_ID, SECRET

3. Compile function

```
cd slack-notifer
GOOS=linux go build -o bin/main
```

4. デプロイ

```
sls deploy
```

5. 確認

```
sls invoke -f slack
````

