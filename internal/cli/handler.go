package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/karmilg/weather_AI-agent/internal/agent"
)

type Cli struct {
	agent *agent.Agent
}

func NewCli(agentInstance *agent.Agent) *Cli {
	return &Cli{
		agent: agentInstance,
	}
}

func (c *Cli) Run() error {
	fmt.Println("╔══════════════════════════════════════╗")
	fmt.Println("║   🌤️  Погодный AI-агент             ║")
	fmt.Println("║   Введите 'exit' для выхода          ║")
	fmt.Println("╚══════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Примеры: какая погода в Москве? | статистика по Питеру")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("👉 Вы: ")
		if !scanner.Scan() {
			break
		}

		question := strings.TrimSpace(scanner.Text())
		if question == "" {
			continue
		}
		if question == "exit" || question == "выход" {
			fmt.Println("👋 До свидания!")
			break
		}

		fmt.Print("🤔 Думаю... ")
		answer, err := c.agent.Run(question)
		if err != nil {
			fmt.Printf("\n❌ Ошибка: %v\n", err)
			continue
		}

		fmt.Printf("\n🤖 AI: %s\n\n", answer)
	}

	return nil
}
