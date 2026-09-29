package wirepod_ttr

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kercre123/wire-pod/chipper/pkg/vars"
)

func TestAzureLunaRequestOmitsRejectedParams(t *testing.T) {
	prev := vars.APIConfig
	t.Cleanup(func() { vars.APIConfig = prev })

	vars.APIConfig.Knowledge.Provider = "custom"
	vars.APIConfig.Knowledge.Endpoint = "https://ksgptsweden.cognitiveservices.azure.com/openai/responses?api-version=2025-04-01-preview"
	vars.APIConfig.Knowledge.Model = "gpt-6-luna"
	vars.APIConfig.STT.Language = "de-DE"
	vars.APIConfig.Knowledge.OpenAIPrompt = ""

	req := CreateAIReq("Wer war Albert Einstein?", "esn", false, true)
	raw, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	body := string(raw)
	for _, bad := range []string{"max_tokens", "temperature", "top_p", "frequency_penalty", "presence_penalty", "reasoning_effort", "verbosity"} {
		if strings.Contains(body, `"`+bad+`"`) {
			t.Fatalf("Azure body contains %s: %s", bad, body)
		}
	}
	if !strings.Contains(body, `"max_completion_tokens"`) || !strings.Contains(body, `"gpt-6-luna"`) {
		t.Fatalf("missing completion budget or model: %s", body)
	}
	if !strings.Contains(body, "ein oder zwei kurzen Saetzen") {
		t.Fatalf("german prompt missing: %s", body)
	}

	got := normalizeChatEndpoint(vars.APIConfig.Knowledge.Endpoint)
	want := "https://ksgptsweden.cognitiveservices.azure.com/openai/v1"
	if got != want {
		t.Fatalf("endpoint %s, want %s", got, want)
	}
}

func TestSpokenDeltaStripsActionTagsWhenCommandsOff(t *testing.T) {
	prev := vars.APIConfig.Knowledge.CommandsEnable
	t.Cleanup(func() { vars.APIConfig.Knowledge.CommandsEnable = prev })
	vars.APIConfig.Knowledge.CommandsEnable = false
	got := spokenDelta("Albert Einstein. {{playAnimationWI||thinking}}")
	if strings.Contains(got, "{{") || strings.Contains(got, "playAnimation") {
		t.Fatalf("tag leaked into speech: %q", got)
	}
	if !strings.Contains(got, "Albert Einstein") {
		t.Fatalf("speech lost: %q", got)
	}
}

func TestRemoveSpecialCharactersKeepsGermanUmlauts(t *testing.T) {
	got := removeSpecialCharacters("für berühmter Relativitätstheorie Größe")
	want := "für berühmter Relativitätstheorie Größe"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
