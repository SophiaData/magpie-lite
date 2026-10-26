package provider

import (
	"cmp"
	"slices"
)

// A preset is a vendor magpie already knows: adding one only asks for the key.

// Kind groups presets in the picker.
type Kind string

const (
	KindVendor Kind = "vendor" // the model's own maker
	KindRelay  Kind = "relay"  // an aggregator / API relay reselling many vendors
	KindLocal  Kind = "local"  // something running on this machine
	// KindRemote is another computer's magpie: not a relay, so the add
	// lists give it a heading of its own, after this machine's
	KindRemote Kind = "remote"
)

// PresetDef describes one preset.
type PresetDef struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Icon      string   `json:"icon"`
	Kind      Kind     `json:"kind"`
	Chat      string   `json:"chat,omitempty"`
	Responses string   `json:"responses,omitempty"`
	Anthropic string   `json:"anthropic,omitempty"`
	Decide    string   `json:"decide,omitempty"` // a decision API: the provider only routes (see decide.go)
	Catalog   string   `json:"catalog,omitempty"`
	Website   string   `json:"website,omitempty"`
	SiteHost  string   `json:"siteHost,omitempty"` // a partner's: the host its usemagpie.ai/go Website sends the browser to, which the editor's link reads
	KeysURL   string   `json:"keysUrl,omitempty"`
	NoKey     bool     `json:"noKey,omitempty"`     // local servers: a key is optional
	KeyHint   string   `json:"keyHint,omitempty"`   // the key field's placeholder, when it says more than NoKey's
	Sponsored bool     `json:"sponsored,omitempty"` // shown first, with a tag
	Note      string   `json:"note,omitempty"`      // one line under the name
	Short     string   `json:"short,omitempty"`     // the add sheet's name for it, when Name is long
	Regions   []Region `json:"regions,omitempty"`   // base-URL choices (a relay's regional endpoints, a vendor's plans)
	// RegionLabel names what the Regions choose between, "Region" if unset.
	RegionLabel string `json:"regionLabel,omitempty"`
	// KeyPlans: each key may be on a Region of its own (KeyAccount.Plan),
	// a vendor whose plans and pay as you go are sold on keys of their own
	// that only their own endpoints take (Volcengine's Coding and Agent
	// Plans beside its pay as you go, yckonnnn on X). Not a relay's
	// regions, where one key works at all of them.
	KeyPlans bool `json:"keyPlans,omitempty"`
	// HeaderHints name optional request headers the vendor documents, which
	// the editor offers to add; their values are the user's to fill in.
	HeaderHints []string `json:"headerHints,omitempty"`
	// Only, for a plan behind a key that also buys much else, is what its
	// models' ids start with: the rest of the vendor's list is left out.
	// Models are the plan's, for when the list has none of them.
	Only   string   `json:"only,omitempty"`
	Models []string `json:"models,omitempty"`
	// NoList: the vendor has no list of models to ask for (Bedrock's
	// runtime serves no /models), so Models are its list
	NoList bool `json:"noList,omitempty"`
	// Endpoint, for a vendor reached at the user's own resource (Azure
	// OpenAI), is an example of its address: the preset gives no URL, and
	// the editor asks for the one the user's resource is at, with
	// EndpointHint under it.
	Endpoint     string `json:"endpoint,omitempty"`
	EndpointHint string `json:"endpointHint,omitempty"`
	// EndpointNeeded is what the editor says when no endpoint was given.
	EndpointNeeded string `json:"endpointNeeded,omitempty"`
	// Hosts: a vendor serving other makers' models as well as its own
	// (Groq, Ollama Cloud, Azure), whose list is no maker's word on theirs
	Hosts bool `json:"-"`
}

// Region is one base-URL option of a preset that offers several. The first
// is the default; picking another in the editor swaps the endpoints.
type Region struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Chat      string `json:"chat,omitempty"`
	Responses string `json:"responses,omitempty"`
	Anthropic string `json:"anthropic,omitempty"`
	// Lists is set on a region that serves a model list although the
	// preset as a whole has none to ask for (NoList): a provider at its
	// endpoints is asked for it. KeysURL, when the region's keys are made
	// on another page than the preset's, is where its Get-a-key link goes.
	Lists   bool   `json:"lists,omitempty"`
	KeysURL string `json:"keysUrl,omitempty"`
	// Website, when the region's docs are another page than the preset's,
	// is where the editor's link goes.
	Website  string `json:"website,omitempty"`
	SiteHost string `json:"siteHost,omitempty"` // as PresetDef's
	// Catalog, for a region whose models another models.dev entry names
	// than the preset's (Tencent Cloud's pay as you go, beside its plan's
	// none), is that entry's id.
	Catalog string `json:"catalog,omitempty"`
	// Models are a plan's text models, for a region that is one beside
	// another plan or pay as you go (PresetDef.Models, for the region alone).
	// Embeddings, Drawers and Videos are its non-chat models, kept apart from
	// agents' model pickers.
	Models     []string `json:"models,omitempty"`
	Embeddings []string `json:"embeddings,omitempty"`
	Drawers    []string `json:"drawers,omitempty"`
	Videos     []string `json:"videos,omitempty"`
	// Decide is the region's decision API, for a preset that routes groups
	Decide string `json:"decide,omitempty"`
}

