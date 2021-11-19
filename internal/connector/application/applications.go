package application

import (
	"context"
	"protocall/application/app"
	"protocall/application/applications"
	"protocall/domain/repository"
	"protocall/domain/services"
	"protocall/infrastructure"
	"protocall/infrastructure/bus"
	"protocall/internal/config"

	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Applications struct {
	Listener        applications.EventListener
	Snoopy          applications.Snoopy
	User            applications.User
	AsteriskAccount applications.AsteriskAccount
	Conference      applications.Conference
	Connector       applications.Connector
	AMI             *app.AMIAsterisk
	Socket          applications.Socket
	Bus             services.Bus
}

func New(reps repository.Repositories) *Applications {
	logrus.Info("start connect")
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARIApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketURL),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	logrus.Info("end connect")
	if err != nil {
		logrus.Fatal("cannot connect to asterisk: ", err)
	}

	socketService := infrastructure.NewCentrifugo()
	socketApp := app.NewSocket(socketService)

	connector := app.NewConnector(ariClient, reps)

	userApp := app.NewUser(reps)
	conferenceApp := app.NewConference(reps, ariClient)
	asteriskApp := app.NewAsteriskAccount(reps, viper.GetString(config.ARIAccountsFile))
	ami, err := app.NewAMIAsterisk(context.Background(), viper.GetString(config.AMIHost)+":"+viper.GetString(config.AMIPort), viper.GetString(config.AMIUser), viper.GetString(config.AMIPassword))
	if err != nil {
		logrus.Fatal("Fail connect to ami: ", err)
	}

	busService := bus.New()

	return &Applications{
		Listener:        app.NewListener(reps, ariClient, app.NewHandler(ariClient, reps, connector), userApp, conferenceApp, asteriskApp, socketApp),
		Snoopy:          app.NewSnoopy(busService),
		Conference:      conferenceApp,
		AsteriskAccount: asteriskApp,
		User:            userApp,
		Connector:       connector,
		AMI:             ami,
		Socket:          socketApp,
		Bus:             busService,
	}
}