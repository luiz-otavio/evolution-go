#!/bin/bash
set -euo pipefail

# Finds packages that are affected by changes based on changed files between commits.
# For Go workspaces with independent modules, it identifies which packages changed.
#
# Usage: ./find-affected-services.sh <base_sha>
#   base_sha: The commit SHA to compare against (e.g., main, HEAD~1, abc123)

BASE_SHA="${1:-}"

if [[ -z "$BASE_SHA" ]]; then
    echo "Usage: $0 <base_sha>" >&2
    echo "  base_sha: Git ref to compare against (e.g., main, HEAD~1)" >&2
    exit 1
fi

# Verify Go is available
if ! command -v go &> /dev/null; then
    echo "Error: 'go' command not found in PATH" >&2
    exit 1
fi

# Discover all modules by finding go.mod files (excluding root if it's a workspace)
echo "Discovering modules..." >&2
module_dirs=$(find . -name "go.mod" -not -path "./go.mod" -exec dirname {} \; | sort -u)

if [[ -z "$module_dirs" ]]; then
    echo "No modules found in workspace" >&2
    exit 0
fi

echo "Found modules:" >&2
echo "$module_dirs" | sed 's/^/  /' >&2
echo "" >&2

# Get list of changed .go files between base and HEAD
changed_files=$(git diff --name-only "$BASE_SHA" HEAD -- '*.go' 'go.mod' 'go.sum' 2>/dev/null || git diff --name-only "$BASE_SHA" -- '*.go' 'go.mod' 'go.sum' 2>/dev/null || true)

if [[ -z "$changed_files" ]]; then
    echo "No Go files changed between $BASE_SHA and HEAD" >&2
    exit 0
fi

echo "Changed files:" >&2
echo "$changed_files" | sed 's/^/  /' >&2
echo "" >&2

# Convert changed files to their package paths
# Extract unique directories from changed files
changed_dirs=$(echo "$changed_files" | while read -r file; do
    dirname "$file"
done | sort -u | grep -v '^\.$')

echo "Changed directories:" >&2
echo "$changed_dirs" | sed 's/^/  /' >&2
echo "" >&2

# For each changed directory, get the Go package import path
# We need to cd into the module directory first to get the correct import path
affected_packages=""
while read -r dir; do
    if [[ -n "$dir" && -d "$dir" ]]; then
        # Find which module this directory belongs to
        for mod_dir in $module_dirs; do
            mod_dir_clean="${mod_dir#./}"
            if [[ "$dir" == "$mod_dir_clean"* ]]; then
                # Get package path by running go list from within the module
                pkg=$(cd "$mod_dir_clean" && go list -e "./${dir#$mod_dir_clean}" 2>/dev/null | head -1 || true)
                if [[ -n "$pkg" && "$pkg" != "." ]]; then
                    affected_packages="${affected_packages}${pkg}"$'\n'
                fi
                break
            fi
        done
    fi
done <<< "$changed_dirs"

affected_packages=$(echo "$affected_packages" | grep -v '^$' | sort -u)

if [[ -z "$affected_packages" ]]; then
    echo "No valid Go packages found in changed files" >&2
    exit 0
fi

echo "Affected packages:" >&2
echo "$affected_packages" | sed 's/^/  /' >&2

# Output affected packages to stdout
echo "$affected_packages"
