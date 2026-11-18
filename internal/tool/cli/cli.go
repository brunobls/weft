// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package cli implements the Weft command line tool. It is shared by the
// "weft" command and its "weaver" alias, which differ only in the name they
// report.
package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	itool "github.com/brunobls/weft/internal/tool"
	"github.com/brunobls/weft/internal/tool/callgraph"
	"github.com/brunobls/weft/internal/tool/generate"
	"github.com/brunobls/weft/internal/tool/multi"
	"github.com/brunobls/weft/internal/tool/single"
	"github.com/brunobls/weft/internal/tool/ssh"
	"github.com/brunobls/weft/runtime/tool"
)

// Usage returns the help text of the command invoked as name.
func Usage(name string) string {
	u := strings.ReplaceAll(usageTmpl, "@", name)
	if name == "weaver" {
		u = strings.Replace(u, `dispatch to a binary called "weaver-<deployer>".`,
			`dispatch to a binary called "weft-<deployer>" or, if there is none,
  "weaver-<deployer>".`, 1)
		u += "\n  weaver is an alias of weft.\n"
	}
	return u
}

const usageTmpl = `USAGE

  @ generate                 // code generator
  @ version                  // show version
  @ single    <command> ...  // for single process deployments
  @ multi     <command> ...  // for multiprocess deployments
  @ ssh       <command> ...  // for multimachine deployments
  @ gke       <command> ...  // for GKE deployments
  @ gke-local <command> ...  // for simulated GKE deployments
  @ kube      <command> ...  // for vanilla Kubernetes deployments

DESCRIPTION

  Use the "@" command to deploy and manage Weft applications.

  The "@ generate", "@ version", "@ single", "@ multi", and
  "@ ssh" subcommands are baked in, but all other subcommands of the form
  "@ <deployer>" dispatch to a binary called "@-<deployer>".
  "@ gke status", for example, dispatches to "@-gke status".
`

// Main runs the command line tool under the given name ("weft" or its
// "weaver" alias). It parses os.Args and never returns on error.
func Main(name string) {
	usage := Usage(name)
	// Parse flags.
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	// Handle the internal deployers.
	internals := map[string]map[string]*tool.Command{
		"single": single.Commands,
		"multi":  multi.Commands,
		"ssh":    ssh.Commands,
	}

	switch flag.Arg(0) {
	case "generate":
		generateFlags := flag.NewFlagSet("generate", flag.ExitOnError)
		tags := generateFlags.String("tags", "", "Optional tags for the generate command")
		generateFlags.Usage = func() {
			fmt.Fprintln(os.Stderr, generate.Usage(name))
		}
		generateFlags.Parse(flag.Args()[1:])
		buildTags := "ignoreWeaverGen"
		if *tags != "" { // tags flag was specified
			// TODO(rgrandl): we assume that the user specify the tags properly. I.e.,
			// a single tag, or a list of tags separated by comma. We may want to do
			// extra validation at some point.
			buildTags = buildTags + "," + *tags
		}
		if err := generate.Generate(".", generateFlags.Args(), generate.Options{BuildTags: buildTags}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return

	case "version":
		cmd := itool.VersionCmd(name)
		if err := cmd.Fn(context.Background(), flag.Args()[1:]); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return

	case "callgraph":
		usage := strings.ReplaceAll(`Generate component callgraphs.

Usage:
  @ callgraph <binary>

Flags:
  -h, --help           Print this help message.

Description:
  "@ callgraph <file>" outputs a component callgraph in mermaid format
  [1]. These graphs can be included in GitHub README files [2].

[1]: https://mermaid.js.org/
[2]: https://github.blog/2022-02-14-include-diagrams-markdown-files-mermaid/`, "@", name)
		flags := flag.NewFlagSet("callgraph", flag.ExitOnError)
		flags.Usage = func() { fmt.Fprintln(os.Stderr, usage) }
		flags.Parse(flag.Args()[1:])
		if flags.NArg() == 0 {
			fmt.Fprintln(os.Stderr, "ERROR: no binary provided.")
		}
		s, err := callgraph.Mermaid(flags.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(s)
		return

	case "single", "multi", "ssh":
		os.Args = os.Args[1:]
		tool.Run(name+" "+flag.Arg(0), internals[flag.Arg(0)])
		return

	case "help":
		n := len(flag.Args())
		command := flag.Arg(1)
		switch {
		case n == 1:
			// help
			fmt.Fprint(os.Stdout, usage)
		case n == 2 && command == "generate":
			// help generate
			fmt.Fprintln(os.Stdout, generate.Usage(name))
		case n == 2 && internals[command] != nil:
			// help <command>
			fmt.Fprintln(os.Stdout, tool.MainHelp(name+" "+command, internals[command]))
		case n == 2:
			// help <external>
			code, err := run(name, command, []string{"--help"})
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				os.Exit(code)
			}
		case n > 2:
			fmt.Fprintf(os.Stderr, "help: too many arguments. Try '%s %s %s --help'\n", name, command, strings.Join(flag.Args()[2:], " "))
		}
		return
	}

	// Handle all other "weaver <deployer>" subcommands.
	code, err := run(name, flag.Args()[0], flag.Args()[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(code)
	}
}

// pluginCandidates returns the binaries that may implement the deployer, in
// lookup order: "weft-<deployer>", then "weaver-<deployer>" for compatibility.
func pluginCandidates(deployer string) []string {
	return []string{"weft-" + deployer, "weaver-" + deployer}
}

// findPlugin returns the first candidate binary found by lookPath.
func findPlugin(deployer string, lookPath func(string) (string, error)) (string, error) {
	cands := pluginCandidates(deployer)
	for _, c := range cands {
		if _, err := lookPath(c); err == nil {
			return c, nil
		}
	}
	return "", fmt.Errorf("none of %q found", cands)
}

// run runs "weft-<deployer> [arg]..." (or "weaver-<deployer>") in a subprocess
// and returns the subprocess' exit code and any error.
func run(name, deployer string, args []string) (int, error) {
	binary, err := findPlugin(deployer, exec.LookPath)
	if err != nil {
		cands := pluginCandidates(deployer)
		msg := fmt.Sprintf(`"%s %s" is not a %s command. See "%s --help". If you're trying to invoke a custom deployer, the %q binary (or %q) was not found. You may need to install the %q binary or add it to your PATH.`, name, deployer, name, name, cands[0], cands[1], cands[0])
		return 1, fmt.Errorf("%s", wrap(msg, 80))
	}
	cmd := exec.Command(binary, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		return 0, nil
	}
	exitError := &exec.ExitError{}
	if errors.As(err, &exitError) {
		return exitError.ExitCode(), err
	}
	return 1, err
}

// wrap trims whitespace in the provided string and wraps it to n characters.
func wrap(s string, n int) string {
	var b strings.Builder
	k := 0
	for i, word := range strings.Fields(s) {
		if i == 0 {
			k = len(word)
			fmt.Fprintf(&b, "%s", word)
		} else if k+len(word)+1 > n {
			k = len(word)
			fmt.Fprintf(&b, "\n%s", word)
		} else {
			k += len(word) + 1
			fmt.Fprintf(&b, " %s", word)
		}
	}
	return b.String()
}
