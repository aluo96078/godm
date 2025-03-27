[English](./Readme.en.md) | [繁體中文](../Readme.md)

# GODM: A Simple Query Mapper for MongoDB in Go

## 🧩 Introduction

GODM (Go Object-Document Mapper) is a lightweight query encapsulation tool for MongoDB, implemented in Go. It provides an ORM-like development experience and simplifies common query conditions and chain operations, helping you quickly build data models and perform CRUD, aggregation, transactions, and other operations.

The core implementation is located in [`pkg/odm`](./pkg/odm), and usage examples can be found in [`examples/`](./examples).

---

## ✨ Features

- 🚀 Chain query syntax (Where, OrWhere, WhereIn, etc.)
- 🔧 Automatic association of data models and collections (supports custom collection names and database names)
- 💾 Supports CRUD and BulkCreate
- 🧠 Supports complex query condition combinations (AND / OR)
- 🔁 Supports MongoDB aggregation pipeline
- 💼 Built-in transaction wrapper `WithTransaction`
- 🧪 Simple and testable, modular design for easy extension

---

## 🛠 Usage (with User model example)

### Method Override (Return Custom Type)

GODM methods default to returning `*GODM`, but if you wish to retain a custom model type (e.g., `*User`) to access fields during chain operations, you can override the corresponding method in the model, for example:

```go
func (o *User) SetCollectionName(name string) *User {
    o.CollectionName = name
    if o.Collection != nil {
        o.Collection = odm.MongoClient.Database(o.DbName).Collection(name)
    }
    return o
}
```

This way, you can maintain type consistency:

```go
u := NewUser().SetCollectionName("custom_users")
fmt.Println(u.Name) // Can directly use *User fields
```

### Creating and Querying

```go
user := NewUser()
user.Name = "Test"
user.Email = "test@example.com"
_ = user.Create()

// Query the first record
err := user.Where("email", "=", "test@example.com").First()
```

### Aggregation and Transaction Operations

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

### More Query Examples

#### Using `WhereID`

```go
// Query document by MongoDB ObjectID
user := NewUser()
_ = user.WhereID("65f74c3a09c7a8f812345678").First()
```

#### Using OR Query

```go
var users []User
err := NewUser().
    OrWhere("name", "=", "Alice").
    OrWhere("email", "=", "bob@example.com").
    All(&users)
```

#### Using WhereIn and Field Selection

```go
var users []User
err := NewUser().
    WhereIn("name", []interface{}{"Alice", "Bob"}).
    Select("name").
    Exclude("email").
    All(&users)
```

#### Using Pagination and Sorting

```go
var users []User
err := NewUser().
    OrderBy("name", true).
    Offset(10).
    Limit(10).
    All(&users)
```

#### Using Custom Context (with Timeout)

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user := NewUser().WithContext(ctx)
_ = user.Where("email", "=", "timeout@example.com").First()
```

---

## 💡 Inspiration

The design of GODM is inspired by [Laravel Eloquent ORM](https://laravel.com/docs/eloquent), aiming to bring a familiar and concise data querying experience to Golang. It is not an ORM, but focuses on query building, result decoding, and transaction wrapping, suitable for users who enjoy chain syntax and lightweight abstraction.

## 📂 Project Structure

```
godm/
├── examples/        # Usage examples: main.go, user.go
├── pkg/
│   └── odm/         # Core implementation of GODM (modularized)
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

## 📄 License

This project is licensed under the [MIT License](./LICENSE).
