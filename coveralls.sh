#!/bin/bash

echo "mode: set" > acc.out
for Dir in $(find ./* -maxdepth 0 -type d -and -not -iname "Godeps");
do
	if ls $Dir/*.go &> /dev/null;
	then
		returnval=`go test -coverprofile=profile.out $Dir`
		echo ${returnval}
		if [[ ${returnval} != *FAIL* ]]
		then
    		if [ -f profile.out ]
    		then
        		cat profile.out | grep -v "mode: set" >> acc.out 
    		fi
    	else
    		exit 1
    	fi	
    fi
done
if [ -n "$COVERALLS_TOKEN" ]
then
	goveralls -coverprofile=acc.out -v -service drone.io -repotoken="$COVERALLS_TOKEN" 
fi	

rm -rf ./profile.out
rm -rf ./acc.out
