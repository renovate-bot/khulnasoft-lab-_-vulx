#!/bin/bash

RELEASES=(wheezy jessie stretch buster trusty xenial bionic)

cd vul-repo/deb

for release in ${RELEASES[@]}; do
  echo "Adding deb package to $release"
  reprepro -A i386 remove $release vul
  reprepro -A amd64 remove $release vul
  reprepro includedeb $release ../../dist/*Linux-64bit.deb
  reprepro includedeb $release ../../dist/*Linux-32bit.deb
done

git add .
git commit -m "Update deb packages"
git push origin master