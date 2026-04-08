// 복잡한 예제: 사용자와 게시글 리소스, 관계 조회, 검색/페이지네이션, 에러 처리
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

// --- 모델 ---

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Post struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
}

// --- 인메모리 DB ---

type DB struct {
	mu     sync.RWMutex
	users  map[int]User
	posts  map[int]Post
	nextID int
}

func NewDB() *DB {
	return &DB{
		users:  make(map[int]User),
		posts:  make(map[int]Post),
		nextID: 1,
	}
}

func (db *DB) allocID() int {
	id := db.nextID
	db.nextID++
	return id
}

// --- User 리소스 (전체 CRUD) ---

type CreateUserReq struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,gte=0"`
}

type UpdateUserReq struct {
	Name *string `json:"name"`
	Age  *int    `json:"age" binding:"omitempty,gte=0"`
}

type UserResource struct {
	db *DB
}

// GET /users?name=... — 이름 검색 지원
func (r *UserResource) List(c *gin.Context) (any, int, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	nameFilter := c.Query("name")
	users := make([]User, 0)
	for _, u := range r.db.users {
		if nameFilter != "" && u.Name != nameFilter {
			continue
		}
		users = append(users, u)
	}
	return gin.H{"users": users, "count": len(users)}, http.StatusOK, nil
}

// GET /users/:id
func (r *UserResource) Get(id string, c *gin.Context) (any, int, error) {
	user, err := r.findUser(id)
	if err != nil {
		return nil, 0, err
	}
	return user, http.StatusOK, nil
}

// POST /users
func (r *UserResource) Post(c *gin.Context) (any, int, error) {
	body, err := restful.Bind[CreateUserReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	user := User{ID: r.db.allocID(), Name: body.Name, Age: body.Age}
	r.db.users[user.ID] = user
	return user, http.StatusCreated, nil
}

// PUT /users/:id — 전체 교체
func (r *UserResource) Put(id string, c *gin.Context) (any, int, error) {
	uid, err := r.parseID(id)
	if err != nil {
		return nil, 0, err
	}

	body, err := restful.Bind[CreateUserReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, ok := r.db.users[uid]; !ok {
		return nil, 0, restful.Abort(http.StatusNotFound, "user not found")
	}

	user := User{ID: uid, Name: body.Name, Age: body.Age}
	r.db.users[uid] = user
	return user, http.StatusOK, nil
}

// PATCH /users/:id — 부분 수정
func (r *UserResource) Patch(id string, c *gin.Context) (any, int, error) {
	uid, err := r.parseID(id)
	if err != nil {
		return nil, 0, err
	}

	body, err := restful.Bind[UpdateUserReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	user, ok := r.db.users[uid]
	if !ok {
		return nil, 0, restful.Abort(http.StatusNotFound, "user not found")
	}

	if body.Name != nil {
		user.Name = *body.Name
	}
	if body.Age != nil {
		user.Age = *body.Age
	}
	r.db.users[uid] = user
	return user, http.StatusOK, nil
}

// DELETE /users/:id
func (r *UserResource) Delete(id string, c *gin.Context) (any, int, error) {
	uid, err := r.parseID(id)
	if err != nil {
		return nil, 0, err
	}

	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, ok := r.db.users[uid]; !ok {
		return nil, 0, restful.Abort(http.StatusNotFound, "user not found")
	}
	delete(r.db.users, uid)
	return nil, http.StatusNoContent, nil
}

func (r *UserResource) findUser(id string) (User, error) {
	uid, err := r.parseID(id)
	if err != nil {
		return User{}, err
	}

	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	user, ok := r.db.users[uid]
	if !ok {
		return User{}, restful.Abort(http.StatusNotFound, "user not found")
	}
	return user, nil
}

func (r *UserResource) parseID(id string) (int, error) {
	uid, err := strconv.Atoi(id)
	if err != nil {
		return 0, restful.Abort(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", id))
	}
	return uid, nil
}

// --- Post 리소스 (읽기 전용 + 생성만) ---

type CreatePostReq struct {
	AuthorID int    `json:"author_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Body     string `json:"body" binding:"required"`
}

type PostResource struct {
	db *DB
}

// GET /posts?author_id=...&page=...&per_page=... — 페이지네이션 + 필터
func (r *PostResource) List(c *gin.Context) (any, int, error) {
	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	authorFilter, _ := strconv.Atoi(c.Query("author_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	all := make([]Post, 0)
	for _, p := range r.db.posts {
		if authorFilter != 0 && p.AuthorID != authorFilter {
			continue
		}
		all = append(all, p)
	}

	start := (page - 1) * perPage
	if start > len(all) {
		start = len(all)
	}
	end := start + perPage
	if end > len(all) {
		end = len(all)
	}

	return gin.H{
		"posts":    all[start:end],
		"total":    len(all),
		"page":     page,
		"per_page": perPage,
	}, http.StatusOK, nil
}

// GET /posts/:id
func (r *PostResource) Get(id string, c *gin.Context) (any, int, error) {
	pid, err := strconv.Atoi(id)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, fmt.Sprintf("invalid id: %s", id))
	}

	r.db.mu.RLock()
	defer r.db.mu.RUnlock()

	post, ok := r.db.posts[pid]
	if !ok {
		return nil, 0, restful.Abort(http.StatusNotFound, "post not found")
	}
	return post, http.StatusOK, nil
}

// POST /posts — 작성자 존재 여부 검증
func (r *PostResource) Post(c *gin.Context) (any, int, error) {
	body, err := restful.Bind[CreatePostReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.db.mu.Lock()
	defer r.db.mu.Unlock()

	if _, ok := r.db.users[body.AuthorID]; !ok {
		return nil, 0, restful.Abort(http.StatusBadRequest, "author not found")
	}

	post := Post{ID: r.db.allocID(), AuthorID: body.AuthorID, Title: body.Title, Body: body.Body}
	r.db.posts[post.ID] = post
	return post, http.StatusCreated, nil
}

// Put, Patch, Delete 미구현 → 해당 라우트 미등록

func main() {
	engine := gin.Default()
	db := NewDB()
	api := restful.NewAPI(engine, "/api/v1")

	api.AddResource("/users", &UserResource{db: db})
	api.AddResource("/posts", &PostResource{db: db})

	if err := engine.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
