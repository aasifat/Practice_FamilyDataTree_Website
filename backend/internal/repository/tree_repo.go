package repository

import (
	"database/sql"
	"family-tree-api/internal/database"
	"family-tree-api/internal/models"
	"fmt"
)

func CreateFamilyTree(tree *models.FamilyTree) error {
	query := `
		INSERT INTO family_trees (user_id, name, created_at, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRow(query, tree.UserID, tree.Name).
		Scan(&tree.ID, &tree.CreatedAt, &tree.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create family tree: %w", err)
	}
	return nil
}

func GetFamilyTreeByID(treeID int) (*models.FamilyTree, error) {
	tree := &models.FamilyTree{}
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM family_trees WHERE id = $1
	`

	err := database.DB.QueryRow(query, treeID).Scan(
		&tree.ID, &tree.UserID, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("family tree not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get family tree: %w", err)
	}
	return tree, nil
}

func GetUserTrees(userID int) ([]models.FamilyTree, error) {
	query := `
		SELECT id, user_id, name, created_at, updated_at
		FROM family_trees WHERE user_id = $1 ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get trees: %w", err)
	}
	defer rows.Close()

	trees := make([]models.FamilyTree, 0)
	for rows.Next() {
		var tree models.FamilyTree
		err := rows.Scan(&tree.ID, &tree.UserID, &tree.Name, &tree.CreatedAt, &tree.UpdatedAt)
		if err != nil {
			return nil, err
		}
		trees = append(trees, tree)
	}

	return trees, rows.Err()
}

func UpdateFamilyTree(tree *models.FamilyTree) error {
	query := `
		UPDATE family_trees 
		SET name = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := database.DB.Exec(query, tree.Name, tree.ID)
	if err != nil {
		return fmt.Errorf("failed to update family tree: %w", err)
	}
	return nil
}

func DeleteFamilyTree(treeID int) error {
	query := `DELETE FROM family_trees WHERE id = $1`
	_, err := database.DB.Exec(query, treeID)
	if err != nil {
		return fmt.Errorf("failed to delete family tree: %w", err)
	}
	return nil
}
