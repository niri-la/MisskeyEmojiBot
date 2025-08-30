# GitHub CLI チートシート（Issue/PRの眺め方・コメント取得）

> 前提：`gh` がインストール済み（`gh --version`）  
> 初回のみ：`gh auth login` で GitHub にログイン  
> リポジトリ直下で実行するか、`-R owner/repo` を付けて実行

## 0. よく使う前提設定

```bash
# ログイン（ブラウザ認証）
gh auth login

# 現在の認証/スコープ確認
gh auth status

# デフォルトの対象リポ（カレントのgitリモートを利用）
cd /path/to/repo

# 明示的に対象リポを指定したい時（どのコマンドにも付けられる）
# 例: -R owner/repo
```

---

## 1) Issue の「眺め方」

### 一覧を見る（フィルタ込み）
```bash
gh issue list
gh issue list -L 10
gh issue list -l "bug"
gh issue list -a "@me"
gh issue list --state open
gh issue list -S "is:open label:bug sort:updated-desc"
```

### 詳細を見る（本文・コメント付き）
```bash
gh issue view <number>
gh issue view <number> --comments
gh issue view <number> --web
```

### タイトルだけ確認
```bash
gh issue list --json number,title --jq '.[] | "\(.number): \(.title)"'
```

---

## 2) Issue コメントの取得

```bash
gh issue view <number> --comments
gh api repos/:owner/:repo/issues/<number>/comments -q '.'
gh api repos/:owner/:repo/issues/<number>/comments \\
  --jq '.[] | [.user.login, .created_at, (.body | gsub("\r";"") | split("\n")[0])] | @tsv'
```

---

## 3) Pull Request（PR）の「眺め方」

```bash
gh pr list
gh pr list -L 10
gh pr list -l "needs-review"
gh pr view <number>
gh pr view <number> --comments
gh pr view <number> --files
gh pr diff <number>
gh pr checks <number>
gh pr checkout <number>
```

---

## 4) PRコメントの取得

### 種別ごとにJSONで取得
```bash
gh api repos/:owner/:repo/issues/<number>/comments -q '.'
gh api repos/:owner/:repo/pulls/<number>/comments -q '.'
gh api repos/:owner/:repo/pulls/<number>/reviews -q '.'
```

### 抽出例
```bash
gh api repos/:owner/:repo/pulls/<number>/comments \\
  --jq '.[] | "\(.path):\(.line) | \(.user.login) | \((.body // "" | gsub("\r";"") | split("\n")[0]))"'

gh api repos/:owner/:repo/pulls/<number>/reviews \\
  --jq '.[] | [.user.login, .state, .submitted_at] | @tsv'

gh api repos/:owner/:repo/issues/<number>/comments \\
  --jq '.[] | [.user.login, .created_at, (.body // "" | gsub("\r";"") | split("\n")[0])] | @tsv'
```

---

## 5) 通知的に「自分が見るべきもの」

```bash
gh pr list -S "is:open review-requested:@me sort:updated-desc"
gh issue list -a "@me" -S "sort:updated-desc"
gh search issues "repo:owner/repo updated:>=2025-08-01" --limit 20
```

---

## 6) よく使うワンライナー

```bash
gh issue list -S "is:open sort:updated-desc" -L 10 --json number,title --jq '.[] | "\(.number)\t\(.title)"'
gh pr list --json baseRefName --jq 'group_by(.baseRefName)[] | {base: .[0].baseRefName, count: length}'
gh pr view <number> --json files --jq '.files[] | "\(.path)\t+\(.additions)\t-\(.deletions)"'
```

---

## 7) エイリアス

```bash
gh alias set il 'issue list -a "@me" -S "sort:updated-desc"'
gh alias set prr 'pr list -S "is:open review-requested:@me sort:updated-desc"'
gh alias set prf 'pr view $1 --files'
```

---

## 8) 小ワザ

- どのコマンドも `-R owner/repo` を付ければ別リポ対象にできる
- `--json ... --jq ...` を積極活用
- PRコメントは「Issueスレッドコメント」と「レビューコメント」がある点に注意

---

## 9) トラブル時

```bash
gh auth status
gh api rate_limit
```
