package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"todo_list/config"
	"todo_list/internal/infrastructure/adapters/logger"

	"golang.org/x/exp/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	codes      = make(map[string]string)
	codesMutex sync.Mutex
)

func generateCode() string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	code := r.Intn(999999)

	return fmt.Sprintf("%06d", code)
}

func sendCode(chatID int64, code, telegramBotToken string) error {
	bot, err := tgbotapi.NewBotAPI(telegramBotToken)
	if err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(chatID, "Your code for enter: "+code)
	msg.ReplyMarkup = createKeyboard()
	_, err = bot.Send(msg)

	if err != nil {
		return fmt.Errorf("error send code: %v", err)
	}

	return nil
}

func requestCode(chatID int64, username, telegramBotToken string) error {
	code := generateCode()
	codesMutex.Lock()
	codes[username] = code
	codesMutex.Unlock()

	err := sendCode(chatID, code, telegramBotToken)
	if err != nil {
		return fmt.Errorf("error send code: %v", err)
	}

	return nil
}

func createKeyboard() tgbotapi.ReplyKeyboardMarkup {
	btn := tgbotapi.NewKeyboardButton("Get code")
	keyboard := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(btn))
	keyboard.OneTimeKeyboard = false
	return keyboard
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error load config: %v", err)
	}

	logger := logger.New(cfg)

	go server(cfg, logger)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("error create bot: %v", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			username := update.Message.From.UserName
			switch text := update.Message.Text; text {

			case "/start":
				welcomeMsg := tgbotapi.NewMessage(chatID, "Welcome! Press the button to receive the code.")
				welcomeMsg.ReplyMarkup = createKeyboard()
				if _, err := bot.Send(welcomeMsg); err != nil {
					logger.Error("Error sending the welcome message: ",
						slog.String("error", err.Error()))
				}

			case "Get code":
				if err := requestCode(chatID, username, cfg.TelegramBotToken); err != nil {
					logger.Error("Error requesting the code: ",
						slog.String("error", err.Error()))
				}

			}
		}
	}
}

func server(cfg *config.Config, logger *slog.Logger) {
	http.HandleFunc("/code/", requestCodeHandler)

	address := fmt.Sprintf(":%s", cfg.TwoFAPort)

	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Error("Error starting the server: ",
			slog.String("error", err.Error()))
	}

}

func requestCodeHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/code/"):]

	slog.Info("Request code", "username", username)

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)

		return
	}

	codesMutex.Lock()
	code, exists := codes[username]
	codesMutex.Unlock()

	if !exists {
		http.Error(w, "Code not found", http.StatusNotFound)

		return
	}

	response := map[string]string{"username": username, "code": code}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)

		return
	}

}
