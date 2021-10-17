package handlers

import (
	"encoding/json"
	"github.com/hashicorp/go-uuid"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"protocall/application"
	"protocall/config"
	"time"
)

func start(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.StartConference(user)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}
	_, err = apps.Connector.CreateBridge(conference.ID)
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	conference.BridgeID = conference.ID

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
	})
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.Response.SetBody(data)
<<<<<<< HEAD
	ctx.Response.Header.SetContentType("application/json")
=======
>>>>>>> 977da2b (rebase inbloud)
}

func join(ctx *fasthttp.RequestCtx, apps *application.Applications) {

	meetID := ctx.UserValue("meetID").(string)
	if !apps.Conference.IsExist(meetID) {
		ctx.Error("Conference does not exist", 404)
		return
	}

	user, account := createSession(ctx, apps)
	if user == nil {
		return
	}

	conference, err := apps.Conference.JoinToConference(user, meetID)
	if err != nil {
		ctx.Error(err.Error(), 400)
		return
	}

	data, err := json.Marshal(map[string]interface{}{
		"conference": conference,
		"account":    account,
	})
	if err != nil {
		ctx.Error(err.Error(), 500)
		return
	}

	ctx.Response.SetBody(data)
<<<<<<< HEAD
	ctx.Response.Header.SetContentType("application/json")
}

func getUser(ctx *fasthttp.RequestCtx, apps *application.Applications) *entity.User {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
		return nil
	}

	return apps.User.Find(string(sessionID))
}

func ready(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
		ctx.Error("no session", 400)
		return
	}

	channel, err := apps.Connector.CallAndConnect(user.AsteriskAccount, user.ConferenceID)
	if err != nil {
		ctx.Error(err.Error(), 400)
		return
	}

	user.Channel = channel
	apps.User.Save(user)

	conference := apps.Conference.Get(user.ConferenceID)
	if conference == nil {
		logrus.Error("fail to get conference ", user.ConferenceID)
		return
	}

	if conference.IsRecording {
		err = apps.Conference.StartRecordUser(user, conference.ID)
		if err != nil {
			logrus.Error("fail to start record user: ", err)
			return
		}
	}
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)
	if user == nil {
=======
}

func leave(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	sessionID := ctx.Request.Header.Cookie(sessionCookie)
	if len(sessionID) == 0 {
>>>>>>> 977da2b (rebase inbloud)
		ctx.SetStatusCode(400)
		return
	}

<<<<<<< HEAD
	err := apps.Connector.Disconnect(user.ConferenceID, user.Channel)
	if err != nil {
		logrus.Error("Fail to disconnect: ", err)
	}
	apps.User.Delete(user.SessionID)
	ctx.Response.Header.DelCookie(sessionCookie)
}

func record(ctx *fasthttp.RequestCtx, apps *application.Applications) {
	user := getUser(ctx, apps)

	err := apps.Conference.StartRecord(user, user.ConferenceID)
	if err != nil {
		ctx.SetStatusCode(403)
		logrus.Error("fail to start record: ", err)
		return
	}
}
=======
	apps.User.Delete(string(sessionID))
	ctx.Response.Header.DelCookie(sessionCookie)
}
>>>>>>> 977da2b (rebase inbloud)
