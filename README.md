# Golang/raspberry pi/linux training 

* Raspberry pi 4
* DHT 22 temperature/humidity sensor
* SIM7600E-H 4G HAT

some notes:

build go app

docker build -t raspigoapp
docker run --rm --privileged raspigoapp


On raspbian i run following commands to get docker up if raspberry pi restarts

$ sudo systemctl enable docker.service
$ sudo systemctl enable containerd.service



