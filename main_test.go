package main

import (
	"net/http"
	"testing"

	"github.com/zyfdegh/fanach/coreserver/api"

	"gopkg.in/kataras/iris.v6/httptest"
)

func TestCoreServer(t *testing.T) {
	server := newCoreServer()
	e := httptest.New(server, t)

	// test if server started
	e.GET("/").
		Expect().
		Status(http.StatusOK).
		Body().Equal("Fanach core server")

	// register user "tom"
	reqJSON := map[string]string{
		"username":  "tom",
		"password":  "secret",
		"wechat_id": "tomwechat",
		"email":     "tom@email.com",
	}
	expectRespJSON := `
		{
			"success": true,
			"errno": 0,
			"errmsg": "",
			"user": {
				"id": "34b7da764b21d298ef307d04d8152dc5",
				"username": "tom",
				"password": "***",
				"wechat_id": "tomwechat",
				"type": "",
				"email": "tom@email.com"
			}
		}
	`
	gotRespJSON := e.POST("/users").
		WithJSON(reqJSON).
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON.Schema(expectRespJSON)
	gotRespJSON.Object().ValueEqual("success", true)
	gotRespJSON.Object().Value("user").Object().ValueEqual("id", "34b7da764b21d298ef307d04d8152dc5")
	gotRespJSON.Object().Value("user").Object().ValueEqual("username", "tom")
	gotRespJSON.Object().Value("user").Object().ValueEqual("password", "***")
	gotRespJSON.Object().Value("user").Object().ValueEqual("wechat_id", "tomwechat")
	gotRespJSON.Object().Value("user").Object().ValueEqual("email", "tom@email.com")

	// query user "tom"
	expectRespJSON2 := `
		{
			"success": true,
			"errno": 0,
			"errmsg": "",
			"user": {
				"id": "34b7da764b21d298ef307d04d8152dc5",
				"username": "tom",
				"password": "***",
				"wechat_id": "tomwechat",
				"type": "",
				"email": "tom@email.com"
			}
		}
	`
	gotRespJSON2 := e.GET("/users/34b7da764b21d298ef307d04d8152dc5").
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON2.Schema(expectRespJSON2)
	gotRespJSON2.Object().ValueEqual("success", true)
	gotRespJSON2.Object().Value("user").Object().ValueEqual("id", "34b7da764b21d298ef307d04d8152dc5")
	gotRespJSON2.Object().Value("user").Object().ValueEqual("username", "tom")
	gotRespJSON2.Object().Value("user").Object().ValueEqual("password", "***")
	gotRespJSON2.Object().Value("user").Object().ValueEqual("wechat_id", "tomwechat")
	gotRespJSON2.Object().Value("user").Object().ValueEqual("email", "tom@email.com")

	// update user "tom"
	reqJSON3 := map[string]string{
		"password":  "strongpassword",
		"wechat_id": "tom123",
		"email":     "tom@outlook.com",
	}
	expectRespJSON3 := `
		{
			"success": true,
			"errno": 0,
			"errmsg": "",
			"user": {
				"id": "34b7da764b21d298ef307d04d8152dc5",
				"username": "tom",
				"password": "***",
				"wechat_id": "tom123",
				"type": "",
				"email": "tom@outlook.com"
			}
		}
	`
	gotRespJSON3 := e.PUT("/users/34b7da764b21d298ef307d04d8152dc5").
		WithJSON(reqJSON3).
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON3.Schema(expectRespJSON3)
	gotRespJSON3.Object().ValueEqual("success", true)
	gotRespJSON3.Object().Value("user").Object().ValueEqual("id", "34b7da764b21d298ef307d04d8152dc5")
	gotRespJSON3.Object().Value("user").Object().ValueEqual("username", "tom")
	gotRespJSON3.Object().Value("user").Object().ValueEqual("password", "***")
	gotRespJSON3.Object().Value("user").Object().ValueEqual("wechat_id", "tom123")
	gotRespJSON3.Object().Value("user").Object().ValueEqual("email", "tom@outlook.com")

	// register again with name "tom"
	reqJSON4 := map[string]string{
		"username": "tom",
		"password": "secret",
	}
	expectRespJSON4 := `
		{
			"success": false,
			"errno": 409,
			"errmsg": "duplicated username",
			"id": ""
		}
	`
	gotRespJSON4 := e.POST("/users").
		WithJSON(reqJSON4).
		Expect().
		Status(http.StatusConflict).
		JSON()

	// validate response JSON
	gotRespJSON4.Schema(expectRespJSON4)
	gotRespJSON4.Object().ValueEqual("success", false)
	gotRespJSON4.Object().ValueEqual("errno", 409)
	gotRespJSON4.Object().ValueEqual("errmsg", "duplicated username")

	// register user "bob"
	reqJSON5 := map[string]string{
		"username":  "bob",
		"password":  "password",
		"wechat_id": "bob123",
		"email":     "bob@email.com",
	}
	expectRespJSON5 := `
			{
				"success": true,
				"errno": 0,
				"errmsg": "",
				"user": {
					"id": "9f9d51bc70ef21ca5c14f307980a29d8",
					"username": "bob",
					"password": "***",
					"wechat_id": "bob123",
					"type": "",
					"email": "bob@email.com"
				}
			}
		`
	gotRespJSON5 := e.POST("/users").
		WithJSON(reqJSON5).
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON5.Schema(expectRespJSON5)
	gotRespJSON5.Object().ValueEqual("success", true)
	gotRespJSON5.Object().Value("user").Object().ValueEqual("id", "9f9d51bc70ef21ca5c14f307980a29d8")
	gotRespJSON5.Object().Value("user").Object().ValueEqual("username", "bob")
	gotRespJSON5.Object().Value("user").Object().ValueEqual("password", "***")
	gotRespJSON5.Object().Value("user").Object().ValueEqual("wechat_id", "bob123")
	gotRespJSON5.Object().Value("user").Object().ValueEqual("email", "bob@email.com")

	// query all users
	expectRespJSON6 := `
			{
				"success": true,
				"errno": 0,
				"errmsg": "",
				"users": [
					{
						"id": "34b7da764b21d298ef307d04d8152dc5",
						"username": "tom",
						"password": "***",
						"wechat_id": "tom123",
						"type": "",
						"email": "tom@outlook.com"
					},
					{
						"id": "9f9d51bc70ef21ca5c14f307980a29d8",
						"username": "bob",
						"password": "***",
						"wechat_id": "bob123",
						"type": "",
						"email": "bob@email.com"
					}
				]
			}
		`
	gotRespJSON6 := e.GET("/users").
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON6.Schema(expectRespJSON6)
	gotRespJSON6.Object().ValueEqual("success", true)
	gotRespJSON6.Object().Value("users").Array().Element(0).Object().ValueEqual("id", "34b7da764b21d298ef307d04d8152dc5")
	gotRespJSON6.Object().Value("users").Array().Element(1).Object().ValueEqual("id", "9f9d51bc70ef21ca5c14f307980a29d8")
	// This is tired and boring work, I won't do anymore

	// create session(login)
	regJSONSess := map[string]string{
		"username": "bob",
		"password": "password",
	}

	e.POST("/sess").
		WithJSON(regJSONSess).
		Expect().
		Status(http.StatusOK)

	// TODO cookie check
	// cookie := e.POST("/sess").
	// 	WithJSON(regJSONSess).
	// 	Expect().
	// 	Status(http.StatusOK).
	// 	Cookie(api.KeySessID)
	//
	// // check cookie
	// now := time.Now()
	// cookie.Expires().InRange(now, now.Add(24*time.Hour))
	// cookie.Path().Equal("/")

	// delete session (logout)
	e.DELETE("/sess/{key}", api.KeySessID).
		Expect().
		Status(http.StatusOK)

	// clean up
	// test delete user "tom"
	e.DELETE("/users/34b7da764b21d298ef307d04d8152dc5").
		Expect().
		Status(http.StatusOK)

	// test delete user "bob"
	e.DELETE("/users/9f9d51bc70ef21ca5c14f307980a29d8").
		Expect().
		Status(http.StatusOK)

	// test get products
	expectRespJSON7 := `
	{
	  "success": true,
	  "errno": 0,
	  "errmsg": "",
	  "products": [
	    {
	      "name": "Free",
	      "description": "Free account",
	      "price": 0,
	      "price_unit": "￥",
	      "dataflow": 1024,
	      "dataflow_unit": "MB",
	      "expire": 1,
	      "expire_unit": "Month"
	    },
	    {
	      "name": "1元包月",
	      "description": "1元包月, 10GB",
	      "price": 1,
	      "price_unit": "￥",
	      "dataflow": 10240,
	      "dataflow_unit": "MB",
	      "expire": 1,
	      "expire_unit": "Month"
	    }
	  ]
	}
	`
	gotRespJSON7 := e.GET("/prods").
		Expect().
		Status(http.StatusOK).
		JSON()

	// validate response JSON
	gotRespJSON7.Schema(expectRespJSON7)
	gotRespJSON7.Object().ValueEqual("success", true)
	gotRespJSON7.Object().Value("products").Array().NotEmpty()
}
