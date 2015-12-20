package controllers

import (
	"github.com/revel/revel"
	"app/app/models"
	"encoding/json"
)

type DocumentCtrl struct {
	GorpController
}

func (c DocumentCtrl) parseDocument() (models.Document, error) {
	document := models.Document{}
	err := json.NewDecoder(c.Request.Body).Decode(&document)
	return document, err
}

func (c DocumentCtrl) List() revel.Result {
	lastID := parseIntOrDefault(c.Params.Get("lid"), -1)
	limit := parseUintOrDefault(c.Params.Get("limit"), uint64(25))
	documents, err := c.Txn.Select(
		models.Document{},
		`SELECT * FROM Documents WHERE ID > ? LIMIT ?`, lastID, limit,
	)
	if err != nil {
		return c.RenderText("Error trying to get records from DB")
	}
	return c.RenderJson(documents)
}

func (c DocumentCtrl) Add() revel.Result {
	if document, err := c.parseDocument(); err != nil {
		return c.RenderText("Unable to parse the document from JSON.")
	} else {
		if err := c.Txn.Insert(&document); err != nil {
			return c.RenderText("Error inserting record into database.")
		} else {
			return c.RenderJson(document)
		}
	}
}

func (c DocumentCtrl) Get(id int64) revel.Result {
	document := new(models.Document)
	err := c.Txn.SelectOne(document, `SELECT * FROM Documents WHERE id = ?`, id)
	if err != nil {
		return c.RenderText("Error, item probably doesn't exist")
	}
	return c.RenderJson(document)
}

func (c DocumentCtrl) Update(id int64) revel.Result {
	document, err := c.parseDocument()
	if err != nil {
		return c.RenderText("Unable to parse document from JSON")
	}
	document.ID = id
	success, err := c.Txn.Update(&document)
	if err != nil || success == 0 {
		return c.RenderText("Unable to update document.")
	}
	return c.RenderText("Updated %v", id)
}

func (c DocumentCtrl) Delete(id int64) revel.Result {
	success, err := c.Txn.Delete(&models.Document{ID: id})
	if err != nil || success == 0 {
		return c.RenderText("Failed to remove document")
	}
	return c.RenderText("Deleted %v", id)
}
