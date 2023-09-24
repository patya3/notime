package note

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

const (
	format string = "%d : %s\n"
)

type Note struct {
	gorm.Model
	Title       string
	Description string
}

type NoteRepo struct {
	DB *gorm.DB
}

func (g *NoteRepo) GetNoteByID(id uint) (Note, error) {
	var note Note
	if err := g.DB.Where("id = ?", id).First(&note).Error; err != nil {
		return note, fmt.Errorf("Cannot find note: %v", err)
	}
	return note, nil
}

func (g *NoteRepo) GetAllNotes() ([]Note, error) {
	var notes []Note
	if err := g.DB.Find(&notes).Error; err != nil {
		return notes, fmt.Errorf("No notes found: %v", err)
	}
	return notes, nil
}

func (g *NoteRepo) CreateNote(title string, description string) (Note, error) {
	note := Note{Title: title, Description: description}
	if err := g.DB.Create(&note).Error; err != nil {
		return note, fmt.Errorf("Cannot create note: %v", err)
	}

	return note, nil
}
