#!/bin/bash

# move to script dir
GIT_HOOKS="$(dirname "$0")/git-hooks"
cd "$GIT_HOOKS"

DEST=$(realpath ../../.git/hooks)

for f in *; do
	# check if symlink or file exists
	if [[ -h "$DEST/$f" || -f "$DEST/$f" ]]; then
		rm "$DEST/$f"
	fi
	# create symlink
	ln -s $f "$DEST/$f"
	if [ $? == 0 ]; then
		echo "Created symlink for $f hook."
	fi
done
