package main

import (
	"github.com/kataras/go-sessions/sessiondb/leveldb"
	"github.com/zyfdegh/fanach/coreserver/api"
	"github.com/zyfdegh/fanach/coreserver/db"
	"github.com/zyfdegh/fanach/coreserver/service"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
)

func main() {
	server := newCoreServer()
	server.Listen(":8080")
}

func newCoreServer() *iris.Framework {
	ir := iris.New()
	ir.Adapt(httprouter.New())
	ir.Adapt(iris.DevLogger())

	ir.Get("/", api.GetRoot)

	ir.Post("/users", api.CreateUser)
	ir.Get("/users", api.GetUsers)
	ir.Get("/users/:id", api.GetUser)
	ir.Put("/users/:id", api.ModifyUser)
	ir.Delete("/users/:id", api.DeleteUser)

	ir.Post("/sess", api.PostSession)
	ir.Delete("/sess/:key", api.DeleteSession)

	ir.Get("/prods", api.GetProducts)

	ir.Post("/orders", api.CreateOrder)

	// save sessions to LevelDB(with GC)
	// import "gopkg.in/kataras/iris.v6/adaptors/sessions"
	// import "github.com/kataras/go-sessions/sessiondb/leveldb"
	sess := sessions.New(sessions.Config{
		Cookie:                      api.KeySessID,
		DecodeCookie:                false,
		Expires:                     service.SessionAliveDuration,
		CookieLength:                32,
		DisableSubdomainPersistence: false,
	})
	db := leveldb.New(leveldb.Config{
		Path:         db.SessDBFile,
		CleanTimeout: service.SessionGCDuration,
		MaxAge:       service.SessionAliveDuration,
	})
	sess.UseDatabase(db)
	ir.Adapt(sess)

	return ir
}
