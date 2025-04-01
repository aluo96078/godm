[English](docs/Readme.en.md) | [繁體中文](Readme.md)

# GODM：MongoDB for Go 的簡易查詢映射器

## 📚 目錄
- [🧩 簡介](#-簡介)
- [✨ 功能特色](#-功能特色)
- [👀 Observer 機制（模型監聽）](#-observer-機制模型監聽)
  - [🎯 支援的事件](#-支援的事件)
  - [📦 使用方式](#-使用方式)
    - [定義 Observer](#定義-observer)
    - [模型自註冊（推薦）](#模型自註冊推薦)
  - [🌐 全域 Observer](#-全域-observer)
  - [🎛️ Observer 擴充功能](#-observer-擴充功能)
    - [✅ 事件過濾](#事件過濾只觀察某些事件)
    - [✅ 模型過濾](#模型過濾只監聽某些模型)
    - [✅ 優先順序](#優先順序priority)
    - [✅ 錯誤處理攔截](#錯誤處理攔截)
- [🛠 使用方式（以 User 模型為例）](#-使用方式以-user-模型為例)
  - [方法覆寫（回傳自定義型別）](#方法覆寫回傳自定義型別)
  - [建立與查詢](#建立與查詢)
  - [聚合與事務操作](#聚合與事務操作)
  - [更多查詢示例](#更多查詢示例)
    - [使用 WhereID](#使用-whereid)
    - [使用 OR 查詢](#使用-or-查詢)
    - [使用 WhereIn 與欄位選取](#使用-wherein-與欄位選取)
    - [使用分頁與排序](#使用分頁與排序)
    - [使用自定義上下文（含超時）](#使用自定義上下文含超時)
    - [判斷指定目標是否存在](#判斷指定目標是否存在)
- [🔗 關聯查詢（with 關聯預載入）](#-關聯查詢with-關聯預載入)
  - [User → Posts](#user--posts)
  - [Post → User](#post--user)
- [💡 靈感來源](#-靈感來源)
- [📂 專案結構](#-專案結構)
- [📄 授權](#-授權)


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
- 🔗 支援 with 預載入關聯資料（eager loading）
- 💼 內建事務封裝 `WithTransaction`
- 👀 內建 Observer 機制，支援模型級、全域、排序與過濾（Inspired by Laravel）
- 🧪 簡潔易測試，模組化設計便於擴展

## 👀 Observer 機制（模型監聽）

GODM 內建 Laravel Eloquent 式的 Observer 系統，可讓你在模型的 `Create`、`Update`、`Delete` 操作前後，自動觸發對應邏輯，適合用於資料驗證、日誌記錄、事件追蹤等情境。

### 🎯 支援的事件

- `creating` / `created`
- `updating` / `updated`
- `deleting` / `deleted`

### 📦 使用方式

#### 定義 Observer

```go
type UserObserver struct{}

func (UserObserver) Creating(model interface{}) error {
	fmt.Println("Creating:", model)
	return nil
}

func (UserObserver) Created(model interface{}) error {
	fmt.Println("Created:", model)
	return nil
}
```

#### 模型自註冊（推薦）

模型可以實作 `ObservedModel` 介面，自動綁定對應的 observer：

```go
func (u User) Observers() []odm.ModelObserver {
	return []odm.ModelObserver{UserObserver{}}
}
```

這樣在呼叫 `user.Create()` 時會自動觸發 observer。

### 🌐 全域 Observer

可以全域註冊 Observer，對所有模型生效：

```go
odm.RegisterGlobalObserver(AuditObserver{})
```

### 🎛️ Observer 擴充功能

#### ✅ 事件過濾（只觀察某些事件）

實作 `EventFilter` 介面：

```go
func (o UserObserver) InterestedIn(stage string) bool {
	return stage == "creating" || stage == "deleted"
}
```

#### ✅ 模型過濾（只監聽某些模型）

實作 `TypedObserver` 介面：

```go
func (o UserObserver) Accepts(model interface{}) bool {
	_, ok := model.(*User)
	return ok
}
```

#### ✅ 優先順序（Priority）

實作 `PrioritizedObserver`，可控制執行順序：

```go
func (o UserObserver) Priority() int {
	return 100 // 數字越大越早執行
}
```

#### ✅ 錯誤處理攔截

可設定全域錯誤攔截器：

```go
odm.RegisterObserverErrorHandler(func(err error, stage string, model interface{}) {
	log.Printf("[observer error] %s: %v", stage, err)
})
```

如果你有更多進階需求（例如事件佇列、非同步 observer），GODM 架構已支援進一步擴展。

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

#### 判斷指定目標是否存在

```go
exists, err := NewUser().
    WhereIn("name", []interface{}{"Alice", "Bob"}).
    Exists()
if err != nil {
	// error process
}
if (exists) {
	// code ...
}
```

## 🔗 關聯查詢（with 關聯預載入）

GODM 支援關聯查詢，可以方便地預載入關聯模型的資料。以下是如何使用 `with` 方法來查詢關聯資料的示例。

### User → Posts

假設一個 `User` 有多個 `Post`，可以這樣查詢：

```go
var users []User
err := NewUser().
    With("posts").
    All(&users)

for _, user := range users {
    fmt.Println("User:", user.Name)
    for _, post := range user.Posts {
        fmt.Println("Post:", post.Title)
    }
}
```

### Post → User

如果要查詢 `Post` 以及其對應的 `User`，可以這樣做：

```go
var posts []Post
err := NewPost().
    With("user").
    All(&posts)

for _, post := range posts {
    fmt.Println("Post:", post.Title)
    fmt.Println("User:", post.User.Name)
}
```

---

## 💡 靈感來源

GODM 的設計靈感來自於 [Laravel Eloquent ORM](https://laravel.com/docs/eloquent)，試圖為 Golang 帶來一種熟悉且簡潔的資料查詢體驗。它並非 ORM，而是專注於查詢構建、結果解碼與事務包裝，適合喜歡鏈式語法與輕量抽象的使用者。

## 📂 專案結構

```
godm/
├── examples/                  		# 使用範例
│	└── model/                  	# 自訂 User / Post 模型
│		├── post.go
│		└── user.go
│   ├── example.go
│   └── observer.go
│   └── relation.go
├── pkg/
│   └── odm/                   		# GODM 核心實作
│       ├── aggregate.go       		# MongoDB 聚合操作輔助工具
│       ├── config.go          		# 組態與全域資料庫客戶端設定
│       ├── context.go         		# Context 處理（自定義 context 注入）
│       ├── crud.go            		# CRUD 方法：建立、更新、刪除等
│       ├── model.go           		# GODM 結構定義與鏈式操作 API
│       ├── operator.go        		# MongoDB 運算子與對應處理
│       ├── query.go           		# 查詢構建邏輯（where, orWhere, select 等）
│       ├── relation.go        		# 關聯預載 (With, SetRelationConfig 等)
│       ├── transaction.go     		# 使用 MongoDB session 的交易包裝
│       ├── util.go            		# 工具函式（如 ObjectID 處理）
│       ├── observer.go        		# Observer 介面與註冊邏輯
│       └── observer_dispatch.go 	# Observer 執行與分派邏輯
├── go.mod                     		# Go module 定義檔
└── README.md                  		# 本說明文件
└── Changelog.md              		# 版本紀錄
```

---

## 📄 授權

本專案採用 [MIT License](./License) 授權。
