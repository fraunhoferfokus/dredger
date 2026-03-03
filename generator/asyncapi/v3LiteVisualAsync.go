package generator

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// LiteSpec is a simplified AsyncAPI representation used for documentation.
// LiteSpec aggregates the most relevant details of an AsyncAPI specification
// so it can be rendered as a static HTML page without any heavy runtime
// dependencies. Only a subset of the full spec is parsed.
type LiteSpec struct {
	Title        string
	Version      string
	Description  string
	ContactName  string
	ContactEmail string
	ContactURL   string
	LicenseName  string
	Servers      []LiteServer
	Channels     []LiteChannel
	Operations   []LiteOperation
	Messages     []LiteMessage
}

// LiteChannel summarizes a single channel with its messages.
type LiteChannel struct {
	Name     string
	Address  string
	Summary  string
	Messages []string
}

// LiteServer describes one server entry.
type LiteServer struct {
	Name        string
	Host        string
	Protocol    string
	Description string
}

// LiteOperation summarizes an operation referencing a channel.
type LiteOperation struct {
	Name     string
	Action   string
	Channel  string
	Summary  string
	Messages []string
}

// LiteMessage describes a message definition in components.
type LiteMessage struct {
	Name        string
	Title       string
	Summary     string
	ContentType string
	Properties  []MessageProperty
}

// MessageProperty holds a single property name/type from a message payload.
type MessageProperty struct {
	Name string
	Type string
}

// ParseLite loads an AsyncAPI v3 specification from path and extracts only
// the basic information needed for the static documentation view.
func ParseLite(path string) (*LiteSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	spec := &LiteSpec{}

	if info, ok := raw["info"].(map[string]any); ok {
		if v, ok := info["title"].(string); ok {
			spec.Title = v
		}
		if v, ok := info["version"].(string); ok {
			spec.Version = v
		}
		if v, ok := info["description"].(string); ok {
			spec.Description = v
		}
		if c, ok := info["contact"].(map[string]any); ok {
			if v, ok := c["name"].(string); ok {
				spec.ContactName = v
			}
			if v, ok := c["email"].(string); ok {
				spec.ContactEmail = v
			}
			if v, ok := c["url"].(string); ok {
				spec.ContactURL = v
			}
		}
		if l, ok := info["license"].(map[string]any); ok {
			if v, ok := l["name"].(string); ok {
				spec.LicenseName = v
			}
		}
	}

	// collect message information from components
	msgSummary := map[string]string{}
	if comp, ok := raw["components"].(map[string]any); ok {
		if msgs, ok := comp["messages"].(map[string]any); ok {
			for name, m := range msgs {
				if mm, ok := m.(map[string]any); ok {
					lm := LiteMessage{Name: name}
					if v, ok := mm["title"].(string); ok {
						lm.Title = v
					}
					if v, ok := mm["summary"].(string); ok {
						lm.Summary = v
						msgSummary[name] = v
					}
					if v, ok := mm["contentType"].(string); ok {
						lm.ContentType = v
					}
					// collect top-level property names
					if payload, ok := mm["payload"].(map[string]any); ok {
						if props, ok := payload["properties"].(map[string]any); ok {
							for pn, pv := range props {
								if pm, ok := pv.(map[string]any); ok {
									prop := MessageProperty{Name: pn}
									if t, ok := pm["type"].(string); ok {
										prop.Type = t
									}
									lm.Properties = append(lm.Properties, prop)
								}
							}
						}
					}
					spec.Messages = append(spec.Messages, lm)
				}
			}
		}
	}

	if servers, ok := raw["servers"].(map[string]any); ok {
		for name, sv := range servers {
			if sm, ok := sv.(map[string]any); ok {
				ls := LiteServer{Name: name}
				if v, ok := sm["host"].(string); ok {
					ls.Host = v
				}
				if v, ok := sm["protocol"].(string); ok {
					ls.Protocol = v
				}
				if v, ok := sm["description"].(string); ok {
					ls.Description = v
				}
				spec.Servers = append(spec.Servers, ls)
			}
		}
	}

	chans, ok := raw["channels"].(map[string]any)
	if !ok {
		chans = map[string]any{}
	}

	for name, ch := range chans {
		chm, ok := ch.(map[string]any)
		if !ok {
			continue
		}
		lc := LiteChannel{Name: name}
		if addr, ok := chm["address"].(string); ok {
			lc.Address = addr
		}
		if s, ok := chm["summary"].(string); ok {
			lc.Summary = s
		} else if s, ok := chm["description"].(string); ok {
			lc.Summary = s
		}

		if msgs, ok := chm["messages"].(map[string]any); ok {
			for msgName := range msgs {
				entry := msgName
				if sum, ok := msgSummary[msgName]; ok && sum != "" {
					entry = msgName + " - " + sum
				}
				lc.Messages = append(lc.Messages, entry)
			}
		}
		spec.Channels = append(spec.Channels, lc)
	}

	if ops, ok := raw["operations"].(map[string]any); ok {
		for name, op := range ops {
			if om, ok := op.(map[string]any); ok {
				lo := LiteOperation{Name: name}
				if v, ok := om["action"].(string); ok {
					lo.Action = v
				}
				if v, ok := om["summary"].(string); ok {
					lo.Summary = v
				} else if v, ok := om["description"].(string); ok {
					lo.Summary = v
				}
				if chref, ok := om["channel"].(map[string]any); ok {
					if ref, ok := chref["$ref"].(string); ok {
						lo.Channel = lastSegment(ref)
					} else if n, ok := chref["name"].(string); ok {
						lo.Channel = n
					}
				}
				if msgs, ok := om["messages"].([]any); ok {
					for _, mv := range msgs {
						if mm, ok := mv.(map[string]any); ok {
							if ref, ok := mm["$ref"].(string); ok {
								lo.Messages = append(lo.Messages, lastSegment(ref))
							} else if name, ok := mm["name"].(string); ok {
								lo.Messages = append(lo.Messages, name)
							}
						}
					}
				}
				spec.Operations = append(spec.Operations, lo)
			}
		}
	}

	return spec, nil
}

// lastSegment returns the part of a JSON Pointer or path after the final '/'.
func lastSegment(ref string) string {
	if i := strings.LastIndex(ref, "/"); i >= 0 && i < len(ref)-1 {
		return ref[i+1:]
	}
	return ref
}