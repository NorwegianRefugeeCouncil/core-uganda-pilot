#!/bin/bash

if [ ! -e "install.sh" ]; then
  echo "This script needs to be run from within its directory"
  exit 1
fi
gitconfig='../.git/config'
# remove hooksPath from git config if set
grep -v "hooksPath" $gitconfig > __tmp && mv __tmp $gitconfig
# point hooksPath to hooks dir
sed -i'' "s/\[core\]/\[core\]\n\thooksPath = actions\/git-hooks/" $gitconfig
echo "Updated repo git config successfully."
