package repository

import (
	"database/sql"
	"family-tree-api/internal/database"
	"family-tree-api/internal/models"
	"fmt"
)

func CreatePerson(person *models.Person) error {
	query := `
		INSERT INTO people (tree_id, name, gender, father_id, mother_id, spouse_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRow(query, person.TreeID, person.Name, person.Gender, person.FatherID, person.MotherID, person.SpouseID).
		Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create person: %w", err)
	}
	return nil
}

func GetPersonByID(personID int) (*models.Person, error) {
	person := &models.Person{}
	query := `
		SELECT id, tree_id, name, gender, father_id, mother_id, spouse_id, image_url, created_at, updated_at
		FROM people WHERE id = $1
	`

	err := database.DB.QueryRow(query, personID).Scan(
		&person.ID, &person.TreeID, &person.Name, &person.Gender, &person.FatherID, &person.MotherID, &person.SpouseID, &person.ImageURL, &person.CreatedAt, &person.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("person not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get person: %w", err)
	}
	return person, nil
}

func GetTreeMembers(treeID int) ([]models.Person, error) {
	query := `
		SELECT id, tree_id, name, gender, father_id, mother_id, spouse_id, image_url, created_at, updated_at
		FROM people WHERE tree_id = $1 ORDER BY id
	`

	rows, err := database.DB.Query(query, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree members: %w", err)
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.TreeID, &person.Name, &person.Gender, &person.FatherID, &person.MotherID, &person.SpouseID, &person.ImageURL, &person.CreatedAt, &person.UpdatedAt)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, rows.Err()
}

func GetChildren(fatherID *int, motherID *int) ([]models.Person, error) {
	query := `
		SELECT id, tree_id, name, gender, father_id, mother_id, spouse_id, image_url, created_at, updated_at
		FROM people WHERE father_id = $1 OR mother_id = $2 ORDER BY id
	`

	rows, err := database.DB.Query(query, fatherID, motherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}
	defer rows.Close()

	var children []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.TreeID, &person.Name, &person.Gender, &person.FatherID, &person.MotherID, &person.SpouseID, &person.ImageURL, &person.CreatedAt, &person.UpdatedAt)
		if err != nil {
			return nil, err
		}
		children = append(children, person)
	}

	return children, rows.Err()
}

func UpdatePerson(person *models.Person) error {
	query := `
		UPDATE people 
		SET name = $1, gender = $2, spouse_id = $3, image_url = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`

	_, err := database.DB.Exec(query, person.Name, person.Gender, person.SpouseID, person.ImageURL, person.ID)
	if err != nil {
		return fmt.Errorf("failed to update person: %w", err)
	}
	return nil
}

func DeletePerson(personID int) error {
	query := `DELETE FROM people WHERE id = $1`
	_, err := database.DB.Exec(query, personID)
	if err != nil {
		return fmt.Errorf("failed to delete person: %w", err)
	}
	return nil
}

func SearchPeople(treeID int, searchTerm string) ([]models.Person, error) {
	query := `
		SELECT id, tree_id, name, gender, father_id, mother_id, spouse_id, image_url, created_at, updated_at
		FROM people WHERE tree_id = $1 AND name ILIKE $2
	`

	rows, err := database.DB.Query(query, treeID, "%"+searchTerm+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search people: %w", err)
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.TreeID, &person.Name, &person.Gender, &person.FatherID, &person.MotherID, &person.SpouseID, &person.ImageURL, &person.CreatedAt, &person.UpdatedAt)
		if err != nil {
			return nil, err
		}
		people = append(people, person)
	}

	return people, rows.Err()
}

// GetPersonByIDWithFamily returns person with spouse loaded
func GetPersonByIDWithFamily(personID int) (*models.Person, error) {
	person, err := GetPersonByID(personID)
	if err != nil {
		return nil, err
	}

	// Load spouse if exists
	if person.SpouseID != nil {
		spouse, err := GetPersonByID(*person.SpouseID)
		if err == nil {
			person.Spouse = spouse
		}
	}

	return person, nil
}

// GetChildrenWithBothParents returns children only if both parents exist
func GetChildrenWithBothParents(fatherID *int, motherID *int) ([]models.Person, error) {
	if fatherID == nil || motherID == nil {
		return []models.Person{}, nil
	}

	query := `
		SELECT id, tree_id, name, gender, father_id, mother_id, spouse_id, image_url, created_at, updated_at
		FROM people WHERE father_id = $1 AND mother_id = $2 ORDER BY id
	`

	rows, err := database.DB.Query(query, fatherID, motherID)
	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}
	defer rows.Close()

	var children []models.Person
	for rows.Next() {
		var person models.Person
		err := rows.Scan(&person.ID, &person.TreeID, &person.Name, &person.Gender, &person.FatherID, &person.MotherID, &person.SpouseID, &person.ImageURL, &person.CreatedAt, &person.UpdatedAt)
		if err != nil {
			return nil, err
		}
		children = append(children, person)
	}

	return children, rows.Err()
}

// GetParents returns both father and mother
func GetParents(personID int) (*models.Person, *models.Person, error) {
	person, err := GetPersonByID(personID)
	if err != nil {
		return nil, nil, err
	}

	var father, mother *models.Person

	if person.FatherID != nil {
		father, _ = GetPersonByID(*person.FatherID)
	}

	if person.MotherID != nil {
		mother, _ = GetPersonByID(*person.MotherID)
	}

	return father, mother, nil
}
