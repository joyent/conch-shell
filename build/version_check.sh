#!/bin/sh

if test -f "VERSION"; then
	VER=$(cat VERSION)
	TAG=$(git describe --tags | sed 's/^v//' | awk -F'-' '{print $1}')

	if test "$VER" != "$TAG"; then
		echo "VERSION ( $VER ) and tag ( $TAG ) do not agree"
		exit 1
	fi
else
	echo "No VERSION file"
	exit 1
fi
exit 0
