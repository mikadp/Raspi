package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
)

const (
	dht22PulseTimeout = 200 * time.Microsecond
	dht22MaxDuration  = 250 * time.Microsecond
	windowSize        = 60 //amount of temps to store
	threshold         = 2  // Temp change threshold
)

// add redis for phonenumbers etc
var rdb *redis.Client

func setRedis() {
	// Connect to Redis
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

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
	tempretureWindow := [windowSize]float32{}
	currentIndex := 0
	currentHighest := float32(-50)
	currentLowest := float32(250)

	pin := "GPIO22" // GPIO pin number

	setRedis()
	phoneNumber, err := rdb.Get("phone_number").Result()
	if err != nil {
		log.Fatalf("Error getting phone number from Redis: %v", err)
	}

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
			communication.SendMessage(fmt.Sprintf("Temperature rising fast: %.2f°C", temp), phoneNumber)
		} else if diff < threshold {
			communication.SendMessage(fmt.Sprintf("Temperature falling fast: %.2f°C", temp), phoneNumber)
		}

		//Update temp to currentHighest and lowest
		currentHighest = temp
		currentLowest = temp

		//update currentLowest if temp is lower
		if temp < currentLowest {
			currentLowest = temp
		}
		//update currentHighes if temp is higher
		if temp > currentHighest {
			currentHighest = temp
		}

		fmt.Printf("Temperature: %.2f°C\n", temp)
		fmt.Printf("Humidity: %.2f%%\n", humidity)

		time.Sleep(1 * time.Minute) //waits 1 minute befor reading again
	}

}
