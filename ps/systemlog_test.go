package ps

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

var latestSystemLogEvents = []byte(`
# Скрипт: получить события из журнал "System" в формате JSON за последние сутки

$filter = @{
    "LogName"   = "System";
    "StartTime" = (Get-Date).AddHours(-24)
}

$queryItems = [System.Collections.ArrayList]@(
    @{Name = "id"; Expression = { $_.id }}
    @{Name = "datetime"; Expression = { (Get-Date $_.TimeCreated.datetime).ToUniversalTime().ToString("yyyy-MM-dd'T'HH:mm:ssZ")}}
    @{Name = "providerName"; Expression = { $_.ProviderName }}
    @{Name = "levelDisplayName"; Expression = { $_.LevelDisplayName }}
    @{Name = "message"; Expression = { $_.Message }}
)
Get-WinEvent -FilterHashtable $filter | Select-Object $queryItems | ConvertTo-Json

`)

type systemLog struct {
	ID               uint64    `json:"id"`
	DateTime         time.Time `json:"datetime"`
	ProviderName     string    `json:"providerName"`
	LevelDisplayName string    `json:"levelDisplayName"`
	Message          string    `json:"message"`
}

func (e *systemLog) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("ID: %v\n", e.ID))
	builder.WriteString(fmt.Sprintf("DateTime: %v\n", e.DateTime))
	builder.WriteString(fmt.Sprintf("ProviderName: %v\n", e.ProviderName))
	builder.WriteString(fmt.Sprintf("LevelDisplayName: %v\n", e.LevelDisplayName))
	builder.WriteString(fmt.Sprintf("Message: %v\n", e.Message))
	return builder.String()
}

func TestRunSystemLogScript(t *testing.T) {
	result, err := Run(latestSystemLogEvents)
	if err != nil {
		t.Fatal(err)
	}

	var entries []*systemLog
	if err := json.Unmarshal(result, &entries); err != nil {
		t.Fatal(err)
	}

	fmt.Println(entries)
}
