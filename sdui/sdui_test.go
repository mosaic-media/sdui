package sdui_test

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/mosaic-media/sdui/sdui"
)

// homeScreen is the kind of tree the Platform's emit-side will build.
func homeScreen() sdui.Node {
	return sdui.Screen(
		sdui.Child(
			sdui.HeroBanner("Spirited Away",
				sdui.Meta("2001", "Anime Film", "PG"),
				sdui.Overview("A young girl wanders into a world of spirits."),
				sdui.Slot("actions",
					sdui.Button("Play", "primary", sdui.Play("demo-part")),
					sdui.Button("Details", "secondary", sdui.Navigate("detail", map[string]any{"title": "Spirited Away"})),
				),
			),
			sdui.Section("Continue watching",
				sdui.Child(sdui.Carousel(sdui.Child(
					sdui.PosterCard("Cowboy Bebop", "Anime Series",
						sdui.Subtitle("S1 · E12"), sdui.Progress(0.6), sdui.BadgeText("12 min left"),
						sdui.Act(sdui.Navigate("detail", map[string]any{"title": "Cowboy Bebop"}))),
					sdui.PosterCard("Dune", "Film", sdui.Progress(0.75)),
				))),
			),
		),
	)
}

// find walks the tree (children and slots) for the first node of the given type.
func find(n sdui.Node, typ string) sdui.Node {
	if n == nil {
		return nil
	}
	if n.GetType() == typ {
		return n
	}
	for _, c := range n.GetChildren() {
		if got := find(c, typ); got != nil {
			return got
		}
	}
	for _, list := range n.GetSlots() {
		for _, c := range list.GetNodes() {
			if got := find(c, typ); got != nil {
				return got
			}
		}
	}
	return nil
}

func TestHomeScreenBuilds(t *testing.T) {
	root := homeScreen()
	if root.GetType() != "Screen" {
		t.Fatalf("root type = %q, want Screen", root.GetType())
	}
	if len(root.GetChildren()) != 2 {
		t.Fatalf("root children = %d, want 2", len(root.GetChildren()))
	}

	// The hero carries its title in props and two buttons in the actions slot.
	hero := find(root, "HeroBanner")
	if hero == nil {
		t.Fatal("no HeroBanner")
	}
	if hero.GetProps().AsMap()["title"] != "Spirited Away" {
		t.Fatalf("hero title = %v", hero.GetProps().AsMap()["title"])
	}
	actions := hero.GetSlots()["actions"]
	if actions == nil || len(actions.GetNodes()) != 2 {
		t.Fatalf("hero actions slot = %v, want 2 nodes", actions)
	}

	// The first button's action is a playPart carrying the part id — an Action
	// riding inside the open props bag.
	play := actions.GetNodes()[0].GetProps().AsMap()["action"].(map[string]any)
	if play["kind"] != "playPart" || play["partId"] != "demo-part" {
		t.Fatalf("play action = %v, want playPart/demo-part", play)
	}

	if find(root, "Carousel") == nil {
		t.Fatal("no Carousel")
	}

	// The first PosterCard (Cowboy Bebop) carries a 0.6 progress in props.
	card := find(root, "PosterCard")
	if card == nil {
		t.Fatal("no PosterCard")
	}
	if got := card.GetProps().AsMap()["progress"]; got != 0.6 {
		t.Fatalf("card progress = %v, want 0.6", got)
	}
}

func TestActionsAreCleanPerKind(t *testing.T) {
	cases := map[string]sdui.Action{
		`{"kind":"navigate","screen":"home"}`:              sdui.Navigate("home", nil),
		`{"kind":"playPart","partId":"p1"}`:                sdui.Play("p1"),
		`{"kind":"toast","message":"hi","tone":"success"}`: sdui.Toast("hi", sdui.ToneSuccess),
		`{"kind":"invoke","mutation":"importContent"}`:     sdui.Invoke("importContent", nil),
		`{"kind":"back"}`:                                  sdui.Back(),
	}
	for want, a := range cases {
		b, err := json.Marshal(a)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		if string(b) != want {
			t.Errorf("action marshalled to %s, want %s", b, want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	b, err := protojson.Marshal(homeScreen())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := sdui.Component("") // an empty node to unmarshal into
	if err := protojson.Unmarshal(b, out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.GetType() != "Screen" || len(out.GetChildren()) != 2 {
		t.Fatalf("round-trip lost structure: type=%q children=%d", out.GetType(), len(out.GetChildren()))
	}
}
