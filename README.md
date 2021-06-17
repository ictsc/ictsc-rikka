# ictsc-rikka

score-server backend in ictsc 2021 summer.

## development

`docker-compose`と`make`を利用して開発することができます。

### 設定ファイル

設定ファイルは`cmd/rikka`以下に作成します。
`config.yaml`、`config.docker.yaml`は`.gitignore`で指定されているため、`config.yaml.example`を元に作成して下さい。

```
# 設定ファイルの例
cmd/rikka/config.yaml.example

# make run(ローカル実行時)に使用します
cmd/rikka/config.yaml

# make run-docker(docker実行時)に使用します
cmd/rikka/config.docker.yaml
```

### docker-compose.yml

| service | image           | 用途         | 永続化                 |
| ------- | --------------- | ------------ | ---------------------- |
| mariadb | mariadb:10.5.10 | データベース | mariadb:/var/lib/mysql |
| redis   | redis:6.2.4     | session 用   | なし                   |
| go      | golang:1.16     | 開発時実行用 | go-pkg:/go/pkg         |

### Makefile

以下のコマンドが使えます。

```
# ローカルでgo run実行
make run

# docker composeでgo run実行
make run-docker

# docker compose ps
make ps

# mariadb, redisを起動
make up

# mariadb, redis, goのコンテナとネットワークを削除
make down

# mariadbにmysqlコマンドでログイン
make mariadb

# redis-cliを起動
make redis-cli

# 自己署名証明書の作成
make self-signed-cert-and-key

```
