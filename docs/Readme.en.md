<p align="right">
  🌐 [English](docs/Readme.en.md) | [繁體中文](README.md)
</p>

# MongoDB GODM Usage Guide (Using User Model as an Example)

## 1. Introduction 

This project provides a lightweight MongoDB GODM (Object-Document Mapper), encapsulated in `odm/godm.go`. Users can quickly endow any data model with MongoDB querying and operation capabilities by embedding this GODM structure.

This guide will use the `User` model defined in `examples/user.go` as an example to demonstrate how to integrate and use this GODM.

---

## 2. Using User Model as an Example

The User model is defined in `user.go`, structured as follows:

```go
type User struct {
    GODM   `bson:"-"`
    ID    primitive.ObjectID `bson:"_id,omitempty"`
    Name  string             `bson:"name"`
    Email string             `bson:"email"`
}
```

The User model embeds GODM, allowing direct invocation of all methods of GODM and providing overridden methods such as `SetCollectionName` and `WithContext`, making it convenient to operate on user-specific fields.

---

## 3. Main Use Case Examples

### 3.1 Initialize GODM and User Model

```go
// Create a new User instance
user := NewUser()
// Optional: Specify the database name (default is "your_db")
user.DbName = "my_database"
// Initialize GODM, automatically determining the collection name based on the model (default is "users")
user.Use(user)
```

### 3.2 Custom Context

```go
// Set a context with a timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Use WithContext to set a custom context
user.WithContext(ctx)
```

### 3.3 Custom Collection Name

```go
// Override the default collection name using SetCollectionName and return *User
customUser := NewUser().SetCollectionName("custom_users")
```

### 3.4 Create Operation

```go
user.Name = "Test User"
user.Email = "test@example.com"
err := user.Create()
if err != nil {
    fmt.Println("Create error:", err)
} else {
    fmt.Println("User created")
}
```

### 3.5 Query First Record

```go
// Query the first user matching the email
err = user.Where("email", "=", "test@example.com").First()
if err != nil {
    fmt.Println("User not found")
} else {
    fmt.Printf("Found user: %s\n", user.Name)
}
```

### 3.6 Update Operation

```go
// Update the user's name that matches the condition
err = user.Where("email", "=", "test@example.com").Update(bson.M{"name": "Updated User"})
if err != nil {
    fmt.Println("Update error:", err)
} else {
    fmt.Println("User updated")
}
```

### 3.7 Bulk Create

```go
user1 := NewUser(); user1.Name = "Alice"; user1.Email = "alice@example.com"
user2 := NewUser(); user2.Name = "Bob"; user2.Email = "bob@example.com"
user3 := NewUser(); user3.Name = "Charlie"; user3.Email = "charlie@example.com"

// Package multiple users and insert them at once
bulkUsers := []interface{}{user1, user2, user3}
err = user.BulkCreate(bulkUsers)
if err != nil {
    fmt.Println("BulkCreate error:", err)
} else {
    fmt.Println("Bulk users created")
}
```

### 3.8 Query Filter Example

```go
// AND condition: Query user with email "test@example.com"
user.Where("email", "=", "test@example.com")

// AND condition (includes): Query user with name "Alice" or "Bob"
user.WhereIn("name", []interface{}{ "Alice", "Bob" })

// OR condition: Query user with email "alice@example.com" or name "Charlie"
user.OrWhere("email", "=", "alice@example.com").
     OrWhere("name", "=", "Charlie")
```

### 3.9 Aggregate Operation

```go
// Define an aggregation pipeline to filter users with non-empty email and group by name for counting
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

### 3.10 Transaction Operation

```go
// Execute create and update operations within a transaction
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
    fmt.Println("Transaction executed successfully")
}
```

---

## Conclusion

Through GODM and the User model, you can use a unified and simple API to perform various document operations with MongoDB. GODM abstracts the MongoDB driver and provides rich support for querying, updating, aggregating, and transactions, allowing developers to focus more on business logic.
