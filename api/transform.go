package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Outputs
type Output struct {
	Type          string   `json:"type,omitempty"`
	ComponentType string   `json:"componentType,omitempty"`
	Groups        []string `json:"groups,omitempty"`
	Connections   []string `json:"connections,omitempty"`
	Icon          string   `json:"icon,omitempty"`
	Label         string   `json:"label,omitempty"`
}

// Inputs
type Datacenter struct {
	Description string `json:"description"`
	Default     bool   `json:"default"`
}

type Resource struct {
	Associations []Association `json:"associations"`
	Location     string        `json:"string"`
}

func (r *Resource) getGroups(associationName, typ, groupName, name string) []string {
	groups := make([]string, 0)
	for _, association := range r.Associations {
		if association.Id == fmt.Sprintf("%s.%s.%s", typ, groupName, name) {
			groups = append(groups, "dc1")

			associationName = strings.TrimPrefix(associationName, "dc1-")
			groups = append(groups, associationName)
		}
	}
	return groups
}

type Association struct {
	Id   string `json:"id"`
	Type string `json:"string"`
}

type Service struct {
	Meta    Meta   `json:"meta"`
	Port    int    `json:"port"`
	Address string `json:"address"`
}

func (s *Service) getComponentType() string {
	return s.Meta.getComponentType()
}

func (s *Service) getConnections(groups []string) []string {
	if len(groups) > 1 {
		client := groups[1]
		number := getNumber(client)

		return []string{"consul" + number}
	}

	return nil
}

func (s *Service) getGroups(input *Input, typ, groupName, name string) []string {
	for _, resources := range input.Resource {
		for associationName, resource := range resources {
			groups := resource.getGroups(associationName, typ, groupName, name)
			if len(groups) > 0 {
				return groups
			}
		}
	}
	return nil
}

func (s *Service) getLabel(name string) string {
	return splitLetterFromNumber(name)
}

func (s *Service) getIcon(name string) string {
	return "user"
}

type Meta struct {
	Version  string            `json:"version"`
	Software string            `json:"software"`
	Extra    map[string]string `json:"extra"`
}

func (m *Meta) getComponentType() string {
	switch m.Software {
	case "nginx":
		return "networking.nginx"
	case "postgres":
		return "database.postgres"
	default:
		return "generic.server"
	}
}

type Input struct {
	Datacenter map[string]Datacenter          `json:"datacenter"`
	Resource   map[string]map[string]Resource `json:"resource"`
	Service    map[string]map[string]Service  `json:"service"`
}

func (w *Input) fill(o map[string]Output) error {
	for groupName, srv := range w.Service {
		for name, data := range srv {
			groups := data.getGroups(w, "service", groupName, name)
			o[name] = Output{
				Type:          "component",
				ComponentType: data.getComponentType(),
				Groups:        groups,
				Connections:   data.getConnections(groups),
				//Icon:          "user",
				Label: data.getLabel(name),
			}
		}
	}

	for name, resource := range w.Resource {
		if name == "consul-client" {
			i := 1
			for client, _ := range resource {
				client = strings.TrimPrefix(client, "dc1-")
				o[client] = Output{
					Type:  "group",
					Label: strings.Title(splitLetterFromNumber(client)),
				}

				j := strconv.Itoa(i)

				o["consul"+j] = Output{
					Type:          "component",
					ComponentType: "devops.hashicorp-consul",
					Groups:        []string{"dc1", "client" + j},
					Label:         "consul",
				}

				i++
			}
		}

		if name == "consul-server" {
			i := 1
			for admin, _ := range resource {
				j := strconv.Itoa(i)
				admin = strings.TrimPrefix(admin, "dc1-")

				o["consulserver"+j] = Output{
					Type:          "component",
					ComponentType: "devops.hashicorp-consul",
					Groups:        []string{"dc1", admin},
					Connections: []string{
						"consul1",
						"consul2",
						"consul3",
						"consul4",
						"consul5",
					},
					Label: "consul",
				}

				o[admin] = Output{
					Type:  "group",
					Label: strings.Title(splitLetterFromNumber(admin)),
				}

				i++
			}
		}
	}

	for name, datacenter := range w.Datacenter {
		o[name] = Output{
			Type:  "group",
			Label: datacenter.Description,
		}
	}

	return nil
}

func splitLetterFromNumber(input string) string {
	var l, n []rune
	for _, r := range input {
		switch {
		case r >= 'A' && r <= 'Z':
			l = append(l, r)
		case r >= 'a' && r <= 'z':
			l = append(l, r)
		case r >= '0' && r <= '9':
			n = append(n, r)
		}
	}
	return fmt.Sprintf("%s %s", string(l), string(n))
}

func getNumber(input string) string {
	var n []rune
	for _, r := range input {
		switch {
		case r >= '0' && r <= '9':
			n = append(n, r)
		}
	}
	return string(n)
}

func transform(input []byte) ([]byte, error) {
	in := new(Input)
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	out := make(map[string]Output)
	if err := in.fill(out); err != nil {
		return nil, err
	}

	bs, err := json.Marshal(out)
	if err != nil {
		return nil, err
	}

	return bs, nil
}
