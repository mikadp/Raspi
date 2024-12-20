package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	serial "go.bug.st/serial"
)

// handle Telegram messages
func handleTelegramMessage(bot *tgbotapi.BotAPI) error {
	//var chatID int64 = 123456 //telegram chat id
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		return fmt.Errorf("error getting Telegram updates channel: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		messageText := strings.ToLower(update.Message.Text)

		if messageText == "temp" || messageText == "status" {
			temperature, humidity, err := readDHT22("GPIO22")
			if err != nil {
				log.Printf("Error reading temperature: %v", err)
				continue
			}

			message := fmt.Sprintf("Current temperature: %.2f°C and humidity: %.2f", temperature, humidity)
			if err := sendTelegramMessage(bot, chatID, message); err != nil {
				log.Printf("Error sending Telegram message: %v", err)
			}
		}
	}

	return nil
}

func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	return err
}

func SendMessage(message string, phoneNumber string) {
	//Serial port opening
	mode := &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	port, err := serial.Open("/dev/ttyS0", mode)
	if err != nil {
		fmt.Println("Error opening serial port: ", err)
		return
	}
	defer port.Close()
	time.Sleep(5 * time.Second)
	port.Write([]byte("AT+CMGF=1\r\n"))
	time.Sleep(500 * time.Millisecond)
	port.Write([]byte(fmt.Sprintf("AT+CMGS=\"%s\"\r\n", phoneNumber)))
	time.Sleep(500 * time.Millisecond)
	port.Write([]byte(message))
	port.Write([]byte{0x1A})
}
