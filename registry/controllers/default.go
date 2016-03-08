package controllers

import (
	"github.com/astaxie/beego"
	"log"
	"github.com/amyangfei/dmux/registry/models"
	"github.com/amyangfei/dmux/store"
)

var m *models.EntryManager = &models.EntryManager{}

type RegistryController struct {
	beego.Controller
}

func InitEntryManager(storage store.Storage) {
	m.Storage = storage
}

func (c *RegistryController) ListEntries() {
	entries, err := m.All()
	if err != nil {
		log.Printf("get all entryies error: %v", err)
		c.Abort("500")
	}
	c.Data["json"] = entries
	c.ServeJSON()
}

func (c *RegistryController) NewEntry() {
	tag := c.GetString("tag", "")
	path := c.GetString("path", "")
	if tag == "" || path == "" {
		c.Abort("400")
	}

	entry := &models.Entry{
		Tag:  tag,
		Path: path,
	}

	if err := m.Save(entry); err != nil {
		log.Printf("save entry with error: %v", err)
		c.Abort("500")
	}
	c.Data["json"] = map[string]string{"status": "ok"}
	c.ServeJSON()
}

func (c *RegistryController) GetEntry() {
	path := c.GetString("path", "")
	if path == "" {
		c.Abort("400")
	}
	entry, err := m.Get(models.EntryPrefix + path)
	if err != nil {
		log.Printf("query entry with error: %v", err)
		c.Abort("500")
	}
	c.Data["json"] = entry
	c.ServeJSON()
}

func (c *RegistryController) DelEntry() {
	path := c.GetString("path", "")
	if path == "" {
		c.Abort("400")
	}
	err := m.Delete(models.EntryPrefix + path)
	if err != nil {
		log.Printf("delete entry with error: %v", err)
		c.Abort("500")
	}
	c.Data["json"] = map[string]string{"status": "ok"}
	c.ServeJSON()
}
