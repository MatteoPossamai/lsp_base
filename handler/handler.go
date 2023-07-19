package handler

import (
	"encoding/json"
	"log"

	"lsp.com/server/external"
	"lsp.com/server/types"
)

func toRawJson(input interface{}) (json.RawMessage, error) {
	jsonified_data, err := json.Marshal(input)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return json.RawMessage(string(jsonified_data)), nil
}

func SendCapabilities() (json.RawMessage, error) {
	capability_list := []string{"textDocument/publishDiagnostics", "..."}
	return toRawJson(capability_list)
}

func SendDiagnostics(consoles []types.Console) (json.RawMessage, error) {
	diagnostics, err := generateDiagnostics(consoles)
	if err != nil {
		return nil, err
	}
	return toRawJson(diagnostics)
}

func generateDiagnostics(consoles []types.Console) ([]types.Diagnostic, error) {
	var diagnostics []types.Diagnostic
	for console := range consoles {
		diagnostics = append(diagnostics, types.Diagnostic{
			Code:    types.Hint,
			Console: consoles[console],
			Message: "test",
			Range: types.Range{
				Start: 0,
				End:   len(consoles[console].Content),
			},
		})
	}
	return diagnostics, nil
}

func SendCompletion(console types.Console, pointer int) (json.RawMessage, error) {
	completion_list := external.Suggest(console.Content, pointer)
	return toRawJson(completion_list)
}
