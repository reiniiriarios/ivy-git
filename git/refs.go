package git

import (
	"strings"
)

type Ref struct {
	Hash      string
	Name      string
	Branch    string
	Remote    string
	AbbrName  string
	Annotated bool
}

// Get ref details from `git show-ref`.
func (g *Git) getRefs() (Refs, error) {
	var refs Refs

	show_refs, err := g.RunCwd("show-ref", "--dereference", "--head")
	if err != nil {
		println(err.Error())
		return refs, err
	}

	// For the purposes of displaying a coherent tree,
	// we're denoting the following:
	// - refs/heads/[...]                 = branches
	// - refs/tags/[...]                  = tags
	// - HEAD and refs/remotes/[...]/HEAD = heads
	// - refs/remotes/[...]/[...]         = remotes
	show_refs_lines := parseLines(show_refs)
	for _, r := range show_refs_lines {
		ref_details := strings.Split(r, " ")
		if len(ref_details) >= 2 {
			hash := ref_details[0]
			name := strings.Join(ref_details[1:], " ")
			if strings.HasPrefix(name, "refs/heads/") {
				refs.Branches = append(refs.Branches, parseRefHead(hash, name[11:]))
			} else if strings.HasPrefix(name, "refs/tags/") {
				refs.Tags = append(refs.Tags, parseRefTag(hash, name[10:]))
			} else if strings.HasPrefix(name, "refs/remotes/") {
				ref := parseRefRemote(hash, name[13:])
				if name[len(name)-4:] == "HEAD" {
					refs.Heads = append(refs.Heads, ref)
				} else {
					refs.Remotes = append(refs.Remotes, ref)
				}
			} else if name == "HEAD" {
				refs.HEAD = Ref{
					Hash: hash,
					Name: name,
				}
			} else if !strings.HasPrefix(name, "refs/stash") {
				// Ignore stash, but anything else log a warning.
				println("Error parsing ref:", name)
			}
		}
	}

	return refs, nil
}

func parseRefHead(hash string, name string) Ref {
	abbr := ""
	if len(name) > REF_MAX_NAME_LENGTH {
		abbr = name[:REF_MAX_NAME_LENGTH] + "..."
	}
	return Ref{
		Hash:     hash,
		Name:     name,
		Branch:   name,
		AbbrName: abbr,
	}
}

func parseRefTag(hash string, name string) Ref {
	annotated := strings.HasSuffix(name, "^{}")
	if annotated {
		name = name[:len(name)-3]
	}
	return Ref{
		Hash:      hash,
		Name:      name,
		Annotated: annotated,
	}
}

func parseRefRemote(hash string, name string) Ref {
	if name[len(name)-4:] == "HEAD" {
		remote := ""
		name_parts := strings.Split(name, "/")
		if len(name_parts) >= 2 {
			remote = name_parts[0]
		}
		return Ref{
			Hash:   hash,
			Name:   name,
			Remote: remote,
		}
	} else {
		branch := name
		remote := ""
		name_parts := strings.Split(name, "/")
		if len(name_parts) >= 2 {
			remote = name_parts[0]
			branch = name[len(remote)+1:]
		}
		abbr := ""
		if len(branch) > REF_MAX_NAME_LENGTH {
			abbr = branch[:REF_MAX_NAME_LENGTH] + "..."
		}
		return Ref{
			Hash:     hash,
			Name:     name,
			Branch:   branch,
			Remote:   remote,
			AbbrName: abbr,
		}
	}
}
