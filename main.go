package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// This is based off of Brodie Robertson's polybattery script
// Dependencies: acpi command

func main() {
	cmd := exec.Command("acpi", "-b")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lines := strings.Split(string(stdout), "\n")

	var chosen_line string
	for _, line := range lines {
		if !strings.Contains(line, "rate information unavailable") {
			chosen_line = line
			break
		}
	}

	args := strings.Split(chosen_line, " ")

	mode := args[2]
	percent := args[3]

	if mode[len(mode)-1] == ',' {
		mode = mode[:len(mode)-1]
	}

	var symbol string
	switch mode {
	case "Discharging":
		symbol = "⚡"
	case "Charging":
		symbol = "🔌"
	case "Unknown":
		symbol = "🔋?"
	case "Not": // When the battery status is "Not charging,"
		if args[3] == "charging," {
			percent = args[4]
		}
		symbol = "🔋?"
	default:
		symbol = "🔋"
	}

	if percent[len(percent)-1] == ',' {
		percent = percent[:len(percent)-1]
	}

	fmt.Print(symbol)
	fmt.Print(" ")

	// Assert
	if percent[len(percent)-1] != '%' {
		panic(fmt.Sprintf("percent is invalidly formatted: %s", percent))
	}

	raw_percent := percent[:len(percent)-1]
	raw_percent_int, err := strconv.Atoi(raw_percent)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if raw_percent_int < 20 {
		// These tell polybar to change the color
		fmt.Print("%{F#ed0b0b}")
	} else if raw_percent_int < 50 {
		fmt.Print("%{F#f2e421}")
	}

	fmt.Print(percent)

	if len(args) >= 5 && mode != "Not" {
		remaining := args[4]

		fmt.Printf(" (%s)", remaining)
	}
}
