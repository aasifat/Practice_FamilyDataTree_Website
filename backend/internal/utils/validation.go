package utils

import (
	"family-tree-api/internal/models"
	"family-tree-api/internal/repository"
	"fmt"
)

// ValidatePersonCreation checks all business rules for creating a person
func ValidatePersonCreation(person *models.Person) error {
	// Validate gender
	if person.Gender != "male" && person.Gender != "female" {
		return fmt.Errorf("gender must be 'male' or 'female'")
	}

	// If person has both father and mother, that's valid (couple)
	if person.FatherID != nil && person.MotherID != nil {
		// Validate that both parents exist
		father, _ := repository.GetPersonByID(*person.FatherID)
		if father == nil {
			return fmt.Errorf("father with ID %d not found", *person.FatherID)
		}

		mother, _ := repository.GetPersonByID(*person.MotherID)
		if mother == nil {
			return fmt.Errorf("mother with ID %d not found", *person.MotherID)
		}

		// Validate that father and mother are actually spouses
		if father.SpouseID != nil && *father.SpouseID != mother.ID {
			return fmt.Errorf("father and mother are not spouses")
		}

		return nil
	}

	// If person has only one parent, that's only valid for root level
	if person.FatherID != nil && person.MotherID == nil {
		return fmt.Errorf("person must have both father and mother, or neither (root level)")
	}

	if person.MotherID != nil && person.FatherID == nil {
		return fmt.Errorf("person must have both father and mother, or neither (root level)")
	}

	// Valid: person is at root level (no parents)
	return nil
}

// ValidateSpouseRelationship validates that spouse linkage is correct
func ValidateSpouseRelationship(person1ID int, person2ID int) error {
	person1, err := repository.GetPersonByID(person1ID)
	if err != nil {
		return fmt.Errorf("person 1 not found: %w", err)
	}

	person2, err := repository.GetPersonByID(person2ID)
	if err != nil {
		return fmt.Errorf("person 2 not found: %w", err)
	}

	// Can't marry self
	if person1.ID == person2.ID {
		return fmt.Errorf("a person cannot be their own spouse")
	}

	// Optional: Add gender preference rule (optional based on requirements)
	if person1.Gender == person2.Gender {
		// This could be allowed or forbidden depending on requirements
		// For now, allowing it
	}

	return nil
}

// EnsureSpouseBidirectional ensures both spouses point to each other
func EnsureSpouseBidirectional(personID int, spouseID int) error {
	person, err := repository.GetPersonByID(personID)
	if err != nil {
		return err
	}

	spouse, err := repository.GetPersonByID(spouseID)
	if err != nil {
		return err
	}

	// Update person to point to spouse
	person.SpouseID = &spouseID
	repository.UpdatePerson(person)

	// Update spouse to point to person
	spouse.SpouseID = &personID
	repository.UpdatePerson(spouse)

	return nil
}

// ValidateNoChildrenWithoutBothParents checks that a person has both parents if they have children
func ValidateNoChildrenWithoutBothParents(personID int) error {
	person, err := repository.GetPersonByID(personID)
	if err != nil {
		return err
	}

	// If person is trying to have children without both parents, that's invalid
	children, err := repository.GetChildren(person.FatherID, person.MotherID)
	if err != nil {
		return err
	}

	if len(children) > 0 {
		if person.FatherID == nil || person.MotherID == nil {
			return fmt.Errorf("a person must have both father and mother to have children")
		}
	}

	return nil
}
