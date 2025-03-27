<file name=0 path=/Users/aluo/project/gotest/Readme.md><p align="right">
  🌐 [English](docs/README.en.md) | [繁體中文](README.md)
</p>
# MongoDB GODM 使用說明（以 User 模型為範例）

本專案提供一個輕量級的 MongoDB GODM（Object-Document Mapper），封裝在 `odm/godm.go` 中。使用者可透過內嵌該 GODM 結構來快速為任意資料模型賦予 MongoDB 查詢與操作能力。

本說明將透過 `examples/user.go` 所定義的 `User` 模型作為範例，示範如何整合與使用此 GODM。

---

## 2. 以 User 模型為範例

User 模型定義在 `user.go` 中，結構如下：

```go
type User struct {
    GODM   `bson:"-"`
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string             `bson:"name"`
    Email string             `bson:"email"`
}
```

User 模型內嵌 GODM，因此能直接調用 GODM 的所有方法，並提供覆蓋方法如 `SetCollectionName` 與 `WithContext`，方便直接操作使用者專屬欄位。

---

## 3. 主要用例示例

### 3.1 初始化 GODM 與 User 模型

```go
// 建立一個新的 User 實例
user := NewUser()
// 可選：指定資料庫名稱 (預設為 "your_db")
user.DbName = "my_database"
// 初始化 GODM，依據模型自動決定集合名稱（預設為 "users"）
user.Use(user)
```

### 3.2 自定義上下文

```go
// 設定一個帶有超時的上下文
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 使用 WithContext 設定自定義上下文
user.WithContext(ctx)
```

### 3.3 自定義集合名稱

```go
// 透過 SetCollectionName 覆蓋預設集合名稱，並回傳 *User
customUser := NewUser().SetCollectionName("custom_users")
```

### 3.4 建立 (Create) 操作

```go
user.Name = "Test User"
user.Email = "test@example.com"
err := user.Create()
if err != nil {
    fmt.Println("Create error:", err)
} else {
    fmt.Println("使用者已建立 (User created)")
}
```

### 3.5 查詢第一筆資料 (First)

```go
// 根據 email 查詢第一筆符合條件的使用者
err = user.Where("email", "=", "test@example.com").First()
if err != nil {
    fmt.Println("找不到使用者 (User not found)")
} else {
    fmt.Printf("找到使用者: %s\n", user.Name)
}
```

### 3.6 更新 (Update) 操作

```go
// 更新符合條件的使用者名稱
err = user.Where("email", "=", "test@example.com").Update(bson.M{"name": "Updated User"})
if err != nil {
    fmt.Println("Update error:", err)
} else {
    fmt.Println("使用者已更新 (User updated)")
}
```

### 3.7 批量建立 (BulkCreate)

```go
user1 := NewUser(); user1.Name = "Alice"; user1.Email = "alice@example.com"
user2 := NewUser(); user2.Name = "Bob"; user2.Email = "bob@example.com"
user3 := NewUser(); user3.Name = "Charlie"; user3.Email = "charlie@example.com"

// 將多個使用者打包後一次插入
bulkUsers := []interface{}{user1, user2, user3}
err = user.BulkCreate(bulkUsers)
if err != nil {
    fmt.Println("BulkCreate error:", err)
} else {
    fmt.Println("批量使用者已建立 (Bulk users created)")
}
```

### 3.8 查詢過濾器示例

```go
// AND 條件：查詢 email 為 "test@example.com" 的使用者
user.Where("email", "=", "test@example.com")

// AND 條件 (包含)：查詢 name 為 "Alice" 或 "Bob"
user.WhereIn("name", []interface{}{ "Alice", "Bob" })

// OR 條件：查詢 email 為 "alice@example.com" 或 name 為 "Charlie"
user.OrWhere("email", "=", "alice@example.com").
     OrWhere("name", "=", "Charlie")
```

### 3.9 聚合 (Aggregate) 操作

```go
// 定義一個聚合管道，過濾 email 不為空並按 name 分組計數
pipeline := mongo.Pipeline{
    {{"$match", bson.M{"email": bson.M{"$ne": ""}}}},
    {{"$group", bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}}},
}
var aggregateResults []bson.M
err = user.Aggregate(pipeline, &aggregateResults)
if err != nil {
    fmt.Println("Aggregate error:", err)
} else {
    for _, res := range aggregateResults {
        fmt.Printf("Name: %v, Count: %v\n", res["_id"], res["count"])
    }
}
```

### 3.10 事務 (Transaction) 操作

```go
// 在事務中執行建立及更新操作
err = user.WithTransaction(func(sessCtx mongo.SessionContext) error {
    if err := user.Create(); err != nil {
        return err
    }
    if err := user.Where("email", "=", "test@example.com").Update(bson.M{"name": "Tx Updated User"}); err != nil {
        return err
    }
    return nil
})
if err != nil {
    fmt.Println("WithTransaction error:", err)
} else {
    fmt.Println("事務操作成功 (Transaction executed successfully)")
}
```

---

## 結論

透過 GODM 與 User 模型，您可以使用統一且簡單的 API 與 MongoDB 進行各種文件操作。GODM 抽象化了 MongoDB 的驅動，並提供了豐富的查詢、更新、聚合以及事務支持，使得開發者能夠更專注於業務邏輯。
</file>
