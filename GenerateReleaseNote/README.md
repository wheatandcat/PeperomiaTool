# リリースノート自動生成

Pull Requestからリリースノートを自動生成する

## 準備

```
$ go mod download
```

## 設定

以下のコマンドで設定ファイルをコピー

```
$ mv config.template.toml config.toml 
```

config.tomlを書き換え

```
[GitHub]
token = ""
owner = ""
repository = ""
label = "リリースノート"
startText = "## 対応内容"
endText = "## 開発用メモ"
```



| 名前 | 内容 |
----|---- 
|  token  |  GitHub APIのトークンを設定  |
|  owner  |  オーナー名を設定  |
|  repository  |  レポジトリー名を設定  |
|  label  |  リリースノートで収集するのPull Requestラベル名を設定  |
|  startText  |  Pull RequestのBodyからの切り抜きの開始文字  |
|  endText  |  Pull RequestのBodyからの切り抜きの終了文字  |

## 実行

MILESTONEはタグ名を指定

```
$ MILESTONE=v2.0.4 go run main.go | pbcopy
```

↑を貼り付けるとマークダウン形式でこんな感じに表示

<img src="./doc/screen.png" width="800">