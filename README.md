# Misskey Emoji Bot
Misskeyの絵文字登録をDiscord上で申請できるDiscord用Botです！

## 機能
✅絵文字申請をボタンで簡単に作成

✅絵文字の各種情報を簡単に設定可能

✅絵文字の申請から承認・不認可、アップロードまでを自動化！

## コマンド
`/init`: 初期設定を行います。

`/ni_rilana`: このBotに関する情報を表示します

`/change_emoji_detail <property> <value>`: 申請後に申請可否チャンネルにおいて絵文字の属性を変更できます。propertyには任意の属性を指定し、valueにはその属性に合う値を入力してください。（NSFWについては現状設定できません。)


今後も機能を追加予定です！

## 使用方法
1. docker imageまたはgit cloneにてリポジトリをローカル環境へ。
- docker環境の方はdocker-compose.yaml及びDockerfileを参照ください。
- ローカルの場合、`go run ./cmd/`で実行できます

2. settings.env内の設定値を埋めてください
- `guild_id`: DiscordのサーバーID
- `bot_token`: BotのToken(事前にDiscordアプリケーションを作成してください。)
- `application_id`: BotのアプリケーションID
- `moderator_role_id`: 絵文字申請処理を行うサーバー内ロールID。絵文字申請の可否をモデレーターのロール数で制御する際に判定します。
- `bot_role_id`: 絵文字Bot用ロールID
- `moderation_channel_name`: 絵文字が申請された際に内容が送信されるチャンネル名。重複する場合そのチャンネルに設定されます。
- `misskey_token`: 絵文字アップロードに利用します。絵文字を操作する都合上、MisskeyのTokenは必要な権限が付与されていることを確認してください。
- `save_path`: 絵文字が申請された際に画像を保存するパス
- `debug`: デバッグ表示を行うか。(Discord上には反映されません。)


## 使用ライブラリ
- [go-misskey](https://github.com/niwaniwa/go-misskey)

## License
- [GPL-3.0 license](https://github.com/niwaniwa/MisskeyEmojiBot/blob/main/LICENSE)

## Contact
- Misskey: [@ni_rilana](https://misskey.niri.la/@ni_rilana)
- Twitter: [@ni_rilana](https://twitter.com/ni_rilana)
