package sdui_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/mosaic-media/sdui/sdui"
)

// These tests keep schema/sdui.schema.json honest for the parts still expressed
// as JSON: the Action envelope (which rides inside the open props bag) and the
// standard definition files. The built UINode tree itself is now the protobuf
// mosaic.sdui.v1.UINode (ADR 0044) — a typed message, not JSON-Schema-shaped — so
// it is exercised structurally in sdui_test.go rather than validated here.

func compile(t *testing.T, ptr string) *jsonschema.Schema {
	t.Helper()
	sch, err := jsonschema.Compile("../schema/sdui.schema.json#/$defs/" + ptr)
	if err != nil {
		t.Fatalf("compile %s: %v", ptr, err)
	}
	return sch
}

// asAny round-trips a value through JSON into the generic form the validator wants.
func asAny(t *testing.T, v any) any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestActionsConformToSchema(t *testing.T) {
	sch := compile(t, "Action")
	actions := []sdui.Action{
		sdui.Navigate("home", map[string]any{"q": "x"}),
		sdui.Play("p1"),
		sdui.Toast("hi", sdui.ToneSuccess),
		sdui.Invoke("importContent", map[string]any{"id": 1}),
		sdui.OpenOverlay(sdui.SurfaceSheet, sdui.Screen()),
		sdui.Sequence(sdui.Back(), sdui.Toast("done", sdui.ToneInfo)),
	}
	for i, a := range actions {
		if err := sch.Validate(asAny(t, a)); err != nil {
			t.Errorf("action %d does not conform:\n%v", i, err)
		}
	}
}

func TestStandardDefinitionsConformToSchema(t *testing.T) {
	sch := compile(t, "ComponentDefinition")
	files, err := filepath.Glob("../definitions/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no definition files found: %v", err)
	}
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		var inst any
		if err := json.Unmarshal(data, &inst); err != nil {
			t.Fatalf("parse %s: %v", f, err)
		}
		if err := sch.Validate(inst); err != nil {
			t.Errorf("%s does not conform to the ComponentDefinition schema:\n%v", filepath.Base(f), err)
		}
	}
}
