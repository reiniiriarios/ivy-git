package git

import (
	"strings"
)

type Ref struct {
	Hash          string
	Name          string
	Branch        string
	Remote        string
	SyncedRemotes []string
	SyncedLocally bool
	Head          bool
	Upstream      string
	AbbrName      string
	Annotated     bool
}

type Refs struct {
	HEAD           Ref
	Branches       []Ref
	Tags           []Ref
	RemoteBranches []Ref
	Heads          []Ref
}

// Get the hash of the last commit on main.
func (g *Git) lastCommitOnMain() string {
	h, err := g.run("--no-pager", "log", "-n", "1", "--format=%H", g.Repo.Main)
	if err != nil {
		return ""
	}
	return parseOneLine(h)
}

// Get the hash of the initial commit.
func (g *Git) getInitialCommit() string {
	h, err := g.run("rev-list", "--max-parents=0", "HEAD")
	if err != nil {
		return ""
	}
	// It's possible to return more than one result (rarely).
	lines := parseLines(h)
	// The correct commit should be the oldest, which is the last in the list.
	return lines[len(lines)-1]
}

// Get the symbolic ref for the HEAD or an empty string if it isn't symbolic.
func (g *Git) symbolicRefHead() string {
	h, err := g.run("symbolic-ref", "HEAD", "--")
	// e.g.
	//   refs/heads/main
	// or
	//   fatal: ref HEAD is not a symbolic ref
	if err != nil {
		return ""
	}
	return parseOneLine(h)
}

// Get a Ref struct for the HEAD.
func (g *Git) headRef() Ref {
	head := Ref{
		Name: "HEAD",
		Head: true,
	}
	h := g.symbolicRefHead()
	if strings.HasPrefix(h, "refs/heads/") {
		head.Branch = h[11:]
	}

	return head
}

