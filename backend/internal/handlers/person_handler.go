package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"family-tree-api/internal/models"
	"family-tree-api/internal/repository"
	"family-tree-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func CreatePerson(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tree, err := repository.GetFamilyTreeByID(treeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tree not found"})
		return
	}

	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req models.PersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person := &models.Person{
		TreeID:   treeID,
		Name:     req.Name,
		Gender:   req.Gender,
		FatherID: req.FatherID,
		MotherID: req.MotherID,
		SpouseID: req.SpouseID,
	}

	// Validate person creation
	if err := utils.ValidatePersonCreation(person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.CreatePerson(person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create person"})
		return
	}

	c.JSON(http.StatusCreated, person)
}

func UpdatePerson(c *gin.Context) {
	personIDStr := c.Param("personId")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	person, err := repository.GetPersonByID(personID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tree, _ := repository.GetFamilyTreeByID(person.TreeID)
	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var req models.PersonUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		person.Name = req.Name
	}
	if req.Gender != "" {
		person.Gender = req.Gender
	}

	if err := repository.UpdatePerson(person); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update person"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func GetPerson(c *gin.Context) {
	personIDStr := c.Param("personId")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	person, err := repository.GetPersonByID(personID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func DeletePerson(c *gin.Context) {
	personIDStr := c.Param("personId")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	person, err := repository.GetPersonByID(personID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	// Check authorization
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tree, _ := repository.GetFamilyTreeByID(person.TreeID)
	role, _ := c.Get("role")
	if tree.UserID != userID.(int) && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	if err := repository.DeletePerson(personID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete person"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person deleted successfully"})
}

func GetChildren(c *gin.Context) {
	personIDStr := c.Param("personId")
	personID, err := strconv.Atoi(personIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	person, err := repository.GetPersonByID(personID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	// Get children only if person has both mother and father (is part of a couple)
	children, err := repository.GetChildrenWithBothParents(person.FatherID, person.MotherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get children"})
		return
	}

	c.JSON(http.StatusOK, children)
}

func SearchPeople(c *gin.Context) {
	treeIDStr := c.Param("id")
	treeID, err := strconv.Atoi(treeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tree id"})
		return
	}

	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search term required"})
		return
	}

	people, err := repository.SearchPeople(treeID, searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to search: %v", err)})
		return
	}

	c.JSON(http.StatusOK, people)
}
