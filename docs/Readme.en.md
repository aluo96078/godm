<p align="right">
  🌐 [English](docs/Readme.en.md) | [繁體中文](README.md)
</p>

# GODM：MongoDB for Go 的簡易查詢映射器

## 🧩 簡介

GODM（Go Object-Document Mapper）是一個用於 MongoDB 的輕量級查詢封裝工具，使用 Go 語言實作。它提供類似 ORM 的開發體驗，並針對常見的查詢條件與鏈式操作做了簡化，幫助你快速構建資料模型並執行 CRUD、聚合、事務等操作。

核心實作位於 [`pkg/odm`](./pkg/odm)，使用範例可見於 [`examples/`](./examples)。

---

## ✨ 功能特色

- 🚀 鏈式查詢語法（Where, OrWhere, WhereIn 等）
- 🔧 自動關聯資料模型與集合（支援自定義集合名與資料庫名）
- 💾 支援 CRUD 與 BulkCreate
- 🧠 支援複雜查詢條件組合（AND / OR）
- 🔁 支援 MongoDB 聚合管道
- 💼 內建事務封裝 `WithTransaction`
- 🧪 簡潔易測試，模組化設計便於擴展

---

## 🛠 使用方式（以 User 模型為例）

### 方法覆寫（回傳自定義型別）

GODM 方法預設回傳 `*GODM`，但若您希望保留自定義模型型別（例如 `*User`）以便鏈式操作時能存取欄位，可以在模型中覆寫對應方法，例如：

```go
func (o *User) SetCollectionName(name string) *User {
    o.CollectionName = name
    if o.Collection != nil {
        o.Collection = odm.MongoClient.Database(o.DbName).Collection(name)
    }
    return o
}
```

這樣就可以保留類型一致性：

```go
u := NewUser().SetCollectionName("custom_users")
fmt.Println(u.Name) // 可直接使用 *User 欄位
```

### 建立與查詢

```go
user := NewUser()
user.Name = "Test"
user.Email = "test@example.com"
_ = user.Create()

// 查詢第一筆資料
err := user.Where("email", "=", "test@example.com").First()
```

### 聚合與事務操作

```go
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"email": bson.M{"$ne": ""}}}},
    {{"$group", bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}}},
}
var result []bson.M
_ = user.Aggregate(pipeline, &result)

_ = user.WithTransaction(func(sess mongo.SessionContext) error {
    return user.Update(bson.M{"name": "Updated"})
})
```

### 更多查詢示例

#### 使用 `WhereID`

```go
// 根據 MongoDB ObjectID 查詢文件
user := NewUser()
_ = user.WhereID("65f74c3a09c7a8f812345678").First()
```

#### 使用 OR 查詢

```go
var users []User
err := NewUser().
    OrWhere("name", "=", "Alice").
    OrWhere("email", "=", "bob@example.com").
    All(&users)
```

#### 使用 WhereIn 與欄位選取

```go
var users []User
err := NewUser().
    WhereIn("name", []interface{}{"Alice", "Bob"}).
    Select("name").
    Exclude("email").
    All(&users)
```

#### 使用分頁與排序

```go
var users []User
err := NewUser().
    OrderBy("name", true).
    Offset(10).
    Limit(10).
    All(&users)
```

#### 使用自定義上下文（含超時）

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user := NewUser().WithContext(ctx)
_ = user.Where("email", "=", "timeout@example.com").First()
```

---

## 💡 靈感來源

GODM 的設計靈感來自於 [Laravel Eloquent ORM](https://laravel.com/docs/eloquent)，試圖為 Golang 帶來一種熟悉且簡潔的資料查詢體驗。它並非 ORM，而是專注於查詢構建、結果解碼與事務包裝，適合喜歡鏈式語法與輕量抽象的使用者。

## 📂 專案結構

```
godm/
├── examples/        # 使用範例：main.go, user.go
├── pkg/
│   └── odm/         # GODM 核心實作（已模組化）
│       ├── aggregate.go
│       ├── config.go
│       ├── context.go
│       ├── crud.go
│       ├── model.go
│       ├── operator.go
│       ├── query.go
│       ├── transaction.go
│       └── util.go
├── go.mod
└── README.md
```

---

## 📄 授權

本專案採用 [MIT License](./LICENSE) 授權。
