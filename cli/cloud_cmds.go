package cli

import (
	"github.com/ShawnPana/browser/core"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func cmdCloud(args []string) {
	if len(args) == 0 {
		Fatal("usage: browser cloud <login|logout|METHOD|poll|--help> ...")
	}

	switch args[0] {
	case "login":
		cloudLogin(args[1:])
	case "logout":
		cloudLogout(args[1:])
	case "poll":
		cloudPoll(args[1:])
	case "--help", "-h", "help":
		cloudHelp()
	default:
		cloudREST(args)
	}
}

func cloudLogin(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser cloud login <api-key>")
	}
	ctx := core.NewContext()
	cfg := core.CloudConfig{APIKey: args[0]}
	if err := core.SaveCloudConfig(ctx, cfg); err != nil {
		Fatal("failed to save cloud config: %v", err)
	}
	fmt.Println("logged in")
}

func cloudLogout(args []string) {
	ctx := core.NewContext()
	if err := core.RemoveCloudConfig(ctx); err != nil {
		// Ignore if file doesn't exist
	}
	fmt.Println("logged out")
}

func cloudREST(args []string) {
	if len(args) < 2 {
		Fatal("usage: browser cloud <METHOD> <path> [json-body]")
	}

	method := strings.ToUpper(args[0])
	path := args[1]
	var bodyStr string
	if len(args) > 2 {
		bodyStr = strings.Join(args[2:], " ")
	}

	ctx := core.NewContext()
	cfg := core.LoadCloudConfig(ctx)
	if cfg.APIKey == "" {
		Fatal("no API key configured (use 'browser cloud login <key>' or set BROWSER_USE_API_KEY)")
	}

	url := core.CloudBaseURL + path
	if !strings.HasPrefix(path, "/") {
		url = core.CloudBaseURL + "/" + path
	}

	var body io.Reader
	if bodyStr != "" {
		body = strings.NewReader(bodyStr)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		Fatal("failed to create request: %v", err)
	}
	req.Header.Set("X-Browser-Use-API-Key", cfg.APIKey)
	if bodyStr != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		Fatal("request failed: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		Fatal("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP %d\n", resp.StatusCode)
		PrintJSON(respBody)
		os.Exit(2)
	}

	PrintJSON(respBody)
}

func cloudPoll(args []string) {
	if len(args) < 1 {
		Fatal("usage: browser cloud poll <task-id>")
	}
	taskID := args[0]

	ctx := core.NewContext()
	cfg := core.LoadCloudConfig(ctx)
	if cfg.APIKey == "" {
		Fatal("no API key configured")
	}

	url := core.CloudBaseURL + "/tasks/" + taskID

	client := &http.Client{Timeout: 10 * time.Second}

	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			Fatal("%v", err)
		}
		req.Header.Set("X-Browser-Use-API-Key", cfg.APIKey)

		resp, err := client.Do(req)
		if err != nil {
			Fatal("request failed: %v", err)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var status struct {
			Status  string `json:"status"`
			CostStr string `json:"cost"`
		}
		json.Unmarshal(body, &status)

		cost, _ := strconv.ParseFloat(status.CostStr, 64)
		fmt.Fprintf(os.Stderr, "\rstatus: %s  cost: $%.4f", status.Status, cost)

		if status.Status == "finished" || status.Status == "stopped" || status.Status == "failed" {
			fmt.Fprintln(os.Stderr)
			PrintJSON(body)
			if status.Status == "failed" {
				os.Exit(2)
			}
			return
		}

		time.Sleep(2 * time.Second)
	}
}

func cloudHelp() {
	ctx := core.NewContext()
	cfg := core.LoadCloudConfig(ctx)

	// Try to fetch OpenAPI spec
	if cfg.APIKey != "" {
		if help := fetchOpenAPIHelp(cfg.APIKey); help != "" {
			fmt.Print(help)
			return
		}
	}

	// Static fallback
	fmt.Print(`Browser Use Cloud API

Usage: browser cloud <METHOD> <path> [json-body]

Examples:
  browser cloud GET /browsers
  browser cloud POST /browsers '{"headless":true}'
  browser cloud GET /browsers/<id>
  browser cloud PATCH /browsers/<id> '{"action":"stop"}'
  browser cloud GET /tasks/<id>/status
  browser cloud POST /tasks '{"url":"https://example.com","task":"..."}'

Management:
  browser cloud login <api-key>     Save API key
  browser cloud logout              Remove API key
  browser cloud poll <task-id>      Poll task until completion

API docs: https://docs.cloud.browser-use.com/api-reference/api-v-2
`)
}

func fetchOpenAPIHelp(apiKey string) string {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("GET", "https://api.browser-use.com/api/v2/openapi.json", nil)
	req.Header.Set("X-Browser-Use-API-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return formatOpenAPIHelp(body)
}

func formatOpenAPIHelp(spec []byte) string {
	var doc struct {
		Paths map[string]map[string]struct {
			Summary    string   `json:"summary"`
			Tags       []string `json:"tags"`
			Parameters []struct {
				Name     string `json:"name"`
				In       string `json:"in"`
				Required bool   `json:"required"`
				Schema   struct {
					Type string `json:"type"`
				} `json:"schema"`
			} `json:"parameters"`
		} `json:"paths"`
	}

	if err := json.Unmarshal(spec, &doc); err != nil {
		return ""
	}

	type endpoint struct {
		method, path, summary string
		tag                   string
		params                []string
	}

	var endpoints []endpoint
	for path, methods := range doc.Paths {
		for method, op := range methods {
			tag := "Other"
			if len(op.Tags) > 0 {
				tag = op.Tags[0]
			}
			var params []string
			for _, p := range op.Parameters {
				marker := ""
				if p.Required {
					marker = "*"
				}
				params = append(params, fmt.Sprintf("%s%s (%s, %s)", p.Name, marker, p.In, p.Schema.Type))
			}
			endpoints = append(endpoints, endpoint{
				method:  strings.ToUpper(method),
				path:    path,
				summary: op.Summary,
				tag:     tag,
				params:  params,
			})
		}
	}

	// Group by tag
	groups := make(map[string][]endpoint)
	var tagOrder []string
	for _, ep := range endpoints {
		if _, exists := groups[ep.tag]; !exists {
			tagOrder = append(tagOrder, ep.tag)
		}
		groups[ep.tag] = append(groups[ep.tag], ep)
	}

	var b strings.Builder
	b.WriteString("Browser Use Cloud API Endpoints\n\n")

	for _, tag := range tagOrder {
		eps := groups[tag]
		fmt.Fprintf(&b, "## %s\n", tag)
		for _, ep := range eps {
			fmt.Fprintf(&b, "  %-7s %s", ep.method, ep.path)
			if ep.summary != "" {
				fmt.Fprintf(&b, "  — %s", ep.summary)
			}
			b.WriteString("\n")
			for _, p := range ep.params {
				fmt.Fprintf(&b, "          param: %s\n", p)
			}
		}
		b.WriteString("\n")
	}

	return b.String()
}
