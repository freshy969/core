package feed

import (
	"github.com/fernandez14/spartangeek-blacker/model"
	"github.com/fernandez14/spartangeek-blacker/modules/exceptions"
	"github.com/fernandez14/spartangeek-blacker/modules/user"
	"github.com/fernandez14/spartangeek-blacker/modules/search"
	"github.com/fernandez14/spartangeek-blacker/mongo"
	"github.com/xuyu/goredis"
	"gopkg.in/mgo.v2/bson"
)

var lightPostFields bson.M = bson.M{"_id": 1, "title": 1, "slug": 1, "category": 1, "user_id": 1, "pinned": 1, "created_at": 1, "updated_at": 1, "type": 1, "content": 1, "comments.set": bson.M{"$elemMatch": bson.M{"chosen": true}}}

type FeedModule struct {
	Mongo        *mongo.Service               `inject:""`
	Errors       *exceptions.ExceptionsModule `inject:""`
	CacheService *goredis.Redis               `inject:""`
	Search       *search.Module               `inject:""`
	User         *user.Module                 `inject:""`
}

func (self *FeedModule) Post(post interface{}) (*Post, error) {

	module := self

	switch post.(type) {
	case bson.ObjectId:

		this := model.Post{}
		database := self.Mongo.Database

		// Use user module reference to get the user and then create the user gaming instance
		err := database.C("posts").FindId(post.(bson.ObjectId)).One(&this)

		if err != nil {

			return nil, exceptions.NotFound{"Invalid post id. Not found."}
		}

		post_object := &Post{data: this, di: module}

		return post_object, nil

	case model.Post:

		this := post.(model.Post)
		post_object := &Post{data: this, di: module}

		return post_object, nil

	default:
		panic("Unkown argument")
	}
}

func (module *FeedModule) LightPost(post interface{}) (*LightPost, error) {

	switch post.(type) {
	case bson.ObjectId:

		scope := LightPostModel{}
		database := module.Mongo.Database

		// Use light post model 
		err := database.C("posts").FindId(post.(bson.ObjectId)).Select(lightPostFields).One(&scope)

		if err != nil {

			return nil, exceptions.NotFound{"Invalid post id. Not found."}
		}

		post_object := &LightPost{data: scope, di: module}

		return post_object, nil

	default:
		panic("Unkown argument")
	}
}

func (module *FeedModule) LightPosts(posts interface{}) ([]LightPostModel, error) {

	switch posts.(type) {
	case []bson.ObjectId:

		var list []LightPostModel

		database := module.Mongo.Database

		// Use light post model 
		err := database.C("posts").Find(bson.M{"_id": bson.M{"$in": posts.([]bson.ObjectId)}}).Select(lightPostFields).All(&list)

		if err != nil {

			return nil, exceptions.NotFound{"Invalid posts id. Not found."}
		}

		return list, nil

	case bson.M:

		var list []LightPostModel

		database := module.Mongo.Database

		// Use light post model 
		err := database.C("posts").Find(posts.(bson.M)).Select(lightPostFields).All(&list)

		if err != nil {

			return nil, exceptions.NotFound{"Invalid posts criteria. Not found."}
		}

		return list, nil

	default:
		panic("Unkown argument")
	}
}

func (module *FeedModule) Posts(limit, offset int) List {

	list := List{
		module: module,
		limit: limit,
		offset: offset,
	}

	return list
}