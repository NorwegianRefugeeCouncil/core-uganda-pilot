#!/bin/bash

if [ ! -e "install_githooks.sh" ]; then
	echo "This script needs to be run from within its directory"
	exit 1
fi

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
