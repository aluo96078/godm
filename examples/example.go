package examples

import (
	"context"
	"fmt"
	examples "godm/examples/model"
	"godm/pkg/odm"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Examples() {
	// 1. 連線至 MongoDB
	// 1. Connect to MongoDB
	// 資料庫連線資訊應該替換爲讀取環境變量
	// Database connection information should be replaced with environment variables
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://root:1145141919810@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	// 指定全域 MongoClient
	// Set the global MongoClient
	odm.MongoClient = client
	// 指定全域資料庫名稱
	// Set the global database name
	odm.DBName = "test"

	// 註冊全局觀察者
	// Register global observers
	odm.RegisterGlobalObserver(&UserObserver{})

	// -------------------------------------------------------
	// 2. Create 操作 - 建立一個新使用者
	// 2. Create operation - Create a new user
	fmt.Println("\n--- Create Operation ---")
	user := examples.NewUser()
	user.Name = "Test User"
	user.Email = "test@example.com"
	err = user.Create()
	if err != nil {
		fmt.Println("Create error:", err)
	} else {
		fmt.Println("使用者已建立 (User created)")
	}

	// -------------------------------------------------------
	// 3. First 操作 - 根據條件查詢第一個使用者
	// 3. First operation - Retrieve the first user matching the condition
	fmt.Println("\n--- First Operation ---")
	user = examples.NewUser()
	err = user.Where("email", "=", "test@example.com").First()
	if err != nil {
		fmt.Println("找不到使用者 (User not found)")
	} else {
		fmt.Printf("找到使用者: %s\n", user.Name)
	}

	// -------------------------------------------------------
	// 4. Update 操作 - 更新使用者的名稱
	// 4. Update operation - Update the user's name
	fmt.Println("\n--- Update Operation ---")
	err = user.Where("email", "=", "test@example.com").Update(bson.M{"name": "Updated User"})
	if err != nil {
		fmt.Println("Update error:", err)
	} else {
		fmt.Println("使用者已更新 (User updated)")
	}

	// -------------------------------------------------------
	// 5. BulkCreate 操作 - 批量建立多個使用者
	// 5. BulkCreate operation - Bulk insert multiple users
	fmt.Println("\n--- BulkCreate Operation ---")
	// 建立多個使用者資料
	user1 := examples.NewUser()
	user1.Name = "Alice"
	user1.Email = "alice@example.com"

	user2 := examples.NewUser()
	user2.Name = "Bob"
	user2.Email = "bob@example.com"

	user3 := examples.NewUser()
	user3.Name = "Charlie"
	user3.Email = "charlie@example.com"

	// 使用 BulkCreate 插入
	bulkUsers := []interface{}{user1, user2, user3}
	err = user.BulkCreate(bulkUsers)
	if err != nil {
		fmt.Println("BulkCreate error:", err)
	} else {
		fmt.Println("批量使用者已建立 (Bulk users created)")
	}

	// -------------------------------------------------------
	// 6. WhereIn 操作 - 根據包含條件查詢使用者
	// 6. WhereIn operation - Query users with an inclusion condition
	fmt.Println("\n--- WhereIn Operation ---")
	var usersIn []examples.User
	// 查詢 email 為 alice@example.com 或 bob@example.com 的使用者
	user = examples.NewUser()
	err = user.WhereIn("email", []interface{}{"alice@example.com", "bob@example.com"}).All(&usersIn)
	if err != nil {
		fmt.Println("WhereIn error:", err)
	} else {
		fmt.Println("使用者 (WhereIn 查詢):")
		for _, u := range usersIn {
			fmt.Printf("- %s (%s)\n", u.Name, u.Email)
		}
	}

	// -------------------------------------------------------
	// 7. OrWhere 操作 - 使用 OR 條件查詢使用者
	// 7. OrWhere operation - Query users with OR conditions
	fmt.Println("\n--- OrWhere Operation ---")
	var usersOr []examples.User
	// 查詢 email 為 alice@example.com 或 name 為 "Charlie" 的使用者
	user = examples.NewUser()
	err = user.OrWhere("email", "=", "alice@example.com").OrWhere("name", "=", "Charlie").All(&usersOr)
	if err != nil {
		fmt.Println("OrWhere error:", err)
	} else {
		fmt.Println("使用者 (OrWhere 查詢):")
		for _, u := range usersOr {
			fmt.Printf("- %s (%s)\n", u.Name, u.Email)
		}
	}

	// -------------------------------------------------------
	// 8. Offset, OrderBy, Select 與 Exclude 操作 - 查詢並排序、分頁、選擇及排除特定欄位
	// 8. Offset, OrderBy, Select, and Exclude operations - Query with sorting, pagination, selecting, and excluding specific fields
	fmt.Println("\n--- Offset, OrderBy, Select & Exclude Operation ---")
	var usersPaginated []examples.User
	user = examples.NewUser()
	// 設定排序（依名稱升序）、限制筆數、跳過第一筆，並選擇只包含 name 欄位，排除 _id 欄位
	err = user.Select("name").Exclude("_id").OrderBy("name", true).Limit(2).Offset(1).All(&usersPaginated)
	if err != nil {
		fmt.Println("Paginated query error:", err)
	} else {
		fmt.Println("分頁查詢結果:")
		for _, u := range usersPaginated {
			fmt.Printf("- %s (Email: %s)\n", u.Name, u.Email)
		}
	}

	// -------------------------------------------------------
	// 9. Count 操作 - 統計符合條件的使用者數量
	// 9. Count operation - Count the number of users matching the condition
	fmt.Println("\n--- Count Operation ---")
	user = examples.NewUser()
	count, err := user.Where("email", "=", "alice@example.com").Count()
	if err != nil {
		fmt.Println("Count error:", err)
	} else {
		fmt.Println("符合條件的使用者數量 (Count):", count)
	}

	// -------------------------------------------------------
	// 10. Aggregate 操作 - 執行聚合管道，例如按名稱分組並計算數量
	// 10. Aggregate operation - Run an aggregation pipeline, e.g., group by name and count
	fmt.Println("\n--- Aggregate Operation ---")
	// 定義聚合管道: 過濾 email 不為空，並根據 name 分組計數
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"email": bson.M{"$ne": ""}}}},
		{{Key: "$group", Value: bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}}},
	}
	var aggregateResults []bson.M
	user = examples.NewUser()
	err = user.Aggregate(pipeline, &aggregateResults)
	if err != nil {
		fmt.Println("Aggregate error:", err)
	} else {
		fmt.Println("Aggregate 結果:")
		for _, res := range aggregateResults {
			fmt.Printf("Name: %v, Count: %v\n", res["_id"], res["count"])
		}
	}

	// -------------------------------------------------------
	// 11. WithTransaction 操作 - 在事務中執行多個操作 (建立並更新使用者)
	// 11. WithTransaction operation - Execute multiple operations in a transaction (create and update a user)
	fmt.Println("\n--- WithTransaction Operation ---")
	userTx := examples.NewUser()
	userTx.Name = "Transaction User"
	userTx.Email = "txuser@example.com"

	err = userTx.WithTransaction(func(sessCtx mongo.SessionContext) error {
		// 在交易中建立使用者
		if err := userTx.Create(); err != nil {
			return err
		}
		// 在交易中更新使用者名稱
		if err := userTx.Where("email", "=", "txuser@example.com").Update(bson.M{"name": "Tx Updated User"}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("WithTransaction error:", err)
	} else {
		fmt.Println("事務操作成功 (Transaction executed successfully)")
	}

	// -------------------------------------------------------
	// 12. WithContext 操作 - 使用自定義上下文，例如設定超時
	// 12. WithContext operation - Use a custom context, e.g., with a timeout
	fmt.Println("\n--- WithContext Operation ---")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userCtx := examples.NewUser().WithContext(ctx)
	userCtx.Name = "Context User"
	userCtx.Email = "ctxuser@example.com"
	err = userCtx.Create()
	if err != nil {
		fmt.Println("WithContext Create error:", err)
	} else {
		fmt.Println("使用自定義上下文建立使用者 (User created with custom context)")
	}

	// -------------------------------------------------------
	// 13. SetCollectionName 操作 - 使用自定義集合名稱
	// 13. SetCollectionName operation - Use a custom collection name
	fmt.Println("\n--- SetCollectionName Operation ---")
	customUser := examples.NewUser().SetCollectionName("custom_users")
	customUser.Name = "Custom User"
	customUser.Email = "custom@example.com"
	err = customUser.Create()
	if err != nil {
		fmt.Println("SetCollectionName Create error:", err)
	} else {
		fmt.Println("使用自定義集合建立使用者 (User created in custom collection)")
	}

	// 注意: 為保持示例簡單，本例未執行清理操作 (例如刪除測試數據)。
	// Note: For simplicity, this example does not perform cleanup of test data.

	// 14. Relation Example - 關聯操作 (多對多關聯)
	RelationExample()
}
