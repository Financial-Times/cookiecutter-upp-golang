#!/bin/bash

echo "mode: atomic" > $CIRCLE_ARTIFACTS/coverage.out # covermode is set to atomic when running tests using the -race flag

PACKAGES=`govendor list -no-status +local`
PACKAGES_LIST=`govendor list -no-status +local |  paste -sd ',' -`
EXIT_CODE=0

for PKG in $PACKAGES; do
  echo =-= $PKG

  govendor test -v -race -coverpkg=$PACKAGES_LIST -coverprofile=$CIRCLE_ARTIFACTS/profile.out $PKG; __EXIT_CODE__=$?

  if [ "$__EXIT_CODE__" -ne "0" ]; then
    EXIT_CODE=$__EXIT_CODE__
  fi

  if [ -f $CIRCLE_ARTIFACTS/profile.out ]; then
    tail -n +2 $CIRCLE_ARTIFACTS/profile.out >> $CIRCLE_ARTIFACTS/coverage.out; rm $CIRCLE_ARTIFACTS/profile.out
  fi
done

exit $EXIT_CODE