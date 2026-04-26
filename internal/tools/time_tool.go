package tools

import "time"

type TimeTool struct {
}

func (t *TimeTool) Name() string {
	return "get_current_time"
}

func (t *TimeTool) Description() string {
	return "Возвращает текущее время"
}

func (t *TimeTool) Execute(input string) (string, error) {
	return time.Now().Format("03.04.2006 14:06:05"), nil
}