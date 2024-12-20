#!/bin/bash
# Replaces the '-' in the version number with a '.'.
version=`grep "current_version =" .bumpversion.toml | sed --expression 's/-/./g' | cut -f 3 -d " " | sed "s/[^a-zA-Z\.0-9]//g"`
echo -n $version