// Get ref details from `git show-ref`, build Refs struct for
// HEAD, branches, tags, remote branches, and heads. This is
// the primary data used to build the labels in the commit list.
func (g *Git) getRefs() (Refs, error) {
	var refs Refs

	show_refs, err := g.run("show-ref", "--dereference", "--head")
	// e.g.
	// a67ea1dbf2b31ebd354604cdc60574950c7fe905 HEAD
	// a67ea1dbf2b31ebd354604cdc60574950c7fe905 refs/heads/main
	// 11522e4f2d94b5043861f5b4e6899ffd5482ac6d refs/heads/test
	// a67ea1dbf2b31ebd354604cdc60574950c7fe905 refs/remotes/origin/HEAD
	// a67ea1dbf2b31ebd354604cdc60574950c7fe905 refs/remotes/origin/main
	// e35785a1e71efbb7a48a1be286236f93f5aeded6 refs/stash
	// e1a3558374dbe85a7eab5094185b1b3e30391f96 refs/tags/testTag
	if err != nil {
		return refs, err
	}

	upstream, err := g.getUpstreamsForRefs()
	if err != nil {
		return refs, err
	}

	refs.HEAD = g.headRef()

	tag_lookup := make(map[string]int)
	// For the purposes of displaying a coherent tree,
	// we're denoting the following:
	//   refs/heads/[...]                 = branches
	//   refs/tags/[...]                  = tags
	//   HEAD and refs/remotes/[...]/HEAD = heads
	//   refs/remotes/[...]/[...]         = remotes
	show_refs_lines := parseLines(show_refs)
	for _, r := range show_refs_lines {
		ref_details := strings.Split(r, " ")
		if len(ref_details) >= 2 {
			hash := ref_details[0]
			name := strings.Join(ref_details[1:], " ")
			if strings.HasPrefix(name, "refs/heads/") {
				ref := parseRefHead(hash, name[11:])
				if up, exists := upstream[name]; exists {
					ref.Upstream = up
				}
				if ref.Branch == refs.HEAD.Branch {
					ref.Head = true
				}
				refs.Branches = append(refs.Branches, ref)
			} else if strings.HasPrefix(name, "refs/tags/") {
				ref := parseRefTag(hash, name[10:])
				if up, exists := upstream[name]; exists {
					ref.Upstream = up
				}
				tag_lookup[name] = len(refs.Tags)
				refs.Tags = append(refs.Tags, ref)
			} else if strings.HasPrefix(name, "refs/remotes/") {
				ref := parseRefRemote(hash, name[13:])
				if up, exists := upstream[name]; exists {
					ref.Upstream = up
				}
				if name[len(name)-4:] == "HEAD" {
					refs.Heads = append(refs.Heads, ref)
				} else {
					refs.RemoteBranches = append(refs.RemoteBranches, ref)
				}
			} else if name == "HEAD" {
				refs.HEAD.Hash = hash
			} else if !strings.HasPrefix(name, "refs/stash") {
				// Ignore stash, but anything else log a warning.
				println("Error parsing ref:", name)
			}
		}
	}

	for n := range refs.RemoteBranches {
		// Add to remote branches which local branches are in sync.
		// Add to local branches which remote branches are in sync.
		for i := range refs.Branches {
			if refs.Branches[i].Hash == refs.RemoteBranches[n].Hash && (refs.Branches[i].Branch == refs.RemoteBranches[n].Branch || refs.Branches[i].Upstream == refs.RemoteBranches[n].Name) {
				refs.Branches[i].SyncedRemotes = append(refs.Branches[i].SyncedRemotes, refs.RemoteBranches[n].Remote)
				refs.RemoteBranches[n].SyncedLocally = true
			}
		}
		// Add to remote branches which remote heads are in sync.
		for h := range refs.Heads {
			if refs.Heads[h].Hash == refs.RemoteBranches[n].Hash {
				refs.RemoteBranches[n].Head = true
			}
		}
		// Add to HEAD which remote branches are in sync.
		if refs.RemoteBranches[n].Hash == refs.HEAD.Hash {
			refs.HEAD.SyncedRemotes = append(refs.HEAD.SyncedRemotes, refs.RemoteBranches[n].Remote)
		}
	}

	// Add to tags which remotes they are on.
	// todo: THIS IS SLOW as it's running ls-remote. Async or faster way of doing this?
	// remote_names, err := g.getRemoteNames()
	// if err == nil && len(remote_names) > 0 {
	// 	for _, remote := range remote_names {
	// 		tags, err := g.getRemoteTags(remote)
	// 		if err == nil && len(tags) > 0 {
	// 			for _, tag := range tags {
	// 				if i, exists := tag_lookup[tag]; exists {
	// 					refs.Tags[i].SyncedRemotes = append(refs.Tags[i].SyncedRemotes, remote)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

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
			Head:   true,
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

func (g *Git) ShowRefAll() (string, error) {
	refs, err := g.run("show-ref", "--dereference", "--head")
	if err != nil {
		if errorCode(err) == NoCommitsYet || errorCode(err) == BadRevision || errorCode(err) == UnknownRevisionOrPath || errorCode(err) == ExitStatus1 {
			return "", nil
		}
		return "", err
	}
	return refs, nil
}

func (g *Git) getUpstreamsForRefs() (map[string]string, error) {
	upstream := make(map[string]string)

	refs, err := g.run("for-each-ref", "--format=%(refname)"+GIT_LOG_SEP+"%(upstream:short)")
	if err != nil {
		if errorCode(err) == NoCommitsYet || errorCode(err) == BadRevision || errorCode(err) == UnknownRevisionOrPath {
			return upstream, nil
		}
		return upstream, err
	}

	// refs/remotes/origin/*
	origin_remotes := []string{}
	// refs/remotes/[not origin]*
	remotes := []string{}

	refs_lines := parseLines(refs)
	for _, line := range refs_lines {
		parts := strings.Split(line, GIT_LOG_SEP)
		if len(parts) == 2 {
			upstream[parts[0]] = parts[1]
		} else {
			upstream[parts[0]] = ""
		}
		if strings.HasPrefix(parts[0], "refs/remotes/origin/") {
			origin_remotes = append(origin_remotes, parts[0][13:])
		} else if strings.HasPrefix(parts[0], "refs/remotes/") {
			remotes = append(remotes, parts[0][13:])
		}
	}

	// If no upstream was found for any heads, loop through the remote refs and see
	// if there's a match. If so, we assume the two are related. Start with origin
	// in order to default to that remote, then loop through any remaining.
	for ref := range upstream {
		if upstream[ref] == "" && strings.HasPrefix(ref, "refs/heads/") {
			name := ref[11:]
			for i := range origin_remotes {
				if strings.HasSuffix(origin_remotes[i], "/"+name) {
					upstream[ref] = origin_remotes[i]
					break
				}
			}
			if upstream[ref] == "" {
				for i := range remotes {
					if strings.HasSuffix(remotes[i], "/"+name) {
						upstream[ref] = remotes[i]
						break
					}
				}
			}
		}
	}

	return upstream, nil
}

// In the case when a branch isn't tracking a remote, search for one that
// matches its name and make the assumption that the two are related.
func (g *Git) findRemoteBranch(branch string) (string, error) {
	r, err := g.run("rev-parse", "--abbrev-ref", "--remotes", "*/"+branch)
	if err != nil {
		return "", err
	}
	remotes := parseLines(r)
	if len(remotes) == 0 {
		return "", nil
	}
	// Loop through and see if there's an origin remote. Default to that.
	for i := range remotes {
		if strings.HasPrefix(remotes[i], "origin/") {
			return remotes[i], nil
		}
	}
	// Otherwise return the first remote branch found.
	return remotes[0], nil
}
