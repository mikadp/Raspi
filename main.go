package raspi

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

func sendTelegramMessage(bot *tgbotapi.BotAPI, chatID, int64, message string) error {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	return err
}

func main() {
	tempretureWindow := [windowSize]float32{}
	currentIndex := 0
	currentHighest := float32(-50)
	currentLowest := float32(250)

	pin := "GPIO22" // GPIO pin number

	//set ip db connection
	setupDatabaseConnection()

	//get phone numbers from the db
	phoneNumber, err := getPhoneNumbers()
	if err != nil {
		log.Fatalf("Error getting phone numbers: %v", err)
	}

	// Initilizing Telegram bot
	botToken, err := getTelegramAPI()
	if err != nil {
		log.Fatalf("Error loading Telegram API: %v", err)
	}

	var chatID int64 = 123456
	bot, err := tgbotapi.NewBotAPI(botToken)

	for {
		temp, humidity, err := readDHT22(pin)
		if err != nil {
			log.Fatalf("Error reading temperature and humidity: %v", err)
		}

		tempretureWindow[currentIndex] = temp
		currentIndex = (currentIndex + 1) % windowSize

		oldIndex := (currentIndex - windowSize/2 + windowSize) % windowSize
		diff := tempretureWindow[currentIndex] - tempretureWindow[oldIndex]

		//Comparing temperature changes
		if diff > threshold {
			SendMessage(fmt.Sprintf("Temperature rising fast: %.2f°C", temp), phoneNumber[0])
		} else if diff < threshold {
			SendMessage(fmt.Sprintf("Temperature falling fast: %.2f°C", temp), phoneNumber[0])
		}

		//Update temp to currentHighest and lowest
		currentHighest = temp
		currentLowest = temp

		//update currentLowest if temp is lower than current
		if temp < currentLowest {
			currentLowest = temp
		}
		//update currentHighes if temp is higher
		if temp > currentHighest {
			currentHighest = temp
		}
		//Print temp
		fmt.Printf("Temperature: %.2f°C\n", temp)
		fmt.Printf("Humidity: %.2f%%\n", humidity)

		time.Sleep(1 * time.Minute) //waits 1 minute befor reading again
	}

}
