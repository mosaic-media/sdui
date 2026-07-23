// Code generated from schema/sdui.schema.json by quicktype. DO NOT EDIT.
/**
 * The single source of truth for the Mosaic Server-Driven-UI contract. Language bindings
 * (Go, TypeScript, Dart) are GENERATED from this file — do not hand-edit them. The root
 * object exists only so a generator reaches every top-level type; the useful types are in
 * $defs.
 */

/**
 * A declarative behaviour envelope. Data, never code — the client interprets the kind. Each
 * kind uses a subset of the fields.
 */
export interface Action {
    actions?:  Action[];
    input?:    { [key: string]: any };
    kind:      ActionKind;
    message?:  string;
    mutation?: string;
    node?:     UINode;
    nodeId?:   string;
    params?:   { [key: string]: any };
    partId?:   string;
    screen?:   string;
    surface?:  Surface;
    tone?:     Tone;
    url?:      string;
}

export enum ActionKind {
    Back = "back",
    CloseOverlay = "closeOverlay",
    Invoke = "invoke",
    Navigate = "navigate",
    OpenOverlay = "openOverlay",
    OpenURL = "openUrl",
    PlayPart = "playPart",
    Sequence = "sequence",
    Toast = "toast",
}

/**
 * One element of a server-driven UI tree. The `type` is an open vocabulary: a client that
 * does not recognise a type renders a placeholder rather than failing.
 */
export interface UINode {
    children?: UINode[];
    id?:       string;
    /**
     * Component-specific data. Open by design.
     */
    props?: { [key: string]: any };
    slots?: { [key: string]: UINode[] };
    /**
     * Component discriminator, e.g. "PosterCard".
     */
    type: string;
}

export enum Surface {
    Drawer = "drawer",
    Modal = "modal",
    Sheet = "sheet",
}

export enum Tone {
    Accent = "accent",
    Danger = "danger",
    Info = "info",
    Neutral = "neutral",
    Success = "success",
    Warning = "warning",
}

/**
 * A component expressed as data: a name, default params, and a template of primitives.
 * Clients register definitions and expand them; this is how a module contributes a
 * component without shipping client code. A template node's props may hold binding objects
 * ({"$bind":"path"} / {"$match":{…}}) and control keys ($if / $ifNot / $each / $as); a node
 * of type "Outlet" renders the caller's children or a named slot.
 */
export interface ComponentDefinition {
    /**
     * The node type this definition provides.
     */
    name: string;
    /**
     * Default param values, overridden by the caller's props.
     */
    params?:  { [key: string]: any };
    template: UINode;
}
