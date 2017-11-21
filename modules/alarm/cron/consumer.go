package cron

import (
	"encoding/json"
	"github.com/Cepave/open-falcon-backend/common/model"
	"github.com/Cepave/open-falcon-backend/modules/alarm/api"
	"github.com/Cepave/open-falcon-backend/modules/alarm/g"
	"github.com/Cepave/open-falcon-backend/modules/alarm/redis"
)

func consume(event *model.Event, isHigh bool) {
	actionId := event.ActionId()
	if actionId <= 0 {
		return
	}

	action := api.GetAction(actionId)
	if action == nil {
		return
	}

	if action.Callback == 1 {
		HandleCallback(event, action)
		return
	}

	if isHigh {
		consumeHighEvents(event, action)
	} else {
		consumeLowEvents(event, action)
	}
}

// 高优先级的不做报警合并
func consumeHighEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	phones, mails := api.ParseTeams(action.Uic)

	smsContent := GenerateSmsContent(event)
	mailContent := GenerateMailContent(event)
	QQContent := GenerateQQContent(event)

	if event.Priority() < 3 {
		redis.WriteSms(phones, smsContent)
	}

	redis.WriteMail(mails, smsContent, mailContent)
	redis.WriteQQ(mails, smsContent, QQContent)
	ParseUserServerchan(event, action)
}

// 低优先级的做报警合并
func consumeLowEvents(event *model.Event, action *api.Action) {
	if action.Uic == "" {
		return
	}

	if event.Priority() < 3 {
		ParseUserSms(event, action)
	}

	ParseUserMail(event, action)
	ParseUserQQ(event, action)
	ParseUserServerchan(event, action)
}

func ParseUserSms(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	content := GenerateSmsContent(event)
	metric := event.Metric()
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserSmsQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := SmsDto{
			Priority: priority,
			Metric:   metric,
			Content:  content,
			Phone:    user.Phone,
			Status:   status,
		}
		// if Phone field with empty string, just skip it.
		if dto.Phone == "" {
			log.Println("found SmsDto has phone field with empty string.")
			continue
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Errorf("json marshal SmsDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("LPUSH redis: %v. Fail: %v. dto: %s", queue, err, bs)
		}
	}
}

func ParseUserMail(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateMailContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserMailQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := MailDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Errorf("json marshal MailDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("LPUSH redis %v. Fail: %v. dto: %s", queue, err, bs)
		}
	}
}

func ParseUserQQ(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)

	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateQQContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserQQQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := QQDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Email:    user.Email,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Errorf("json marshal QQDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("LPUSH redis: %v. fail: %v. dto: %s", queue, err, bs)
		}
	}
}

func ParseUserServerchan(event *model.Event, action *api.Action) {
	userMap := api.GetUsers(action.Uic)
	metric := event.Metric()
	subject := GenerateSmsContent(event)
	content := GenerateServerchanContent(event)
	status := event.Status
	priority := event.Priority()

	queue := g.Config().Redis.UserServerchanQueue

	rc := g.RedisConnPool.Get()
	defer rc.Close()

	for _, user := range userMap {
		dto := ServerchanDto{
			Priority: priority,
			Metric:   metric,
			Subject:  subject,
			Content:  content,
			Username: user.Name,
			Sckey:    user.IM,
			Status:   status,
		}
		bs, err := json.Marshal(dto)
		if err != nil {
			log.Errorf("json marshal ServerchanDto fail: %v", err)
			continue
		}

		_, err = rc.Do("LPUSH", queue, string(bs))
		if err != nil {
			log.Errorf("LPUSH redis: %v. fail: %v. dto: %s", queue, err, bs)
		}
	}
}
