#!/bin/bash

# move to script dir
SRC_DIR="$(realpath "$(dirname "$0")")/git-hooks"
GIT_DIR="$(realpath "$(dirname "$(dirname "$0")")")/.git"
# remove hooksPath from git config if set
grep -v "hooksPath" $GIT_DIR/config > __tmp && mv __tmp $GIT_DIR/config
cd "$SRC_DIR" || exit
for f in *; do
	DEST="$GIT_DIR/hooks/$f"
	# ensure hook is executable
	chmod +x $f
	# check if symlink or file exists
	if [ -h "$DEST" ] || [ -f "$DEST" ]; then
		echo "Removing old file."
		rm "$DEST"
	fi
	# create symlink
	if ln -s -f "$SRC_DIR/$f" "$DEST"; then
		echo "Created symlink for $f hook."
	fi
done

