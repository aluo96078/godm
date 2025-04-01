[English](./Readme.en.md) | [繁體中文](../Readme.md)

# GODM: A Lightweight Query Mapper for MongoDB in Go

## 📚 Table of Contents
- [🧩 Introduction](#-introduction)
- [✨ Features](#-features)
- [👀 Observer Mechanism (Model Listening)](#-observer-mechanism-model-listening)
  - [🎯 Supported Events](#-supported-events)
  - [📦 Usage](#-usage)
    - [Defining Observer](#defining-observer)
    - [Model Self-Registration (Recommended)](#model-self-registration-recommended)
  - [🌐 Global Observer](#-global-observer)
  - [🎛️ Observer Extensions](#-observer-extensions)
    - [✅ Event Filtering](#event-filtering-only-observe-some-events)
    - [✅ Model Filtering](#model-filtering-only-listen-to-some-models)
    - [✅ Priority](#priority)
    - [✅ Error Handling Interception](#error-handling-interception)
- [🛠 Usage (with User model)](#-usage-with-user-model)
  - [Method Overriding (Return Custom Type)](#method-overriding-return-custom-type)
  - [Create and Query](#create-and-query)
  - [Aggregation and Transaction Operations](#aggregation-and-transaction-operations)
  - [More Query Examples](#more-query-examples)
    - [Using WhereID](#using-whereid)
    - [Using OR Queries](#using-or-queries)
    - [Using WhereIn and Field Selection](#using-wherein-and-field-selection)
    - [Using Pagination and Sorting](#using-pagination-and-sorting)
    - [Using Custom Context (with Timeout)](#using-custom-context-with-timeout)
    - [Check If Target Exists](#check-if-target-exists)
- [🔗 Association Queries (with Eager Loading)](#-association-queries-with-eager-loading)
  - [User → Posts](#user--posts)
  - [Post → User](#post--user)
- [💡 Inspiration](#-inspiration)
- [📂 Project Structure](#-project-structure)
- [📄 License](#-license)

## 🧩 Introduction

GODM (Go Object-Document Mapper) is a lightweight query wrapper tool for MongoDB, implemented in Go. It provides an ORM-like development experience and simplifies common query conditions and chain operations, helping you quickly build data models and perform CRUD, aggregation, transaction, and other operations.

The core implementation is located in [`pkg/odm`](./pkg/odm), and usage examples can be found in [`examples/`](./examples).

---

## ✨ Features

- 🚀 Chain query syntax (Where, OrWhere, WhereIn, etc.)
- 🔧 Automatic association of data models and collections (supports custom collection names and database names)
- 💾 Supports CRUD and BulkCreate
- 🧠 Supports complex query condition combinations (AND / OR)
- 🔁 Supports MongoDB aggregation pipelines
- 🔗 Supports eager loading of associated data with `with`
- 💼 Built-in transaction wrapper `WithTransaction`
- 👀 Built-in Observer mechanism, supporting model-level, global, sorting, and filtering (Inspired by Laravel)
- 🧪 Concise and easy to test, modular design for easy extension

## 👀 Observer Mechanism (Model Listening)

GODM has a built-in Observer system similar to Laravel Eloquent, allowing you to automatically trigger corresponding logic before and after the model's `Create`, `Update`, and `Delete` operations, suitable for scenarios like data validation, logging, and event tracking.

### 🎯 Supported Events

- `creating` / `created`
- `updating` / `updated`
- `deleting` / `deleted`

### 📦 Usage

#### Defining Observer

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

Models can implement the `ObservedModel` interface to automatically bind to the corresponding observer:

```go
func (u User) Observers() []odm.ModelObserver {
    return []odm.ModelObserver{UserObserver{}}
}
```

This way, calling `user.Create()` will automatically trigger the observer.

### 🌐 Global Observer

You can globally register an Observer that applies to all models:

```go
odm.RegisterGlobalObserver(AuditObserver{})
```

### 🎛️ Observer Extensions

#### ✅ Event Filtering (Only Observe Some Events)

Implement the `EventFilter` interface:

```go
func (o UserObserver) InterestedIn(stage string) bool {
    return stage == "creating" || stage == "deleted"
}
```

#### ✅ Model Filtering (Only Listen to Some Models)

Implement the `TypedObserver` interface:

```go
func (o UserObserver) Accepts(model interface{}) bool {
    _, ok := model.(*User)
    return ok
}
```

#### ✅ Priority

Implement `PrioritizedObserver` to control execution order:

```go
func (o UserObserver) Priority() int {
    return 100 // The larger the number, the earlier it executes
}
```

#### ✅ Error Handling Interception

You can set a global error interceptor:

```go
odm.RegisterObserverErrorHandler(func(err error, stage string, model interface{}) {
    log.Printf("[observer error] %s: %v", stage, err)
})
```

If you have more advanced requirements (such as event queues, asynchronous observers), the architecture of GODM supports further extensions.

## 🛠 Usage (with User model)

### Method Overriding (Return Custom Type)

GODM methods default to returning `*GODM`, but if you want to retain a custom model type (e.g., `*User`) so that you can access fields during chain operations, you can override the corresponding method in the model, for example:

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

### Create and Query

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

#### Using OR Queries

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

#### Check If Target Exists

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

## 🔗 Association Queries (with Eager Loading)

GODM supports association queries, allowing you to conveniently preload associated model data. Here’s how to use the `with` method to query associated data.

### User → Posts

Assuming a `User` has multiple `Post`, you can query like this:

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

If you want to query `Post` and its corresponding `User`, you can do this:

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

## 💡 Inspiration

The design of GODM is inspired by [Laravel Eloquent ORM](https://laravel.com/docs/eloquent), aiming to bring a familiar and concise data query experience to Golang. It is not an ORM but focuses on query building, result decoding, and transaction wrapping, suitable for users who enjoy chain syntax and lightweight abstraction.

## 📂 Project Structure

```
godm/
├── examples/                  		# Usage examples
│	└── model/                  	# Custom User / Post models
│		├── post.go
│		└── user.go
│   ├── example.go
│   └── observer.go
│   └── relation.go
├── pkg/
│   └── odm/                   		# Core implementation of GODM
│       ├── aggregate.go       		# MongoDB aggregation operation helpers
│       ├── config.go          		# Configuration and global database client settings
│       ├── context.go         		# Context handling (custom context injection)
│       ├── crud.go            		# CRUD methods: create, update, delete, etc.
│       ├── model.go           		# GODM structure definitions and chain operation API
│       ├── operator.go        		# MongoDB operators and corresponding handling
│       ├── query.go           		# Query building logic (where, orWhere, select, etc.)
│       ├── relation.go        		# Association preloading (With, SetRelationConfig, etc.)
│       ├── transaction.go     		# Transaction wrapping using MongoDB sessions
│       ├── util.go            		# Utility functions (e.g., ObjectID handling)
│       ├── observer.go        		# Observer interface and registration logic
│       └── observer_dispatch.go 	# Observer execution and dispatch logic
├── go.mod                     		# Go module definition file
└── README.md                  		# This documentation file
└── Changelog.md              		# Version history
```

---

## 📄 License

This project is licensed under the [MIT License](../License).
