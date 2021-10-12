#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

gitconfig='../.git/config'

# remove hooksPath from git config if set
grep -v "hooksPath" $gitconfig > __tmp && mv __tmp $gitconfig

# point hooksPath to hooks dir
uname=$(uname)
if [ "$uname" == "Linux" ]; then
	sed -i'' "s/\[core\]/\[core\]\n\thooksPath = githooks/" $gitconfig
elif [ "$uname" == "Darwin" ]; then
	# add a space after the '-i' flag on mac...
	sed -i '' "s/\[core\]/\[core\]\n\thooksPath = githooks/" $gitconfig
fi

echo "Installed githooks successfully."
