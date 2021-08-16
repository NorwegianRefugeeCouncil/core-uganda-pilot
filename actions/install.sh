#!/bin/bash

# move to script dir
SRC_DIR="$(realpath "$(dirname $0)")/git-hooks"
echo "$SRC_DIR"
cd "$SRC_DIR" || exit


for f in *; do
	echo "$f"
	DEST="../../.git/hooks/$f"
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

