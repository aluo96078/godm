[English](./Readme.en.md) | [繁體中文](../Readme.md)

# GODM: A Simple Query Mapper for MongoDB in Go

## 📚 Table of Contents

- [🧩 Introduction](#🧩-Introduction)
- [✨ Features](#✨-Features)
- [🛠 Usage (Example with User Model)](#🛠-Usage-Example-with-User-Model)
  - [Define Model](#Define-Model)
  - [Load Database Connection](#Load-Database-Connection)
  - [Specify Model to Use Another Database](#Specify-Model-to-Use-Another-Database)
  - [Method Overriding (Return Custom Type)](#Method-Overriding-Return-Custom-Type)
  - [Create and Query](#Create-and-Query)
  - [Aggregation and Transaction Operations](#Aggregation-and-Transaction-Operations)
  - [More Query Examples](#More-Query-Examples)
    - [Using WhereID](#Using-WhereID)
    - [Using OR Query](#Using-OR-Query)
    - [Using WhereIn and Field Selection](#Using-WhereIn-and-Field-Selection)
    - [Using Pagination and Sorting](#Using-Pagination-and-Sorting)
    - [Using Custom Context (Including Timeout)](#Using-Custom-Context-Including-Timeout)
    - [Check if a Target Exists](#Check-if-a-Target-Exists)
- [🔗 Relationship Queries (with Preloading)](#🔗-Relationship-Queries-with-Preloading)
  - [Model Definition](#Model-Definition)
  - [Relationship Settings](#Relationship-Settings)
  - [User → Posts](#User--Posts)
  - [Post → User](#Post--User)
- [👀 Observer Mechanism (Model Listening)](#👀-Observer-Mechanism-Model-Listening)
  - [Supported Events](#Supported-Events)
  - [Usage](#Usage-1)
    - [Define Observer](#Define-Observer)
    - [Model Self-Registration (Recommended)](#Model-Self-Registration-Recommendation)
    - [Global Observer](#Global-Observer)
    - [Observer Extension Features](#Observer-Extension-Features)
      - [Event Filtering](#Event-Filtering-Only-Observe-Certain-Events)
      - [Model Filtering](#Model-Filtering-Only-Monitor-Certain-Models)
      - [Priority](#Priority)
      - [Error Handling Interception](#Error-Handling-Interception)
- [💡 Inspiration](#💡-Inspiration)
- [🖥 System Architecture](#🖥-System-Architecture)
- [📂 Project Structure](#📂-Project-Structure)
- [📝 Usage Notes and Extensions](#📝-Usage-Notes-and-Extensions)
- [📄 License](#📄-License)


## 🧩 Introduction

GODM (Go Object-Document Mapper) is a lightweight query wrapper tool for MongoDB, implemented in Go. It provides an ORM-like development experience and simplifies common query conditions and chained operations, helping you quickly build data models and perform CRUD, aggregation, and transaction operations.

The core implementation is located in [`pkg/odm`](./pkg/odm), and usage examples can be found in [`examples/`](./examples).

---

## ✨ Features

- 🚀 Chained query syntax (Where, OrWhere, WhereIn, etc.)
- 🔧 Automatic association of data models and collections (supports custom collection names and database names)
- 💾 Supports CRUD and BulkCreate
- 🧠 Supports complex query condition combinations (AND / OR)
- 🔁 Supports MongoDB aggregation pipeline
- 🔗 Supports eager loading of related data (with)
- 💼 Built-in transaction wrapper `WithTransaction`
- 👀 Built-in Observer mechanism, supporting model-level, global, sorting, and filtering (Inspired by Laravel)
- 🧪 Simple and testable, modular design for easy extension

## 🛠 Usage (Example with User Model)

### Define Model
```go
// User defines the user model, embedding ODM and including user-specific fields.
type User struct {
    odm.GODM `bson:"-"` // Used to inherit GODM-related properties
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
    Email    string             `bson:"email"`
    // Model required for With relationship
    Posts []Post `bson:"posts,omitempty"`
}

// NewUser creates a new User instance and initializes ODM.
func NewUser() *User {
    u := &User{}
    // Complete GODM initialization
    u.Use(u)
    return u
}
```

### Load Database Connection
```go
// Replace with database connection information from configuration
client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:1145141919810@localhost:27017"))
if err != nil {
    log.Fatal(err)
}
// Specify global MongoClient
odm.MongoClient = client
// Specify global database name
odm.DBName = "test"

// You must set odm.DBName and odm.MongoClient first, otherwise an error will occur
u := NewUser()
```

### Specify Model to Use Another Database
```go
// Database connection settings...
u := NewUser()
u.SetDBName("db_name")
```

### Method Overriding (Return Custom Type)

GODM methods return `*GODM` by default, but if you want to retain the custom model type (e.g., `*User`) to access fields during chained operations, you can override the corresponding method in the model, for example:

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

#### Using Custom Context (Including Timeout)

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user := NewUser().WithContext(ctx)
_ = user.Where("email", "=", "timeout@example.com").First()
```

#### Check if a Target Exists

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

#### With Relationship Preloading

##### Model Definition

```go
// user.go
type User struct {
    odm.GODM `bson:"-"` // Used to inherit GODM-related properties
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Name     string             `bson:"name"`
    Email    string             `bson:"email"`
    // Model required for With relationship
    Posts []Post `bson:"posts,omitempty"`
}
```

```go
// post.go

type Post struct {
    odm.GODM `bson:"-"`
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    UserID   primitive.ObjectID `bson:"user_id"`
    Title    string             `bson:"title"`
    Body     string             `bson:"body"`
    // Model required for With relationship
    User *User `bson:"user,omitempty"`
}
```
#### Relationship Settings
```go
// NewUserModel creates a User model with a one-to-many with relationship setting
func NewUserModel() *odm.GODM {
    user := NewUser() // Reference the NewUser() above
    return user.SetRelationConfig(map[string]odm.RelationConfig{
        "posts": {
            // Target table name for the relationship
            From:         "posts",
            // Local field for the relationship
            LocalField:   "_id",
            // Foreign key in the external table
            ForeignField: "user_id",
            // Name of the related data field
            As:           "posts",
            // Whether there are multiple records
            IsArray:      true,
        },
    })
}
// NewPostModel creates a Post model with a one-to-one with relationship setting
func NewPostModel() *odm.GODM {
    post := examples.NewPost()
    return post.SetRelationConfig(map[string]odm.RelationConfig{
        "user": {
            // Target table name for the relationship
            From:         "users",
            // Local field for the relationship
            LocalField:   "user_id",
            // Foreign key in the external table
            ForeignField: "_id",
            // Name of the related data field
            As:           "user",
            // Whether there are multiple records
            IsArray:      false,
        },
    })
}
```

##### User → Posts

Assuming a `User` has multiple `Post`, you can query like this:

```go
var users []User
err := NewUserModel().
    With("posts").
    All(&users)

for _, user := range users {
    fmt.Println("User:", user.Name)
    for _, post := range user.Posts {
        fmt.Println("Post:", post.Title)
    }
}
```

##### Post → User

If you want to query `Post` and its corresponding `User`, you can do this:

```go
var posts []Post
err := NewPostModel().
    With("user").
    All(&posts)

for _, post := range posts {
    fmt.Println("Post:", post.Title)
    fmt.Println("User:", post.User.Name)
}
```

## 👀 Observer Mechanism (Model Listening)

GODM has a built-in Observer system similar to Laravel Eloquent, allowing you to automatically trigger corresponding logic before and after model operations such as `Create`, `Update`, and `Delete`, making it suitable for data validation, logging, event tracking, and other scenarios.

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

This way, when calling `user.Create()`, the observer will be automatically triggered.

### 🌐 Global Observer

You can register a global Observer that applies to all models:

```go
odm.RegisterGlobalObserver(AuditObserver{})
```

### 🎛️ Observer Extension Features

#### ✅ Event Filtering (Only Observe Certain Events)

Implement the `EventFilter` interface:

```go
func (o UserObserver) InterestedIn(stage string) bool {
    return stage == "creating" || stage == "deleted"
}
```

#### ✅ Model Filtering (Only Monitor Certain Models)

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

If you have more advanced needs (such as event queues, asynchronous observers), the GODM architecture already supports further extensions.

---

## 💡 Inspiration

The design inspiration for GODM comes from [Laravel Eloquent ORM](https://laravel.com/docs/eloquent), aiming to bring a familiar and concise data querying experience to Golang. It is not an ORM but focuses on query construction, result decoding, and transaction wrapping, suitable for users who prefer chained syntax and lightweight abstraction.

## 🖥 System Architecture

The core modules of GODM include:

- Query Builder: Responsible for handling chained query conditions
- CRUD Operations: Handles create, query, update, and delete
- Aggregation Operations: Supports MongoDB aggregation pipeline
- Transaction Processing: Wraps MongoDB session transactions
- Observer Module: Responsible for listening to and dispatching events before and after model operations
- With Relationship Preloading: Loads related data through a single query

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
│       ├── aggregate.go       		# MongoDB aggregation operation helper
│       ├── config.go          		# Configuration and global database client settings
│       ├── context.go         		# Context handling (custom context injection)
│       ├── crud.go            		# CRUD methods: create, update, delete, etc.
│       ├── model.go           		# GODM structure definitions and chained operation APIs
│       ├── operator.go        		# MongoDB operators and corresponding handling
│       ├── query.go           		# Query construction logic (where, orWhere, select, etc.)
│       ├── relation.go        		# Relationship preloading (With, SetRelationConfig, etc.)
│       ├── transaction.go     		# Transaction wrapping using MongoDB session
│       ├── util.go            		# Utility functions (such as ObjectID handling)
│       ├── observer.go        		# Observer interface and registration logic
│       └── observer_dispatch.go 	# Observer execution and dispatch logic
├── go.mod                     		# Go module definition file
└── README.md                  		# This documentation file
└── Changelog.md              		# Version history
```

---

## 📝 Usage Notes and Extensions

Notes:

- GODM is a lightweight query wrapper tool, not a complete ORM.
- Focused on query construction, result decoding, and transaction wrapping.
- Designed to be modular for easy testing and extension, allowing for the addition of event queues or asynchronous processing as needed.

## 📄 License

This project is licensed under the [MIT License](./License).