// WorkspaceID stands in a decision API's address for the user's Bailian
// workspace id (#647), which the editor and `workspace=` fill in.
const WorkspaceID = "{WorkspaceId}"

// bailianWorkspace is Bailian's System One root in a workspace of region.
func bailianWorkspace(region string) string {
	return "https://" + WorkspaceID + "." + region + ".maas.aliyuncs.com/compatible-mode/v1"
}

// alibabaPlanModels are the text models of Alibaba's Token Plan, the same
// on the Qwen AI platform's and on Bailian's, as their overviews list them.
var alibabaPlanModels = []string{"auto", "qwen3.8-max", "qwen3.8-flash", "qwen3.7-max", "qwen3.7-plus", "qwen3.6-flash",
	"deepseek-v4.1-flash", "deepseek-v4-pro", "deepseek-v4-pro-0813", "deepseek-v4-flash-0731", "glm-5.3", "glm-5.2"}

// presets are ordered as they appear in the picker.
var presets = []PresetDef{
	// Vendors
	{ID: "deepseek", Name: "DeepSeek", Icon: "deepseek-color", Kind: KindVendor, Catalog: "deepseek",
		Chat: "https://api.deepseek.com/v1", Responses: "https://api.deepseek.com/v1", Anthropic: "https://api.deepseek.com/anthropic",
		Website: "https://platform.deepseek.com", KeysURL: "https://platform.deepseek.com/api_keys"},
	{ID: "anthropic", Name: "Anthropic", Icon: "claude-color", Kind: KindVendor, Catalog: "anthropic",
		Anthropic: "https://api.anthropic.com",
		Website:   "https://console.anthropic.com", KeysURL: "https://console.anthropic.com/settings/keys",
		HeaderHints: []string{"anthropic-workspace-id"}},
	{ID: "openai", Name: "OpenAI", Icon: "openai", Kind: KindVendor, Catalog: "openai",
		Chat: "https://api.openai.com/v1", Responses: "https://api.openai.com/v1",
		Website: "https://platform.openai.com", KeysURL: "https://platform.openai.com/api-keys"},
	{ID: "qwen-cn", Name: "Qwen (DashScope)", Icon: "qwen-color", Kind: KindVendor, Catalog: "alibaba",
		Chat: "https://dashscope.aliyuncs.com/compatible-mode/v1", Anthropic: "https://dashscope.aliyuncs.com/apps/anthropic",
		Note:    "DashScope · China",
		Website: "https://bailian.console.aliyun.com", KeysURL: "https://bailian.console.aliyun.com/?tab=model#/api-key"},
	{ID: "zhipu", Name: "Zhipu GLM", Icon: "zhipu-color", Kind: KindVendor, Catalog: "zhipuai",
		Chat: "https://open.bigmodel.cn/api/paas/v4", Anthropic: "https://open.bigmodel.cn/api/anthropic",
		Website: "https://open.bigmodel.cn", KeysURL: "https://open.bigmodel.cn/usercenter/proj-mgmt/apikeys"},
	{ID: "moonshot-cn", Name: "Kimi (Moonshot)", Icon: "kimi", Kind: KindVendor, Catalog: "moonshotai",
		Chat: "https://api.moonshot.cn/v1", Anthropic: "https://api.moonshot.cn/anthropic",
		Website: "https://platform.moonshot.cn", KeysURL: "https://platform.moonshot.cn/console/api-keys"},
	{ID: "google", Name: "Google Gemini", Icon: "gemini-color", Kind: KindVendor, Catalog: "google",
		Chat:    "https://generativelanguage.googleapis.com/v1beta/openai",
		Note:    "Gemini Developer API",
		Website: "https://aistudio.google.com", KeysURL: "https://aistudio.google.com/apikey"},

	// Relays
	{ID: "openrouter", Name: "OpenRouter", Icon: "openrouter", Kind: KindRelay, Catalog: "openrouter",
		Chat: "https://openrouter.ai/api/v1", Anthropic: "https://openrouter.ai/api",
		Website: "https://openrouter.ai", KeysURL: "https://openrouter.ai/keys",
		HeaderHints: []string{"HTTP-Referer", "X-OpenRouter-Title"}},
	{ID: "opencode-zen", Name: "OpenCode Zen", Icon: "opencode", Kind: KindRelay, Catalog: "opencode", NoKey: true,
		KeyHint: "optional: free models need no key",
		Chat:    "https://opencode.ai/zen/v1", Responses: "https://opencode.ai/zen/v1", Anthropic: "https://opencode.ai/zen",
		Website: "https://opencode.ai/docs/zen", KeysURL: "https://opencode.ai/auth"},

	// Local
	{ID: "ollama", Name: "Ollama", Icon: "ollama", Kind: KindLocal, NoKey: true,
		Chat: "http://localhost:11434/v1", Anthropic: "http://localhost:11434",
		Note:   "your local models", Website: "https://ollama.com"},
	{ID: "lmstudio", Name: "LM Studio", Icon: "lmstudio", Kind: KindLocal, NoKey: true,
		Chat: "http://localhost:1234/v1",
		Note: "local server, port 1234 by default", Website: "https://lmstudio.ai"},
}

