# リリースノート自動生成

## セットアップ

```
$ go mod download
```

or

```
$ go mod tidy
```

## 実行

```
$ MILESTONE_NUMBER=9 go run main.go | pbcopy
```

↑を貼り付けるとマークダウン形式でこんな感じに表示

<img src="./doc/screen.png" width="800">