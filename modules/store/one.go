package store

import (
	"github.com/fernandez14/spartangeek-blacker/modules/mail"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strings"
)

type One struct {
	di   *Module
	data *OrderModel
}

func (self *One) Data() *OrderModel {

	return self.data
}

func (self *One) PushAnswer(text, kind string) {

	if kind != "text" && kind != "note" {

		return
	}

	database := self.di.Mongo.Database

	message := MessageModel{
		Content: text,
		Type:    kind,
		Created: time.Now(),
		Updated: time.Now(),
	}

	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$push": bson.M{"messages": message}, "$set": bson.M{"updated_at": time.Now()}})

	if err != nil {
		panic(err)
	}

	if kind == "text" {

		// Send an email async
		go func() {

			mailing := self.di.Mail
			text = strings.Replace(text, "\n", "<br>", -1)
			
			compose := mail.Mail{
				Subject:  "PC Spartana",
				Template: "simple",
				Recipient: []mail.MailRecipient{
					{
						Name:  self.data.User.Name,
						Email: self.data.User.Email,
					},
				},
				FromEmail: "pc@pedidos.spartangeek.com",
				FromName: "Drak Spartan",
				Variables: map[string]string{
					"content": text,
				},
			}

			mailing.Send(compose)
		}()
	}
}

func (self *One) PushTag(tag string) {

	database := self.di.Mongo.Database
	item := TagModel{
		Name: tag,
		Created: time.Now(),
	}

	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$push": bson.M{"tags": item}})

	if err != nil {
		panic(err)
	}
}

func (self *One) PushActivity(name, description string, due_at time.Time) {

	database := self.di.Mongo.Database
	activity := ActivityModel{
		Name: name,
		Description: description,
		Done: false,
		Due: due_at,
		Created: time.Now(),
		Updated: time.Now(),
	}

	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$push": bson.M{"activities": activity}})

	if err != nil {
		panic(err)
	}
}

func (self *One) PushInboundAnswer(text string, mail bson.ObjectId) {

	database := self.di.Mongo.Database

	message := MessageModel{
		Content: text,
		Type:    "inbound",
		RelatedId: mail,
		Created: time.Now(),
		Updated: time.Now(),
	}

	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$push": bson.M{"messages": message}, "$set": bson.M{"unreaded": true, "updated_at": time.Now()}})

	if err != nil {
		panic(err)
	}
}

func (self *One) Stage(name string) {

	// Temp way to validate the name of the stage
	if name != "estimate" && name != "negotiation" && name != "accepted" && name != "awaiting" && name != "closed" {
		return
	}

	database := self.di.Mongo.Database

	// Define steps in order
	steps := []string{"estimate", "negotiation", "accepted", "awaiting", "closed"}
	current := self.data.Pipeline.Step

	if current > 0 {
		current = current-1
	}

	target := 0

	for index, step := range steps {

		if step == name {

			target = index
		}
	}

	named := steps[target]
	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$set": bson.M{"pipeline.step": target+1, "pipeline.current": named, "pipeline.updated_at": time.Now(), "updated_at": time.Now()}})

	if err != nil {
		panic(err)
	}
}

func (self *One) Touch() {

	database := self.di.Mongo.Database
	
	err := database.C("orders").Update(bson.M{"_id": self.data.Id}, bson.M{"$set": bson.M{"unreaded": false}})

	if err != nil {
		panic(err)
	}
}