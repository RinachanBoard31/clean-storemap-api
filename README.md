# clean-restaurant-map

## Setup

### DB

```
$ docker compose -f environments/docker-compose.yml up -d
```

### OAuth

### GOOGLE_CLIENT_IDとGOOGLE_CLIENT_SECRETの発行方法
一例になります。制限等は全く同じではなくても構いません。

#### Oauth同意画面の構成(「認証情報の追加」ボタンが存在するなら不要)
1. GoogleCloudの「APIとサービス」->「認証情報」にて画面上部の「Oauth同意画面を構成」をクリック
2. GoogleWorkspaceユーザではないのでUserTypeを「外部」にする
3. 「アプリ名」「ユーザーサポートメール」を入力する。アプリのロゴ、アプリのドメインは必要があれば。「デベロッパーの連絡先情報」を入力する
4. スコープ、テストユーザは特に入力なし

#### 認証情報の追加
1. 「アプリケーションの種類」をウェブアプリケーションにする
2. 「承認済みの JavaScript 生成元」をバックエンドのものにする(http://localhost:PORT番号)
3. 「承認済みのリダイレクト URI」をバックエンドのものにする(http://localhost:PORT番号/auth/signup)
4. 「クライアント ID」と「クライアント シークレット」を環境変数にセットする
5. 環境変数のBACKEND_URLに"http://localhost:PORT番号"をセットする

**PORT番号はバックエンドのものにする**

## Run

```
$ go build cmd/main.go
$ ./main
```

## Usage

### Get store
```
$ curl http://localhost:8080
```

### Login user
- dbにlogin_user_api_example.jsonに記載されているemailがある => {}
- dbにlogin_user_api_example.jsonに記載されているemailがない => エラーメッセージとなる
```
$ curl -H "Content-Type: application/json" -X POST -d "@example/login_user_api_example.json" http://localhost:8080/login
```

### Auth
```
# ログイン用
$ curl http://localhost:8080/auth?accessedType=login
# サインアップ用
$ curl http://localhost:8080/auth?accessedType=signup
```
