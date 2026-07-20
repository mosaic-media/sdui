package sdui

// Standard-component constructors. These build Nodes of the standard vocabulary
// — the same types the client renders from the shared definition library
// (../definitions). A producer uses these to compose pages; the props keys match
// each definition's contract.
//
// Node type names.
const (
	TypeScreen          = "Screen"
	TypeSection         = "Section"
	TypeStack           = "Stack"
	TypeGrid            = "Grid"
	TypeCarousel        = "Carousel"
	TypeDivider         = "Divider"
	TypePosterCard      = "PosterCard"
	TypeHeroBanner      = "HeroBanner"
	TypeDetailHeader    = "DetailHeader"
	TypeEpisodeRow      = "EpisodeRow"
	TypeSeasonSelector  = "SeasonSelector"
	TypeRelatedRail     = "RelatedRail"
	TypeSourcePicker    = "SourcePicker"
	TypePlaybackBar     = "PlaybackBar"
	TypePersonChip      = "PersonChip"
	TypeGenreTag        = "GenreTag"
	TypeButton          = "Button"
	TypeIconButton      = "IconButton"
	TypeBadge           = "Badge"
	TypeBanner          = "Banner"
	TypeStatusIndicator = "StatusIndicator"
	TypeEmptyState      = "EmptyState"
	TypeSearchBar       = "SearchBar"
	TypeTextField       = "TextField"
	TypeToggle          = "Toggle"
	TypeSelect          = "Select"
	TypeSlider          = "Slider"
	TypeProgressBar     = "ProgressBar"
	TypePagination      = "Pagination"
)

// Option configures a Node after its required fields are set.
type Option func(*Node)

// Prop sets an arbitrary prop. The escape hatch for anything without sugar.
func Prop(key string, val any) Option {
	return func(n *Node) {
		if n.Props == nil {
			n.Props = Props{}
		}
		n.Props[key] = val
	}
}

// ID sets a stable node id.
func ID(id string) Option { return func(n *Node) { n.ID = &id } }

// Act sets the node's primary action.
func Act(a Action) Option { return Prop("action", a) }

// Child appends child nodes.
func Child(nodes ...Node) Option {
	return func(n *Node) { n.Children = append(n.Children, nodes...) }
}

// Slot sets a named slot's nodes.
func Slot(name string, nodes ...Node) Option {
	return func(n *Node) {
		if n.Slots == nil {
			n.Slots = map[string][]Node{}
		}
		n.Slots[name] = append(n.Slots[name], nodes...)
	}
}

// Component is the generic constructor — for a type without a helper, or a
// module's own component.
func Component(nodeType string, opts ...Option) Node {
	n := Node{Type: nodeType}
	for _, o := range opts {
		o(&n)
	}
	return n
}

func build(nodeType string, props Props, opts []Option) Node {
	n := Node{Type: nodeType, Props: props}
	for _, o := range opts {
		o(&n)
	}
	return n
}

// ── containers ───────────────────────────────────────────────────────────────

// Screen is the root of a server-defined page. Title/subtitle via With options.
func Screen(opts ...Option) Node { return build(TypeScreen, nil, opts) }

// Section is a titled band; add a "see all" with Act + Prop("actionLabel", …).
func Section(title string, opts ...Option) Node {
	return build(TypeSection, Props{"title": title}, opts)
}

// Stack arranges children. direction is "horizontal" or "vertical".
func Stack(direction string, gap int, opts ...Option) Node {
	return build(TypeStack, Props{"direction": direction, "gap": gap}, opts)
}

// Grid is a responsive auto-fill grid.
func Grid(opts ...Option) Node { return build(TypeGrid, nil, opts) }

// Carousel is a horizontal snap-scrolling rail.
func Carousel(opts ...Option) Node { return build(TypeCarousel, nil, opts) }

// Divider is a hair rule, optionally labelled via Prop("label", …).
func Divider(opts ...Option) Node { return build(TypeDivider, nil, opts) }

// ── media ────────────────────────────────────────────────────────────────────

// PosterCard renders a Node (work/collection/item). Optional: Subtitle, Poster,
// Progress, BadgeText, Act.
func PosterCard(title, mediaType string, opts ...Option) Node {
	return build(TypePosterCard, Props{"title": title, "mediaType": mediaType}, opts)
}

