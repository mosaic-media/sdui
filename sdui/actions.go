package sdui

// Action constructors. The Action uses pointer fields for optionals (so the JSON
// omits absent ones); these constructors hide that.

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
)

func strp(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Navigate pushes another server-defined screen.
func Navigate(screen string, params map[string]any) Action {
	return Action{Kind: KindNavigate, Screen: strp(screen), Params: params}
}

// Back pops the client's navigation stack.
func Back() Action { return Action{Kind: KindBack} }

// OpenURL opens an external URL (the client validates the scheme).
func OpenURL(url string) Action { return Action{Kind: KindOpenURL, URL: strp(url)} }

// Invoke runs a Platform mutation by name.
func Invoke(mutation string, input map[string]any) Action {
	return Action{Kind: KindInvoke, Mutation: strp(mutation), Input: input}
}

// Query runs a Platform query, optionally refreshing a named region.
func Query(query string, variables map[string]any, into string) Action {
	return Action{Kind: KindQuery, Query: strp(query), Variables: variables, Into: strp(into)}
}

// OpenOverlay presents a node as a modal/sheet/drawer. The node rides inside the
// action's open bag, so it is rendered to its JSON object form.
func OpenOverlay(surface Surface, node Node) Action {
	return Action{Kind: KindOpenOverlay, Surface: &surface, Node: nodeToMap(node)}
}

// CloseOverlay dismisses the topmost overlay.
func CloseOverlay() Action { return Action{Kind: KindCloseOverlay} }

// Play asks the client to resolve and play a content Part.
func Play(partID string) Action { return Action{Kind: KindPlayPart, PartID: strp(partID)} }

// Toast shows a transient message.
func Toast(message string, tone Tone) Action {
	a := Action{Kind: KindToast, Message: strp(message)}
	if tone != "" {
		a.Tone = &tone
	}
	return a
}

// Sequence runs several actions in order.
func Sequence(actions ...Action) Action {
	return Action{Kind: KindSequence, Actions: actions}
}

// nodeToMap renders a node to the JSON object form that rides inside an action's
// open bag (OpenOverlay). Nodes embedded in an action are rare; when absent this
// is never called.
func nodeToMap(node Node) map[string]any {
	if node == nil {
		return nil
	}
	b, err := protojson.Marshal(node)
	if err != nil {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}
