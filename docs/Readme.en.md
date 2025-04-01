[English](./Readme.en.md) | [繁體中文](../Readme.md)

# GODM: A Simple Query Mapper for MongoDB in Go

## 🧩 Introduction

GODM (Go Object-Document Mapper) is a lightweight query encapsulation tool for MongoDB, implemented in the Go language. It provides an ORM-like development experience and simplifies common query conditions and chain operations, helping you quickly build data models and perform CRUD, aggregation, transaction, and other operations.

The core implementation is located in [`pkg/odm`](./pkg/odm), and usage examples can be found in [`examples/`](./examples).

---

## ✨ Features

- 🚀 Chain query syntax (Where, OrWhere, WhereIn, etc.)
- 🔧 Automatic association of data models and collections (supports custom collection names and database names)
- 💾 Supports CRUD and BulkCreate
- 🧠 Supports complex query condition combinations (AND / OR)
- 🔁 Supports MongoDB aggregation pipelines
- 💼 Built-in transaction encapsulation `WithTransaction`
- 👀 Built-in Observer mechanism, supporting model-level, global, sorting, and filtering (Inspired by Laravel)
- 🧪 Simple and easy to test, modular design facilitates expansion

## 👀 Observer Mechanism

GODM has a built-in Observer system similar to Laravel Eloquent, allowing you to automatically trigger corresponding logic before and after the `Create`, `Update`, and `Delete` operations on models, suitable for scenarios such as data validation, logging, and event tracking.

### 🎯 Supported Events

- `creating` / `created`
- `updating` / `updated`
- `deleting` / `deleted`

### 📦 Usage

#### Define Observer

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

#### Model Self-Registration (Recommended)

Models can implement the `ObservedModel` interface to automatically bind the corresponding observer:

```go
func (u User) Observers() []odm.ModelObserver {
    return []odm.ModelObserver{UserObserver{}}
}
```

This way, when calling `user.Create()`, the observer will be triggered automatically.

### 🌐 Global Observer

You can register a global Observer that applies to all models:

```go
odm.RegisterGlobalObserver(AuditObserver{})
```

### 🎛️ Observer Extensions

#### ✅ Event Filtering (Observe Only Certain Events)

Implement the `EventFilter` interface:

```go
func (o UserObserver) InterestedIn(stage string) bool {
    return stage == "creating" || stage == "deleted"
}
```

#### ✅ Model Filtering (Listen Only to Certain Models)

Implement the `TypedObserver` interface:

```go
func (o UserObserver) Accepts(model interface{}) bool {
    _, ok := model.(*User)
    return ok
}
```

#### ✅ Priority

Implement `PrioritizedObserver` to control the execution order:

```go
func (o UserObserver) Priority() int {
    return 100 // Higher numbers execute earlier
}
```

#### ✅ Error Handling Interception

A global error handler can be set:

```go
odm.RegisterObserverErrorHandler(func(err error, stage string, model interface{}) {
    log.Printf("[observer error] %s: %v", stage, err)
})
```

If you have more advanced requirements (such as event queues, asynchronous observers), the GODM architecture supports further expansion.

## 🛠 Usage (Using User Model as an Example)

### Method Overriding (Return Custom Type)

GODM methods default to returning `*GODM`, but if you want to retain custom model types (e.g., `*User`) to access fields during chain operations, you can override the corresponding methods in your model, for example:

```go
func (o *User) SetCollectionName(name string) *User {
    o.CollectionName = name
    if o.Collection != nil {
        o.Collection = odm.MongoClient.Database(o.DbName).Collection(name)
    }
    return o
}
```

This way, type consistency can be maintained:

```go
u := NewUser().SetCollectionName("custom_users")
fmt.Println(u.Name) // Can directly use *User fields
```

### Create and Query

```go
user := NewUser()
user.Name = "Test"
user.Email = "test@example.com"
_ = user.Create()

// Query the first piece of data
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

#### Check if Specific Targets Exist

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

---

## 💡 Inspiration

The design of GODM is inspired by [Laravel Eloquent ORM](https://laravel.com/docs/eloquent), aiming to bring a familiar and concise data query experience to Golang. It is not an ORM but focuses on query building, result decoding, and transaction wrapping, suitable for users who enjoy chain syntax and lightweight abstraction.

## 📂 Project Structure

```
godm/
├── examples/                  # Usage examples, including entry points and custom User models
│   ├── example.go
│   └── user_observer.go
│   └── user.go
├── pkg/
│   └── odm/                   # Core implementation of GODM
│       ├── aggregate.go       # Helper for MongoDB aggregation operations
│       ├── config.go          # Configuration and global database client settings
│       ├── context.go         # Context handling (custom context injection)
│       ├── crud.go            # CRUD methods: create, update, delete, etc.
│       ├── model.go           # GODM structure definitions and chain operation API
│       ├── operator.go        # MongoDB operators and corresponding handling
│       ├── query.go           # Query building logic (where, orWhere, select, etc.)
│       ├── transaction.go     # Transaction wrapping using MongoDB sessions
│       ├── util.go            # Utility functions (e.g., ObjectID handling)
│       ├── observer.go        # Observer interface and registration logic
│       └── observer_dispatch.go # Observer execution and dispatch logic
├── go.mod                     # Go module definition file
└── README.md                  # This documentation file
└── Changelog.md               # Version history
```

---

## 📄 License

This project is licensed under the [MIT License](./License).
