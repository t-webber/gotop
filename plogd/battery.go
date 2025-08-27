package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

func get_battery_level() int {
	capacity_bytes, err := os.ReadFile("/sys/class/power_supply/BAT1/capacity")
	if err != nil {
		log.Fatalf("Failed to open battery file")
	}

	capacity_string := string(capacity_bytes)
	capacity_parsed := capacity_string[:len(capacity_string)-1]

	battery_level, err := strconv.Atoi(capacity_parsed)
	if err != nil {
		log.Fatalf("Invalid content for battery status file: %s. %s", capacity_parsed, err)
	}

	return battery_level
}

func is_battery_charging() bool {
	data, err := os.ReadFile("/sys/class/power_supply/BAT1/status")
	if err != nil {
		log.Fatalf("Failed to open battery file")
	}

	return string(data) == "Charging\n"
}

func battery(db *db) {
	for {
		battery_level := get_battery_level()
		is_charging := is_battery_charging()

		if !is_charging && battery_level < 20 {
			if err := exec.Command("/bin/notify-send", "Plug your computer").Run(); err != nil {
				log.Fatalf("Failed to notify: %s", err)
			}
			time.Sleep(time.Minute)
			if err := syscall.Exec("/bin/sudo", []string{"sudo", "systemctl", "suspend"}, []string{}); err != nil {
				log.Fatalf("Failed to suspend device: %s", err)
			}
		}

		db.mutex.Lock()
		if _, err := db.db.Exec("INSERT INTO battery(time, level, charging) VALUES(?, ?, ?)", time.Now(), battery_level, is_charging); err != nil {
			log.Fatalf("Failed to insert into battery: %s", err)
		}

		db.mutex.Unlock()

		time.Sleep(time.Minute)
	}
}
