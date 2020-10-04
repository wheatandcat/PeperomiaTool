# Firestoreのマイグレーション

## 新設計


元の設計

```
├──users/:id
├──calendars/:id
├──items/:id
├──itemDetails/:id
└──expoPushTokens/:id
```

以下みたいにマイグレーション

```
/version/1/users/:id
    ├──expoPushTokens/:id
    └──calendars/:date
        └──items/:id
            └──itemDetails/:id
```

## 準備

```
$ go mod download
```

or

```
$ go mod tidy
```

## マイグレーション実行

```
$ go run main.go
```
