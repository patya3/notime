package note

import (
	"fmt"

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
	if err := g.DB.Order("created_at desc").Find(&notes).Error; err != nil {
		return notes, fmt.Errorf("No notes found: %v", err)
	}
	return notes, nil
}

func (g *NoteRepo) GetAllNoteIds() ([]int, error) {
	var noteIds []int
	var note Note
	if err := g.DB.Model(&note).Order("created_at desc").Select([]string{"id"}).Find(&noteIds).Error; err != nil {
		return noteIds, fmt.Errorf("No notes found: %v", err)
	}
	return noteIds, nil
}

func (g *NoteRepo) CreateNote(note Note) (Note, error) {
	if err := g.DB.Create(&note).Error; err != nil {
		return note, fmt.Errorf("Cannot create note: %v", err)
	}

	return note, nil
}
