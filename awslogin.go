package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"gopkg.in/ini.v1"
)

const version = "1.0.1"

// Profile represents the AWS profile data
type Profile struct {
	Name    string
	Region  string
	Account string
	Role    string
}

func main() {
	// 1. Configurar Defaults para o Help
	home, _ := os.UserHomeDir()
	defaultConfig := filepath.Join(home, ".aws", "config")

	// 2. Definir Flags
	var (
		flagConnect string
		flagConfig  string
		flagList    bool
		flagVersion bool
	)

	flag.StringVar(&flagConnect, "connect", "", "Direct connection by profile name")
	flag.StringVar(&flagConfig, "config", defaultConfig, "Path to config file")
	flag.BoolVar(&flagList, "list", false, "List all configured profiles")
	flag.BoolVar(&flagVersion, "version", false, "Show current version")

	// 3. Custom Usage (Estilo pssql)
	flag.Usage = func() {
		fmt.Printf("aws-sso - AWS SSO Profile Manager\n\nUsage: aws-sso [flags]\n\nFlags:\n")
		fmt.Printf("  --connect string  Direct connection by profile name\n")
		fmt.Printf("  --config string   Path to config file (default: %s)\n", defaultConfig)
		fmt.Printf("  --list            List all configured profiles\n")
		fmt.Printf("  --help            Show this help\n")
		fmt.Printf("  --version         Show current version\n")
	}

	flag.Parse()

	// --- Logica das Flags ---

	if flagVersion {
		fmt.Printf("aws-sso version %s\n", version)
		os.Exit(0)
	}

	// Carregar perfis (usando o caminho default ou o da flag)
	profiles, err := loadProfiles(flagConfig)
	if err != nil {
		fmt.Printf("Error loading profiles: %v\n", err)
		os.Exit(1)
	}

	if flagList {
		fmt.Printf("%-30s | %-15s | %-15s | %s\n", "Name", "Region", "Account", "Role")
		fmt.Println(strings.Repeat("-", 90))
		for _, p := range profiles {
			fmt.Printf("%-30s | %-15s | %-15s | %s\n", p.Name, p.Region, p.Account, p.Role)
		}
		os.Exit(0)
	}

	if flagConnect != "" {
		// Procura o perfil especificado
		found := false
		for _, p := range profiles {
			if p.Name == flagConnect {
				if err := loginSSO(p.Name); err != nil {
					fmt.Printf("Login failed: %v\n", err)
					os.Exit(1)
				}
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Profile '%s' not found.\n", flagConnect)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// --- Modo Interativo (Só executa se nenhuma flag de ação for passada) ---

	// 1. Clear screen to start fresh (like pssql)
	print("\033[H\033[2J")

	selected, err := selectProfile(profiles)
	if err != nil {
		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}
		fmt.Printf("Selection failed: %v\n", err)
		os.Exit(1)
	}

	if err := loginSSO(selected); err != nil {
		fmt.Printf("Login failed: %v\n", err)
		os.Exit(1)
	}
}

// Alterado para receber o path como argumento
func loadProfiles(configPath string) ([]Profile, error) {
	// Se o path vier vazio (segurança), força o default
	if configPath == "" {
		home, _ := os.UserHomeDir()
		configPath = filepath.Join(home, ".aws", "config")
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not load %s: %w", configPath, err)
	}

	var profiles []Profile
	for _, section := range cfg.Sections() {
		name := section.Name()

		// AWS SSO configs usually have these keys
		region := section.Key("sso_region").String()
		account := section.Key("sso_account_id").String()
		role := section.Key("sso_role_name").String()

		// Fallback if keys are named differently
		if role == "" {
			role = section.Key("sso_role").String()
		}

		if name == "DEFAULT" || len(section.Keys()) == 0 {
			continue
		}

		var pName string
		if strings.HasPrefix(name, "profile ") {
			pName = strings.TrimPrefix(name, "profile ")
		} else if name == "default" {
			pName = "default"
		} else {
			continue
		}

		profiles = append(profiles, Profile{
			Name:    pName,
			Region:  region,
			Account: account,
			Role:    role,
		})
	}

	// Sort alphabetically by profile name
	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	if len(profiles) == 0 {
		return nil, fmt.Errorf("no profiles found in %s", configPath)
	}

	return profiles, nil
}

func selectProfile(profiles []Profile) (string, error) {
	// Helper function for ANSI colors to avoid dependency on promptui constants
	color := func(code string, s interface{}) string {
		return fmt.Sprintf("\033[%sm%v\033[0m", code, s)
	}

	// Define the FuncMap manually to ensure complete control over formatting and colors
	funcMap := template.FuncMap{
		"cyan":  func(s interface{}) string { return color("36", s) },
		"white": func(s interface{}) string { return color("37", s) },
		"green": func(s interface{}) string { return color("32", s) },
		"faint": func(s interface{}) string { return color("2", s) }, // Dim/Faint

		// The Column formatter (truncates if too long, pads if too short)
		"Col": func(l int, v interface{}) string {
			s := fmt.Sprintf("%v", v)
			if len(s) > l {
				return s[:l-1] + "…"
			}
			// %-*s pads to the right (left-aligned text)
			return fmt.Sprintf("%-*s", l, s)
		},
	}

	templates := &promptui.SelectTemplates{
		FuncMap: funcMap,
		Label:   "{{ . }}",
		// Note: pipe syntax implies the first arg is passed as last to the function
		// So {{ .Name | Col 30 }} calls Col(30, .Name)
		Active:   ` -> {{ .Name | Col 30 | cyan }} | {{ .Region | Col 15 | cyan }} | {{ .Account | Col 15 | cyan }} | {{ .Role | Col 30 | cyan }}`,
		Inactive: `    {{ .Name | Col 30 | white }} | {{ .Region | Col 15 | white }} | {{ .Account | Col 15 | white }} | {{ .Role | Col 30 | white }}`,
		Selected: `{{ "✔" | green }} Login to: {{ .Name | cyan }}`,
	}

	prompt := promptui.Select{
		Label:     "Select Profile",
		Items:     profiles,
		Templates: templates,
		Size:      20,
		Searcher: func(input string, index int) bool {
			p := profiles[index]
			content := fmt.Sprintf("%s %s %s", p.Name, p.Account, p.Role)
			return strings.Contains(strings.ToLower(content), strings.ToLower(input))
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return profiles[i].Name, nil
}

func loginSSO(profile string) error {
	// Simple output message
	fmt.Printf("Executing: aws sso login --profile %s\n", profile)

	cmd := exec.Command("aws", "sso", "login", "--profile", profile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
