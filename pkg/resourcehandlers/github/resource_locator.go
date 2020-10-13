package github

import (
	"fmt"
	"strings"

	"github.com/gardener/docforge/pkg/util/urls"
)

// ResourceType is an enumeration for GitHub resource types
// Supported types are "tree", "blob" and "wiki"
type ResourceType int

func (s ResourceType) String() string {
	return [...]string{"tree", "blob", "wiki"}[s]
}

// NewResourceType creates a ResourceType enum from string
func NewResourceType(resourceTypeString string) (ResourceType, error) {
	switch resourceTypeString {
	case "tree":
		return Tree, nil
	case "blob":
		return Blob, nil
	case "wiki":
		return Wiki, nil
	}
	return 0, fmt.Errorf("Unknown resource type string %s. Must be one of %v", resourceTypeString, []string{"tree", "blob", "wiki"})
}

const (
	// Tree is GitHub tree objects resource type
	Tree ResourceType = iota
	// Blob is GitHub blob objects resource type
	Blob
	// Wiki is GitHub Wiki resource type
	Wiki
)

var nonSHAPathPrefixes = map[string]struct{}{
	"releases": struct{}{},
	"issues":   struct{}{},
	"issue":    struct{}{},
	"pulls":    struct{}{},
	"pull":     struct{}{},
	"wiki":     struct{}{},
}

// ResourceLocator is an abstraction for GitHub specific Universal Resource Locators (URLs)
// It is an internal structure breaking down the GitHub URLs into more segment types such as
// Repo, Owner or SHA.
// ResourceLocator is a common denominator used to translate between GitHub user-oriented urls
// and API urls
type ResourceLocator struct {
	Host  string
	Owner string
	Repo  string
	SHA   string
	Type  ResourceType
	Path  string
	// branch name (master), tag (v1.2.3), commit hash (1j4h4jh...)
	SHAAlias string
}

// String produces a GitHub website link to a resource from a ResourceLocator.
// That's the format used to link а GitHub resource in the documentation structure and pages.
// Example: https://github.com/gardener/gardener/blob/master/docs/README.md
func (r *ResourceLocator) String() string {
	if r.Type == Wiki {
		return fmt.Sprintf("https://%s/%s%s%s%s", r.Host, r.Owner, "/"+r.Repo, fmt.Sprintf("/%s", r.Type), "/"+r.Path)
	}
	if len(r.SHAAlias) > 0 && len(r.Path) < 1 {
		return fmt.Sprintf("https://%s/%s%s", r.Host, r.Owner, "/"+r.Repo)
	}
	if len(r.SHAAlias) < 1 && len(r.Path) > 0 {
		return fmt.Sprintf("https://%s/%s%s%s", r.Host, r.Owner, "/"+r.Repo, "/"+r.Path)
	}
	return fmt.Sprintf("https://%s/%s%s%s%s%s", r.Host, r.Owner, "/"+r.Repo, fmt.Sprintf("/%s", r.Type), "/"+r.SHAAlias, "/"+r.Path)
}

// GetName returns the Name segment of a resource URL path
func (r *ResourceLocator) GetName() string {
	if len(r.Path) == 0 {
		return ""
	}
	p := strings.Split(r.Path, "/")
	return p[len(p)-1]
}

// Parse a GitHub URL into an incomplete ResourceLocator, without
// the APIUrl property.
func parse(urlString string) (*ResourceLocator, error) {
	var (
		resourceType       ResourceType
		path               string
		err                error
		resourceTypeString string
		shaAlias           string
		u                  *urls.URL
	)

	if u, err = urls.Parse(urlString); err != nil {
		return nil, err
	}

	host := u.Host
	sourceURLSegments := strings.Split(u.Path, "/")

	owner := sourceURLSegments[1]
	repo := sourceURLSegments[2]

	if len(sourceURLSegments) > 3 {
		resourceTypeString = sourceURLSegments[3]
		// {blob|tree|wiki}
		if resourceType, err = NewResourceType(resourceTypeString); err == nil {
			if resourceTypeString != Wiki.String() {
				// that would be wrong url but we make up for that
				if len(sourceURLSegments) < 5 {
					shaAlias = "master"
				} else {
					shaAlias = sourceURLSegments[4]
				}
			}
			s := strings.Join([]string{owner, repo, resourceTypeString, shaAlias}, "/")
			// get the github url "path" part without:
			// - leading "/"
			// - owner, repo and {tree|blob}, shaAlias segments if applicable
			if p := strings.Split(u.Path[1:], s); len(p) > 1 {
				path = strings.TrimPrefix(p[1], "/")
			}
		}
		if err != nil {
			s := strings.Join([]string{owner, repo}, "/")
			if p := strings.Split(u.Path[1:], s); len(p) > 1 {
				path = strings.TrimPrefix(p[1], "/")
			}
		}
	} else {
		resourceType = Tree
		resourceTypeString = Tree.String()
		shaAlias = "master"
	}
	if len(u.Fragment) > 0 {
		path = fmt.Sprintf("%s#%s", path, u.Fragment)
	}
	//TODO: add queries if any
	//TODO: type will always default to 0 (Tree). Introduce nil
	ghRL := &ResourceLocator{
		host,
		owner,
		repo,
		"",
		resourceType,
		path,
		shaAlias,
	}
	return ghRL, nil
}
