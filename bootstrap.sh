#!/bin/sh

set -e

if [ -f /etc/debian_version ]; then
    for package in golang; do
        if [ "$(dpkg --status -- $package 2>/dev/null|sed -n 's/^Status: //p')" != "install ok installed" ]; then
            # add a space after old values
            missing="${missing:+$missing }$package"
        fi
    done
    if [ -n "$missing" ]; then
        echo "$0: missing required DEB packages. Installing via sudo." 1>&2
        sudo apt-get -y install $missing
    fi
elif [ -f /etc/fedora-release ]; then
    for package in golang; do
        if [ "$(rpm -qa $package 2>/dev/null)" == "" ]; then
            missing="${missing:+$missing }$package"
        fi
    done
    if [ -n "$missing" ]; then
        echo "$0: missing required RPM packages. Installing via sudo." 1>&2
        sudo yum -y install $missing
    fi
elif [ -f /etc/redhat-release ]; then
    for package in golang; do
        if [ "$(rpm -qa $package 2>/dev/null)" == "" ]; then
            missing="${missing:+$missing }$package"
        fi
    done
    if [ -n "$missing" ]; then
        echo "$0: missing required RPM packages. Installing via sudo." 1>&2
        sudo yum -y install $missing
    fi
fi
