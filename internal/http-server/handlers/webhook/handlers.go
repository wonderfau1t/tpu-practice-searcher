package webhook

import (
	"github.com/go-chi/render"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"net/http"
)

type Storage interface {
	ApprovePhoneNumber(userID int64, phoneNumber string) error
}

func New(bot *tgbotapi.BotAPI, log *slog.Logger, db Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers./Users/wonderfau1t/GolandProjects/tpu-practice-searcher/internal/bot/handlers.go.New"
		log := log.With(slog.String("fn", fn))

		var update tgbotapi.Update
		if err := render.DecodeJSON(r.Body, &update); err != nil {
			log.Error("Failed to decode update: %v", err)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid JSON"})
			return
		}

		if update.Message != nil && update.Message.IsCommand() {
			if update.Message.Command() == "start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Это сервис для поиска практики для студентов ТПУ\nПереходите в сервис по кнопке \"Запустить\"")
				if _, err := bot.Send(msg); err != nil {
					log.Error("Failed to send message: %v", err)
				}
			}
		}

		if update.Message != nil && update.Message.Contact != nil {
			contact := update.Message.Contact
			err := db.ApprovePhoneNumber(contact.UserID, contact.PhoneNumber)
			if err != nil {
				log.Error(err.Error())
			}
			response := "Спасибо! Номер телефона подтвержден"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			if _, err := bot.Send(msg); err != nil {
				log.Error("Failed to send message: %v", err)
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"status": "ok"})
	}
}
