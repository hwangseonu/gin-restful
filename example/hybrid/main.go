// 하이브리드 예제: restful 리소스와 일반 Gin 핸들러가 공존하는 구조
package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/hwangseonu/gin-restful"
)

// --- 모델 ---

type Task struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
	UserID string `json:"user_id"`
}

// --- 인메모리 저장소 ---

type TaskStore struct {
	mu     sync.RWMutex
	tasks  map[int]Task
	nextID int
}

func NewTaskStore() *TaskStore {
	return &TaskStore{tasks: make(map[int]Task), nextID: 1}
}

// --- Task 리소스 (restful) ---

type CreateTaskReq struct {
	Title string `json:"title" binding:"required"`
}

type TaskResource struct {
	store *TaskStore
}

// GET /tasks — 현재 사용자의 태스크만 반환
func (r *TaskResource) List(c *gin.Context) (any, int, error) {
	userID := c.GetString("user_id")

	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	tasks := make([]Task, 0)
	for _, t := range r.store.tasks {
		if t.UserID == userID {
			tasks = append(tasks, t)
		}
	}
	return gin.H{"tasks": tasks}, http.StatusOK, nil
}

// GET /tasks/:id
func (r *TaskResource) Get(id string, c *gin.Context) (any, int, error) {
	task, err := r.findTask(id, c.GetString("user_id"))
	if err != nil {
		return nil, 0, err
	}
	return task, http.StatusOK, nil
}

// POST /tasks
func (r *TaskResource) Post(c *gin.Context) (any, int, error) {
	body, err := restful.Bind[CreateTaskReq](c)
	if err != nil {
		return nil, 0, restful.Abort(http.StatusBadRequest, err.Error())
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	task := Task{
		ID:     r.store.nextID,
		Title:  body.Title,
		Done:   false,
		UserID: c.GetString("user_id"),
	}
	r.store.tasks[task.ID] = task
	r.store.nextID++
	return task, http.StatusCreated, nil
}

// DELETE /tasks/:id
func (r *TaskResource) Delete(id string, c *gin.Context) (any, int, error) {
	task, err := r.findTask(id, c.GetString("user_id"))
	if err != nil {
		return nil, 0, err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	delete(r.store.tasks, task.ID)
	return nil, http.StatusNoContent, nil
}

func (r *TaskResource) findTask(id string, userID string) (Task, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return Task{}, restful.Abort(http.StatusBadRequest, "invalid id")
	}

	r.store.mu.RLock()
	defer r.store.mu.RUnlock()

	task, ok := r.store.tasks[tid]
	if !ok {
		return Task{}, restful.Abort(http.StatusNotFound, "task not found")
	}
	if task.UserID != userID {
		return Task{}, restful.Abort(http.StatusForbidden, "access denied")
	}
	return task, nil
}

// --- 일반 Gin 핸들러들 ---

// GET /health — 헬스체크
func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GET /stats — 전체 통계 (관리자용, 미들웨어 없는 간단 핸들러)
func statsHandler(store *TaskStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		store.mu.RLock()
		defer store.mu.RUnlock()

		total := len(store.tasks)
		done := 0
		for _, t := range store.tasks {
			if t.Done {
				done++
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"total_tasks":     total,
			"completed_tasks": done,
			"pending_tasks":   total - done,
		})
	}
}

// 간단한 인증 미들웨어 (헤더에서 user_id 추출)
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "X-User-ID header required"})
			return
		}
		c.Set("user_id", userID)
		c.Next()
	}
}

func main() {
	engine := gin.Default()
	store := NewTaskStore()

	// --- 일반 Gin 라우트 (restful 밖) ---
	engine.GET("/health", healthHandler)
	engine.GET("/stats", statsHandler(store))

	// --- restful 리소스 (인증 미들웨어가 적용된 그룹) ---
	authorized := engine.Group("/", authMiddleware())
	api := restful.NewAPI(authorized, "/api/v1")

	api.AddResource("/tasks", &TaskResource{store: store})

	log.Println("Hybrid server running on :8080")
	log.Println("  GET  /health          — no auth")
	log.Println("  GET  /stats           — no auth")
	log.Println("  *    /api/v1/tasks/** — requires X-User-ID header")
	if err := engine.Run(":8080"); err != nil {
		log.Fatalln(err)
	}
}
