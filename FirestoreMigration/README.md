# Firestoreのマイグレーション

## 新設計


一旦、スケジュール周りの以下みたいにマイグレーション

```
/version/1/users/:id
    └──calendars/:id
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
