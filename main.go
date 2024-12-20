package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const (
	dht22PulseTimeout = 200 * time.Microsecond
	dht22MaxDuration  = 250 * time.Microsecond
	windowSize        = 60 //amount of temps to store
	threshold         = 5  // Temp change threshold
)

func readDHT22(pin string) (float32, float32, error) {
	p := gpioreg.ByName(pin)
	if p == nil {
		return 0, 0, errors.New("failed to find GPIO pin")
	}

	// Send start signal to the sensor.
	p.Out(gpio.Low)
	time.Sleep(1 * time.Millisecond)
	p.Out(gpio.High)
	time.Sleep(40 * time.Microsecond)

	// Read response signal from the sensor.
	p.In(gpio.PullUp, gpio.NoEdge)
	bitStream := make([]int, 0, 40)
	for i := 0; i < 40; i++ {
		timeout := dht22MaxDuration
		for p.Read() == gpio.Low {
			timeout -= dht22PulseTimeout
			if timeout < 0 {
				return 0, 0, errors.New("timeout waiting for response signal")
			}
			time.Sleep(dht22PulseTimeout)
		}
		duration := time.Duration(0)
		for p.Read() == gpio.High {
			duration += dht22PulseTimeout
			if duration > dht22MaxDuration {
				return 0, 0, errors.New("response signal duration out of range")
			}
			time.Sleep(dht22PulseTimeout)
		}
		bit := 0
		if duration > 30*time.Microsecond {
			bit = 1
		}
		bitStream = append(bitStream, bit)
	}

	// Convert bit stream to data values.
	humidity := int(bitStream[0])*256 + int(bitStream[1])
	temperature := int(bitStream[2]&0x7F)*256 + int(bitStream[3])
	if bitStream[2]&0x80 != 0 {
		temperature = -temperature
	}
	checksum := bitStream[0] + bitStream[1] + bitStream[2] + bitStream[3]
	if checksum&0xFF != bitStream[4] {
		return 0, 0, errors.New("checksum error")
	}

	return float32(temperature) / 10, float32(humidity) / 10, nil
}

func main() {
	setupDatabaseConnection()
	phoneNumbers, err := getPhoneNumbers()
	if err != nil {
		log.Fatalf("Error getting phone numbers: %v", err)
	}

	botToken, err := getTelegramAPI()
	if err != nil {
		log.Fatalf("Error loading Telegram API: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Error initializing Telegram bot: %v", err)
	}

	go func() {
		err = handleTelegramMessage(bot)
		if err != nil {
			log.Printf("Error handling Telegram messages: %v", err)
		}
	}()

	pin := "GPIO22"
	temperatureWindow := make([]float32, windowSize)
	currentIndex := 0

	for {
		temp, _, err := readDHT22(pin)
		if err != nil {
			log.Printf("Error reading sensor: %v", err)
			continue
		}

		temperatureWindow[currentIndex] = temp
		currentIndex = (currentIndex + 1) % windowSize

		if currentIndex >= threshold {
			diff := temperatureWindow[currentIndex-1] - temperatureWindow[currentIndex-threshold]
			if diff > float32(threshold) {
				SendMessage(fmt.Sprintf("Temperature rising fast: %.2f°C", temp), phoneNumbers[0])
			} else if diff < -float32(threshold) {
				SendMessage(fmt.Sprintf("Temperature falling fast: %.2f°C", temp), phoneNumbers[0])
			}
		}

		fmt.Printf("Temperature: %.2f°C\n", temp)
		time.Sleep(1 * time.Minute)
	}
}
