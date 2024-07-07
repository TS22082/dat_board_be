package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
)

func Logging(c *fiber.Ctx) error {
	// Define color functions
	methodColor := color.New(color.FgCyan).SprintFunc()
	urlColor := color.New(color.FgGreen).SprintFunc()
	ipColor := color.New(color.FgYellow).SprintFunc()
	queryColor := color.New(color.FgMagenta).SprintFunc()
	statusColor := color.New(color.FgRed).SprintFunc()
	durationColor := color.New(color.FgBlue).SprintFunc()

	start := time.Now()
	err := c.Next()
	duration := time.Since(start)

	// Log with colors in a table format
	log.Println("-------------------------------------------------------------")
	log.Printf("| %-15s | %-40s ", "Request", fmt.Sprintf("%s %s", methodColor(c.Method()), urlColor(c.OriginalURL())))
	log.Printf("| %-15s | %-40s ", "IP", ipColor(c.IP()))
	log.Printf("| %-15s | %-40s ", "Queries", queryColor(c.Queries()))
	log.Printf("| %-15s | %-40s ", "Response Status", statusColor(c.Response().StatusCode()))
	log.Printf("| %-15s | %-40s ", "Duration", durationColor(duration))
	log.Println("-------------------------------------------------------------")

	return err
}
