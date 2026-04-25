package handlers

import (
	"net/http"
	"strconv"

	"family-tree-api/internal/models"
	"family-tree-api/internal/repository"

	"github.com/gin-gonic/gin"
)

func CreateFamilyTree(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req models.TreeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tree := &models.FamilyTree{
		UserID: userID.(int),
		Name:   req.Name,
	}

	if err := repository.CreateFamilyTree(tree); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create family tree"})
		return
	}

	c.JSON(http.StatusCreated, tree)
}

func GetUserTrees(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	trees, err := repository.GetUserTrees(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if trees == nil {
		trees = []models.FamilyTree{}
	}

	c.JSON(http.StatusOK, trees)
}

func GetFamilyTree(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	tree, err := repository.GetFamilyTreeByID(treeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tree not found"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Get all members
	members, err := repository.GetTreeMembers(treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}

	tree.Members = members
	c.JSON(http.StatusOK, tree)
}

func UpdateFamilyTree(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	tree, err := repository.GetFamilyTreeByID(treeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tree not found"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req models.TreeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		tree.Name = req.Name
	}

	if err := repository.UpdateFamilyTree(tree); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tree"})
		return
	}

	c.JSON(http.StatusOK, tree)
}

func DeleteFamilyTree(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	tree, err := repository.GetFamilyTreeByID(treeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tree not found"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := repository.DeleteFamilyTree(treeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tree"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tree deleted successfully"})
}

func GetTreeMembers(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	members, err := repository.GetTreeMembers(treeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get members"})
		return
	}

	c.JSON(http.StatusOK, members)
}
