package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/tvvat/project/internal/model"
)

type Handler struct {
	addr   string
	router *gin.Engine
	repo   repository
}

type repository interface {
	GetAllUsers() []model.User
	GetAllGroups() []model.Group
	GetUserByID(id int) (model.User, error)
	GetGroupByID(id int) (model.Group, error)
	DeleteUserByID(id int) error
	DeleteGroupByID(id int) error
	CreateUser(id int, login, password string, groups []int) model.User
	CreateGroup(id int, subgroups []int) model.Group
	UpdateUserByID(id int, login, password string) (model.User, error)
}

func NewHandler(repo repository, addr string) *Handler {
	router := gin.Default()
	return &Handler{
		addr:   addr,
		router: router,
		repo:   repo,
	}
}

func (h *Handler) InitRoutes() {
	h.router.GET("/group/:id", h.getGroupByID)
	h.router.GET("/user/:id", h.getUserByID)
	h.router.GET("/user", h.getUsers)
	h.router.GET("/group", h.getGroup)
	h.router.POST("/user", h.postUser)
	h.router.POST("/group", h.postGroup)
	h.router.PUT("/user", h.updateUser)
	h.router.DELETE("/user/:id", h.deleteUserByID)
	h.router.DELETE("/group/:id", h.deleteGroupByID)
}

func (h *Handler) Run() {
	h.router.Run(h.addr)
}

type UserInput struct {
	ID       int    `json:"id" binding:"required"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
	Groups   []int  `json:"groups"`
}

type GroupInput struct {
	ID     int   `json:"id" binding:"required"`
	Groups []int `json:"groups" `
}

func (h *Handler) getUsers(c *gin.Context) {
	users := h.repo.GetAllUsers()
	c.IndentedJSON(http.StatusOK, gin.H{"data": users})
}

func (h *Handler) getGroup(c *gin.Context) {
	groups := h.repo.GetAllGroups()
	c.IndentedJSON(http.StatusOK, gin.H{"data": groups})
}

func (h *Handler) deleteUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("bad ID: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	h.repo.DeleteUserByID(id)

	c.IndentedJSON(http.StatusOK, "successfully deleted the record")
}

func (h *Handler) deleteGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("bad id: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	h.repo.DeleteGroupByID(id)

	c.IndentedJSON(http.StatusOK, "successfully deleted the record")
}

func (h *Handler) postUser(c *gin.Context) {
	var userInput UserInput
	err := c.ShouldBindJSON(&userInput)
	if err != nil {
		log.Printf("invalid JSON body: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	user := h.repo.CreateUser(
		userInput.ID,
		userInput.Login,
		userInput.Password,
		userInput.Groups,
	)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) postGroup(c *gin.Context) {
	var groupInput GroupInput
	err := c.ShouldBindJSON(&groupInput)
	if err != nil {
		log.Printf("invalid JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	group := h.repo.CreateGroup(
		groupInput.ID,
		groupInput.Groups,
	)

	c.JSON(http.StatusOK, gin.H{"data": group})

}

func (h *Handler) getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("bad id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.repo.GetUserByID(id)
	if err != nil {
		log.Printf("don't find user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "can't get user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (h *Handler) getGroupByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("bad id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	group, err := h.repo.GetGroupByID(id)
	if err != nil {
		log.Printf("don't find group: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "can't get group"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": group})
}

func (h *Handler) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("bad id: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var newUser UserInput
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Printf("invalid JSON body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.UpdateUserByID(id, newUser.Login, newUser.Password)

	c.JSON(http.StatusOK, gin.H{"data": user})
}
