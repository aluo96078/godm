package examples

import (
	"fmt"
	examples "godm/examples/model"
	"godm/pkg/odm"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewUserModel 建立具備 with 關聯設定的 User 模型
func NewUserModel() *odm.GODM {
	user := examples.NewUser()
	return user.SetRelationConfig(map[string]odm.RelationConfig{
		"posts": {
			From:         "posts",
			LocalField:   "_id",
			ForeignField: "user_id",
			As:           "posts",
			IsArray:      true,
		},
	})
}

func NewPostModel() *odm.GODM {
	post := examples.NewPost()
	return post.SetRelationConfig(map[string]odm.RelationConfig{
		"user": {
			From:         "users",
			LocalField:   "user_id",
			ForeignField: "_id",
			As:           "user",
			IsArray:      false,
		},
	})
}

func RelationExample() {
	// 初始化 ODM 用戶模型
	userModel := NewUserModel()

	// 建立測試資料
	userID := primitive.NewObjectID()
	post1 := examples.Post{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Title:  "Post 1",
		Body:   "This is the first post.",
	}
	post2 := examples.Post{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Title:  "Post 2",
		Body:   "This is the second post.",
	}

	// 插入使用者
	userModel.Model = &examples.User{
		ID:    userID,
		Name:  "With Tester",
		Email: "with@test.com",
	}
	err := userModel.Create()
	if err != nil {
		log.Println("插入使用者錯誤:", err)
		return
	}
	// 插入貼文
	postModel := NewPostModel()
	postModel.Model = &post1
	err = postModel.Create()
	if err != nil {
		log.Println("插入貼文錯誤:", err)
		return
	}

	postModel.Model = &post2
	err = postModel.Create()
	if err != nil {
		log.Println("插入貼文錯誤2:", err)
		return
	}
	// 查詢並預載入 posts
	var user examples.User
	err = userModel.WhereID(userID).With("posts").First()
	if err != nil {
		log.Println(userModel.ToBson())
		fmt.Println("查詢錯誤:", err)
		return
	}
	user = *(userModel.Model.(*examples.User)) // 轉型回 user

	fmt.Printf("使用者: %s\n", user.Name)
	for _, p := range user.Posts {
		fmt.Printf("  貼文: %s - %s\n", p.Title, p.Body)
	}

	// 查詢貼文並預載入使用者
	var post examples.Post
	err = postModel.WhereID(post1.ID).With("user").First()
	if err != nil {
		fmt.Println("查詢錯誤:", err)
		return
	}
	post = *(postModel.Model.(*examples.Post)) // 轉型回 post

	fmt.Printf("貼文: %s\n", post.Title)
	fmt.Printf("  使用者: %s - %s\n", post.User.Name, post.User.Email)
}