func Presets() []PresetDef {
	out := make([]PresetDef, 0, len(presets))
	for _, p := range presets {
		if p.Sponsored {
			out = append(out, p)
		}
	}
	for _, p := range presets {
		if !p.Sponsored {
			out = append(out, p)
		}
	}
	return out
}

// presetAliases are ids presets carried before, each the preset since and
// the region it is (none: the preset's default): qianfan's first day's
// (qianfan-token-plan, v0.1.394), and Tencent Cloud's three presets, its
// Token Plan and TokenHub's China and global pay as you go, made one.
var presetAliases = map[string]struct{ preset, region string }{
	"qianfan-token-plan":  {"baidu-qianfan", ""},
	"tencent-token-plan":  {"tencent-cloud", "plan"},
	"tencent-tokenhub-cn": {"tencent-cloud", "cn"},
	"tencent-tokenhub":    {"tencent-cloud", "intl"},
}

// Preset finds a preset by id, or by an id it carried before
// (presetAliases), or a partner listed now or before (partners.go).
func Preset(id string) *PresetDef {
	if pr := builtinPreset(id); pr != nil {
		return pr
	}
	return partnerPreset(id)
}

// builtinPreset is the preset of that id built into magpie, partners'
// aside.
func builtinPreset(id string) *PresetDef {
	if a, ok := presetAliases[id]; ok {
		id = a.preset
	}
	for i := range presets {
		if presets[i].ID == id {
			return &presets[i]
		}
	}
	return nil
}

// FromPreset builds a provider from a preset; the caller adds the key. An
// id the preset carried before builds it at the region that id was.
func FromPreset(id string) (Provider, error) {
	pr := Preset(id)
	if pr == nil {
		return Provider{}, errorf("no preset %q — magpie presets lists them", id)
	}
	p := Provider{
		ID: pr.ID, Name: pr.Name, Icon: pr.Icon, Preset: pr.ID,
		Chat: pr.Chat, Responses: pr.Responses, Anthropic: pr.Anthropic, Decide: pr.Decide,
		Catalog: pr.Catalog, Website: pr.Website, KeysURL: pr.KeysURL,
	}
	if a, ok := presetAliases[id]; ok && a.region != "" {
		for _, r := range pr.Regions {
			if r.ID == a.region {
				p.atRegionOf(pr, r)
			}
		}
	}
	return p, nil
}

// atRegionOf puts the provider at a region of its preset pr: its
// endpoints, and its key page, docs and catalog where it has its own.
func (p *Provider) atRegionOf(pr *PresetDef, r Region) {
	p.Chat, p.Responses, p.Anthropic, p.Decide = r.Chat, r.Responses, r.Anthropic, cmp.Or(r.Decide, pr.Decide)
	p.KeysURL = cmp.Or(r.KeysURL, p.KeysURL)
	p.Website = cmp.Or(r.Website, p.Website)
	p.Catalog = cmp.Or(r.Catalog, pr.Catalog)
}

// IconForCatalog names the logo of the vendor behind a models.dev
// provider id, or "" when no preset covers it.
func IconForCatalog(catalogID string) string {
	for _, p := range presets {
		if p.Catalog == catalogID || slices.ContainsFunc(p.Regions, func(r Region) bool { return r.Catalog == catalogID }) {
			return p.Icon
		}
	}
	return ""
}
