package main

import (
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	apiToken              = ""
	sourceChannelID int64 = -100
	targetChannelID int64 = -100
	botDebug              = true
)

type MediaGroup struct {
	MediaGroupID string
	Media        []interface{}
	Caption      string
	Timer        *time.Timer
}

var (
	mediaGroups = make(map[string]*MediaGroup)
	mu          sync.Mutex
)

func main() {
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = botDebug

	log.Printf("Бот запущен: %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.ChannelPost != nil {
			post := update.ChannelPost

			if post.Chat.ID == sourceChannelID {
				log.Printf("Новый пост из канала: %d", post.Chat.ID)
				if post.MediaGroupID != "" {
					mu.Lock()
					group, exists := mediaGroups[post.MediaGroupID]
					if !exists {
						group = &MediaGroup{
							MediaGroupID: post.MediaGroupID,
							Caption:      post.Caption,
						}
						mediaGroups[post.MediaGroupID] = group
					}

					if post.Photo != nil {
						photo := post.Photo[len(post.Photo)-1]
						group.Media = append(group.Media, tgbotapi.NewInputMediaPhoto(tgbotapi.FileID(photo.FileID)))
					} else if post.Video != nil {
						group.Media = append(group.Media, tgbotapi.NewInputMediaVideo(tgbotapi.FileID(post.Video.FileID)))
					} else if post.Document != nil {
						group.Media = append(group.Media, tgbotapi.NewInputMediaDocument(tgbotapi.FileID(post.Document.FileID)))
					}

					if group.Timer != nil {
						group.Timer.Stop()
					}

					group.Timer = time.AfterFunc(2*time.Second, func() {
						mu.Lock()
						defer mu.Unlock()

						sendMediaGroup(bot, group)

						delete(mediaGroups, group.MediaGroupID)
					})

					mu.Unlock()
				} else {
					if post.Text != "" {
						msg := tgbotapi.NewMessage(targetChannelID, post.Text)
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке текста: %v", err)
						}
					}

					if post.Photo != nil {
						photo := post.Photo[len(post.Photo)-1]
						msg := tgbotapi.NewPhoto(targetChannelID, tgbotapi.FileID(photo.FileID))
						msg.Caption = post.Caption
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке фото: %v", err)
						}
					}

					if post.Video != nil {
						msg := tgbotapi.NewVideo(targetChannelID, tgbotapi.FileID(post.Video.FileID))
						msg.Caption = post.Caption
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке видео: %v", err)
						}
					}

					if post.Document != nil {
						msg := tgbotapi.NewDocument(targetChannelID, tgbotapi.FileID(post.Document.FileID))
						msg.Caption = post.Caption
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке документа: %v", err)
						}
					}

					if post.Voice != nil {
						msg := tgbotapi.NewVoice(targetChannelID, tgbotapi.FileID(post.Voice.FileID))
						msg.Caption = post.Caption
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке голосового сообщения: %v", err)
						}
					}

					if post.VideoNote != nil {
						msg := tgbotapi.NewVideoNote(targetChannelID, 0, tgbotapi.FileID(post.VideoNote.FileID))
						_, err := bot.Send(msg)
						if err != nil {
							log.Printf("Ошибка при отправке видеосообщения: %v", err)
						}
					}
				}
			}
		}
	}
}

func sendMediaGroup(bot *tgbotapi.BotAPI, group *MediaGroup) {
	// Создаем медиагруппу
	mediaGroup := tgbotapi.NewMediaGroup(targetChannelID, group.Media)

	// Добавляем подпись к первому элементу медиагруппы
	if len(mediaGroup.Media) > 0 {
		// Проверяем тип первого элемента и добавляем подпись
		switch firstMedia := mediaGroup.Media[0].(type) {
		case tgbotapi.InputMediaPhoto:
			firstMedia.Caption = group.Caption
			mediaGroup.Media[0] = firstMedia
		case tgbotapi.InputMediaVideo:
			firstMedia.Caption = group.Caption
			mediaGroup.Media[0] = firstMedia
		case tgbotapi.InputMediaDocument:
			firstMedia.Caption = group.Caption
			mediaGroup.Media[0] = firstMedia
		default:
			log.Printf("Тип медиа не поддерживает подпись: %T", firstMedia)
		}
	}

	// Отправляем медиагруппу
	_, err := bot.SendMediaGroup(mediaGroup)
	if err != nil {
		log.Printf("Ошибка при отправке медиагруппы: %v", err)
	} else {
		log.Printf("Медиагруппа отправлена: %s", group.MediaGroupID)
	}
}