// HeroBanner is a featured Node. Provide the CTA row via Slot("actions", …).
func HeroBanner(title string, opts ...Option) Node {
	return build(TypeHeroBanner, Props{"title": title}, opts)
}

// DetailHeader renders a Node's metadata. Provide the CTA row via
// Slot("actions", …) and genres via Genres.
func DetailHeader(title string, opts ...Option) Node {
	return build(TypeDetailHeader, Props{"title": title}, opts)
}

// EpisodeRow renders a Part under a series Node.
func EpisodeRow(title string, opts ...Option) Node {
	return build(TypeEpisodeRow, Props{"title": title}, opts)
}

// RelatedRail renders related Nodes; empty when it has no children.
func RelatedRail(title string, opts ...Option) Node {
	return build(TypeRelatedRail, Props{"title": title}, opts)
}

// PersonChip is a cast/crew chip.
func PersonChip(name string, opts ...Option) Node {
	return build(TypePersonChip, Props{"name": name}, opts)
}

// GenreTag is a genre chip; actionable when given Act.
func GenreTag(label string, opts ...Option) Node {
	return build(TypeGenreTag, Props{"label": label}, opts)
}

// SourcePicker surfaces SourceBindings / RemoteLocation stream Parts.
func SourcePicker(sources []Source, opts ...Option) Node {
	return build(TypeSourcePicker, Props{"sources": sources}, opts)
}

// PlaybackBar shows resume / now-playing state.
func PlaybackBar(title string, progress float64, opts ...Option) Node {
	return build(TypePlaybackBar, Props{"title": title, "progress": progress}, opts)
}

// ── controls & feedback ──────────────────────────────────────────────────────

// Button carries an Action. variant is primary/secondary/ghost/danger.
func Button(label, variant string, action Action, opts ...Option) Node {
	return build(TypeButton, Props{"label": label, "variant": variant, "action": action}, opts)
}

// Badge is a small pill. tone is one of the Tone constants.
func Badge(label string, tone Tone, opts ...Option) Node {
	return build(TypeBadge, Props{"label": label, "tone": tone}, opts)
}

// Banner is an inline message.
func Banner(message string, tone Tone, opts ...Option) Node {
	return build(TypeBanner, Props{"message": message, "tone": tone}, opts)
}

// EmptyState is a titled empty placeholder.
func EmptyState(icon, title string, opts ...Option) Node {
	return build(TypeEmptyState, Props{"icon": icon, "title": title}, opts)
}

// ProgressBar renders a determinate 0..1 value.
func ProgressBar(value float64, opts ...Option) Node {
	return build(TypeProgressBar, Props{"value": value}, opts)
}

// ── sugar options ────────────────────────────────────────────────────────────

// Subtitle sets a card/row subtitle.
func Subtitle(s string) Option { return Prop("subtitle", s) }

// Poster sets a poster image URL.
func Poster(url string) Option { return Prop("poster", url) }

// Backdrop sets a hero backdrop image URL.
func Backdrop(url string) Option { return Prop("backdrop", url) }

// Logo sets a HeroBanner's clearlogo/title-treatment image URL; when set it
// renders in place of the text title (ADR 0034).
func Logo(url string) Option { return Prop("logo", url) }

// Progress sets a 0..1 watched fraction.
func Progress(f float64) Option { return Prop("progress", f) }

// BadgeText sets a corner badge on a card.
func BadgeText(s string) Option { return Prop("badge", s) }

// Overview sets a synopsis/overview string.
func Overview(s string) Option { return Prop("overview", s) }

// Meta sets a HeroBanner's meta line (e.g. year, type, rating).
func Meta(items ...string) Option { return Prop("meta", items) }

// Genres sets a DetailHeader's genre list.
func Genres(items ...string) Option { return Prop("genres", items) }

// Source is one entry in a SourcePicker.
type Source struct {
	Label    string  `json:"label"`
	Provider string  `json:"provider,omitempty"`
	Quality  string  `json:"quality,omitempty"`
	Kind     string  `json:"kind,omitempty"`
	Action   *Action `json:"action,omitempty"`
}
