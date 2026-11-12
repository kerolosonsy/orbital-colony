package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
)

type Pod struct {
	Name    string
	Running bool
}

type Runtime interface {
	List(ctx context.Context) ([]Pod, error)
	Start(ctx context.Context, name string) error
	Stop(ctx context.Context, name string) error
}

const podUsage = "usage: podctl list [--running] | start NAME | stop [--yes] NAME\n"

func RunPodCtl(ctx context.Context, rt Runtime, args []string, in io.Reader, out, errOut io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(errOut, podUsage)
		return 2
	}
	cmd, rest := args[0], args[1:]
	switch {
	case cmd == "list" && (len(rest) == 0 || len(rest) == 1 && rest[0] == "--running"):
		pods, err := rt.List(ctx)
		if err != nil {
			fmt.Fprintf(errOut, "runtime: %v\n", err)
			return 1
		}
		pods = slices.Clone(pods)
		slices.SortFunc(pods, func(a, b Pod) int { return strings.Compare(a.Name, b.Name) })
		for _, p := range pods {
			state := "stopped"
			if p.Running {
				state = "running"
			} else if len(rest) == 1 {
				continue
			}
			fmt.Fprintf(out, "%s\t%s\n", p.Name, state)
		}
		return 0
	case cmd == "start" && len(rest) == 1:
		return podChange(ctx, rt, rest[0], true, false, in, out, errOut)
	case cmd == "stop" && len(rest) == 1 && !strings.HasPrefix(rest[0], "-"):
		return podChange(ctx, rt, rest[0], false, false, in, out, errOut)
	case cmd == "stop" && len(rest) == 2 && rest[0] == "--yes":
		return podChange(ctx, rt, rest[1], false, true, in, out, errOut)
	}
	fmt.Fprint(errOut, podUsage)
	return 2
}

func podChange(ctx context.Context, rt Runtime, name string, start, yes bool, in io.Reader, out, errOut io.Writer) int {
	pods, err := rt.List(ctx)
	if err != nil {
		fmt.Fprintf(errOut, "runtime: %v\n", err)
		return 1
	}
	i := slices.IndexFunc(pods, func(p Pod) bool { return p.Name == name })
	if i < 0 {
		fmt.Fprintf(errOut, "unknown pod %s\n", name)
		return 1
	}
	if pods[i].Running == start {
		if start {
			fmt.Fprintf(out, "already running %s\n", name)
		} else {
			fmt.Fprintf(out, "already stopped %s\n", name)
		}
		return 0
	}
	if start {
		err = rt.Start(ctx, name)
	} else {
		if !yes {
			fmt.Fprintf(out, "stop %s? [y/N] ", name)
			sc := bufio.NewScanner(in)
			sc.Scan()
			answer := strings.ToLower(strings.TrimSpace(sc.Text()))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(out, "aborted")
				return 1
			}
		}
		err = rt.Stop(ctx, name)
	}
	if err != nil {
		fmt.Fprintf(errOut, "runtime: %v\n", err)
		return 1
	}
	if start {
		fmt.Fprintf(out, "started %s\n", name)
	} else {
		fmt.Fprintf(out, "stopped %s\n", name)
	}
	return 0
}

func main() {}
