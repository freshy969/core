package comments

import (
	"html"
	"time"

	"github.com/tryanzu/core/core/content"
	"gopkg.in/mgo.v2/bson"
)

// UpsertComment performs validations before upserting data struct
func UpsertComment(deps Deps, c Comment) (comment Comment, err error) {
	if c.Id.Valid() == false {
		c.Id = bson.NewObjectId()
		c.Created = time.Now()
	}

	c.Content = html.EscapeString(c.Content)
	c.Updated = time.Now()

	// Post process comment content.
	processed, err := content.Postprocess(deps, c)
	if err != nil {
		return
	}

	c = processed.(Comment)
	_, err = deps.Mgo().C("comments").UpsertId(c.Id, bson.M{"$set": c})
	if err != nil {
		return
	}

	if c.ReplyType == "post" {
		err = deps.Mgo().C("posts").UpdateId(c.ReplyTo, bson.M{"$inc": bson.M{"comments.count": 1}})
	}

	comment = c
	return
}
