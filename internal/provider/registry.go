package provider

import (
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	allDataSources []func() datasource.DataSource
	allResources   []func() resource.Resource
)

// registerDataSource may be called during package initialization to register a new data source with the provider.
func registerDataSource(fn func() datasource.DataSource) {
	allDataSources = append(allDataSources, fn)
}

// registerResource may be called during package initialization to register a new resource with the provider.
func registerResource(fn func() resource.Resource) {
	allResources = append(allResources, fn)
}

func splitInternalID(s types.String) (string, string) {
	splitted := strings.Split(s.ValueString(), "/")
	return splitted[0], splitted[1]
}

func DiffTypedStrings(oldList, newList []types.String) (added, removed []types.String) {
	oldMap := make(map[string]struct{})
	newMap := make(map[string]struct{})

	for _, v := range oldList {
		if !v.IsNull() && !v.IsUnknown() {
			oldMap[v.ValueString()] = struct{}{}
		}
	}
	for _, v := range newList {
		if !v.IsNull() && !v.IsUnknown() {
			newMap[v.ValueString()] = struct{}{}
		}
	}

	for _, v := range newList {
		val := v.ValueString()
		if _, found := oldMap[val]; !found {
			added = append(added, v)
		}
	}
	for _, v := range oldList {
		val := v.ValueString()
		if _, found := newMap[val]; !found {
			removed = append(removed, v)
		}
	}

	return
}

// enumValuesList formats enum values as a space-separated list of Markdown code spans.
func enumValuesList[T ~string](values []T) string {
	quoted := make([]string, 0, len(values))
	for _, v := range values {
		quoted = append(quoted, "`"+string(v)+"`")
	}

	return strings.Join(quoted, " ")
}

// enumStrings returns enum values as plain strings.
func enumStrings[T ~string](values []T) []string {
	out := make([]string, 0, len(values))
	for _, v := range values {
		out = append(out, string(v))
	}

	return out
}

// timeValue returns t as an RFC 3339 string value.
func timeValue(t time.Time) types.String {
	return types.StringValue(t.Format(time.RFC3339Nano))
}

// timePointerValue returns t as an RFC 3339 string value, or a null value when t is nil.
func timePointerValue(t *time.Time) types.String {
	if t == nil {
		return types.StringNull()
	}

	return timeValue(*t)
}

// diffStrings returns the sorted, distinct values of newList missing from
// oldList (added) and of oldList missing from newList (removed).
func diffStrings(oldList, newList []string) (added, removed []string) {
	missing := func(from, in []string) []string {
		present := make(map[string]bool, len(in))
		for _, v := range in {
			present[v] = true
		}

		out := []string{}
		for _, v := range from {
			if !present[v] {
				present[v] = true
				out = append(out, v)
			}
		}
		sort.Strings(out)

		return out
	}

	return missing(newList, oldList), missing(oldList, newList)
}
