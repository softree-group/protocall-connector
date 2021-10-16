package snoopy

import (
	"context"
	"github.com/CyCoreSystems/ari/v5"
	"github.com/CyCoreSystems/ari/v5/client/native"
	"github.com/CyCoreSystems/ari/v5/ext/record"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"protocall/config"
)

type Snoopy struct {
	ari ari.Client
}

func New() *Snoopy {
	ariClient, err := native.Connect(&native.Options{
		Application:  viper.GetString(config.ARISnoopyApplication),
		URL:          viper.GetString(config.ARIUrl),
		WebsocketURL: viper.GetString(config.ARIWebsocketUrl),
		Username:     viper.GetString(config.ARIUser),
		Password:     viper.GetString(config.ARIPassword),
	})
	if err != nil {
		logrus.Fatal("Fail to connect snoopy app")
	}
	return &Snoopy{ari: ariClient}
}

func (s Snoopy) channelHandler(channel *ari.ChannelHandle) {
	sub := channel.Subscribe(ari.Events.All)
	end := channel.Subscribe(ari.Events.StasisEnd)

	defer sub.Cancel()
	defer end.Cancel()

	ctx := context.Background()
	rec := record.Record(ctx, channel)

	for {
		select {
		case event := <-sub.Events():
			logrus.Info("In SPY: ", event.GetType())
		case <-end.Events():
			logrus.Info("saving record for ", channel.ID())
			res, err := rec.Result()
			if err != nil {
				logrus.Error("Fail to get result from record for channel ", channel.ID(), ". Error: ", err)
				return
			}
			err = res.Save(channel.ID())
			if err != nil {
				logrus.Error("fail to save result record for channel ", channel.ID(), ". Error: ", err)
				return
			}
			logrus.Info("saved record for ", channel.ID())
			return
		}
	}
}

func (s Snoopy) listen() {
	start := s.ari.Bus().Subscribe(nil, ari.Events.StasisStart)
	for {
		select {
		case event := <-start.Events():
			value := event.(*ari.StasisStart)

			channel := s.ari.Channel().Get(value.Key(ari.ChannelKey, value.Channel.ID))
			logrus.Info("snoop channel: ", channel.ID())
			go s.channelHandler(channel)
		}
	}
}

func (s Snoopy) Snoop() {
	logrus.Info("Start snooping...")
	s.listen()
}
