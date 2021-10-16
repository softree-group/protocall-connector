package handlers

import (
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/application"
	"protocall/config"
	"protocall/domain/entity"
	"time"
)

const (
	sessionCookie = "session_id"
)

func createCookie() *fasthttp.Cookie {
	token, _ := uuid.GenerateUUID()

	authCookie := fasthttp.Cookie{}
	authCookie.SetKey(sessionCookie)
	authCookie.SetValue(token)
	authCookie.SetDomain("." + viper.GetString(config.ServerDomain))
	authCookie.SetPath("/")
	authCookie.SetExpire(time.Now().Add(24 * time.Hour))
	authCookie.SetHTTPOnly(true)
	authCookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	authCookie.SetSecure(false)
	return &authCookie
}

func session(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		ctx.SetStatusCode(204)
		return
	}

	user := apps.User.Find(string(sessionID))
	if user == nil {
		ctx.Response.Header.DelCookie(sessionCookie)
		ctx.SetStatusCode(204)
		return
	}

	account := apps.AsteriskAccount.Get(user.AsteriskAccount)
	if account == nil {
		account = apps.AsteriskAccount.GetFree()
		if account == nil {
			ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
			// TODO: wait free account
			return
		}
	}
}

func createSession(ctx *fasthttp.RequestCtx, apps *application.Applications) (*entity.User, *entity.AsteriskAccount) {
	account := apps.AsteriskAccount.GetFree()
	if account == nil {
		ctx.Error("Sorry, we are busy ;(", fasthttp.StatusServiceUnavailable)
		return nil, nil
	}

	user := &entity.User{}

	err := json.Unmarshal(ctx.PostBody(), user)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return nil, nil
	}

	cookie := createCookie()
	ctx.Response.Header.SetCookie(cookie)

	user.SessionID = string(cookie.Value())
	user.AsteriskAccount = account.Username
	apps.AsteriskAccount.Take(account.Username, user.SessionID)

	apps.User.Save(user)
	return user, account
}
