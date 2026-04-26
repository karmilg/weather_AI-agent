package agent

import (
	"log"
	"strings"

	"github.com/karmilg/weather_AI-agent/internal/llm"
	"github.com/karmilg/weather_AI-agent/internal/tools"
)

type Agent struct {
    llmClient *llm.Client
    tools     map[string]tools.Tool
}

func NewAgent(llmClient *llm.Client, toolList []tools.Tool) *Agent {
    toolMap := make(map[string]tools.Tool)
    for _, t := range toolList {
        toolMap[t.Name()] = t
        log.Printf("✅ Инструмент зарегистрирован: %s", t.Name())
    }
    return &Agent{
        llmClient: llmClient,
        tools:     toolMap,
    }
}

func (a *Agent) Run(userQuestion string) (string, error) {
    log.Printf("📨 Вопрос пользователя: %s", userQuestion)

    fullPrompt := GetSystemPrompt() + "\n\nВопрос: " + userQuestion + "\n\nОТВЕТЬ ТОЛЬКО НА РУССКОМ ЯЗЫКЕ:"

    log.Printf("📤 Отправка в LLM: %s", fullPrompt)

    response, err := a.llmClient.Generate(fullPrompt)
    if err != nil {
        log.Printf("❌ Ошибка LLM: %v", err)
        return "", err
    }

    log.Printf("📥 Ответ от LLM: %s", response)

    response = strings.TrimSpace(response)

    if strings.HasPrefix(response, "TOOL:") {
        toolPart := strings.TrimPrefix(response, "TOOL:")
        toolPart = strings.TrimSpace(toolPart)

        parts := strings.SplitN(toolPart, ":", 2)
        if len(parts) == 0 {
            log.Printf("❌ Неверный формат TOOL: %s", toolPart)
            return "Извините, произошла ошибка. Попробуйте ещё раз.", nil
        }

        toolName := strings.TrimSpace(parts[0])
        toolInput := ""
        if len(parts) > 1 {
            toolInput = strings.TrimSpace(parts[1])
        }

        log.Printf("🔧 Вызов инструмента: %s с параметром: %s", toolName, toolInput)

        if tool, exists := a.tools[toolName]; exists {
            result, err := tool.Execute(toolInput)
            if err != nil {
                log.Printf("❌ Ошибка выполнения инструмента %s: %v", toolName, err)
                return "Извините, не удалось выполнить запрос. Попробуйте ещё раз.", nil
            }
            log.Printf("✅ Результат инструмента %s: %s", toolName, result)
            return result, nil
        }

        log.Printf("❌ Инструмент не найден: %s", toolName)
        log.Printf("🔍 Доступные инструменты: %v", getToolNames(a.tools))
        return "Извините, такой функции нет. Доступны: погода, время, история.", nil
    }

    log.Printf("💬 Обычный ответ от LLM (не TOOL)")
    return response, nil
}

func getToolNames(tools map[string]tools.Tool) []string {
    names := make([]string, 0, len(tools))
    for name := range tools {
        names = append(names, name)
    }
    return names
}
