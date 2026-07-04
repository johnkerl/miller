// Entrypoint for `mlr skill`: puts the Miller Agent Skill (SKILL.md) where a
// coding agent can find it on disk, for tools that read Agent Skills
// directly rather than over MCP.
//
// The content is identical to what `mlr mcp` serves as its "miller-playbook"
// prompt/resource (pkg/terminals/mcp/SKILL.md, exported as mcp.PlaybookText)
// -- this is a second delivery path for the same text, not a second source
// of truth.

package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/johnkerl/miller/v6/pkg/terminals/mcp"
)

const defaultInstallDir = ".claude/skills/miller"

func skillUsage(o *os.File) {
	fmt.Fprintf(o, "Usage: mlr skill {print|install} [options]\n")
	fmt.Fprintf(o, "Puts the Miller Agent Skill (SKILL.md) where a coding agent can find it.\n")
	fmt.Fprintf(o, "This is the same playbook mlr mcp serves as its \"miller-playbook\"\n")
	fmt.Fprintf(o, "prompt/resource, packaged for agents that read Agent Skills from disk.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Subcommands:\n")
	fmt.Fprintf(o, "  print          Write the skill content to stdout.\n")
	fmt.Fprintf(o, "  install [DIR]  Write DIR/SKILL.md, creating DIR if needed.\n")
	fmt.Fprintf(o, "                 Default DIR is %s\n", defaultInstallDir)
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, " -h or --help   Show this message.\n")
}

// SkillMain is the entrypoint called by the terminals dispatcher for `mlr skill`.
func SkillMain(args []string) int {
	args = args[1:] // strip "skill"

	if len(args) == 0 {
		skillUsage(os.Stderr)
		return 1
	}

	switch args[0] {
	case "-h", "--help":
		skillUsage(os.Stdout)
		return 0
	case "print":
		return printMain(args[1:])
	case "install":
		return installMain(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "mlr skill: subcommand \"%s\" not recognized.\n", args[0])
		return 1
	}
}

func printMain(args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "mlr skill print: takes no arguments.\n")
		return 1
	}
	fmt.Print(mcp.PlaybookText)
	return 0
}

func installMain(args []string) int {
	dir := defaultInstallDir
	switch len(args) {
	case 0:
	case 1:
		dir = args[0]
	default:
		fmt.Fprintf(os.Stderr, "mlr skill install: takes at most one argument (the target directory).\n")
		return 1
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "mlr skill install: could not create %s: %v\n", dir, err)
		return 1
	}

	path := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(path, []byte(mcp.PlaybookText), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "mlr skill install: could not write %s: %v\n", path, err)
		return 1
	}

	fmt.Printf("Wrote %s\n", path)
	return 0
}
