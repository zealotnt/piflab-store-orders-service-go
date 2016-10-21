#!/bin/bash -e

echo "mode: atomic" > piflab-store-api-go.coverprofile

packages=(
	"handlers"
	"models"
	"models/form"
	"models/repository"
	"services"
)

if [ "$COVERALLS_TOKEN" == "" ]
then
	reset
	# ginkgo -r -cover -skipPackage=handlers,models,services,repository,form	# skip all
	ginkgo -r -cover
fi

for package in ${packages[@]};
do
	path=./$package/$(basename $package).coverprofile
	if [ -f $path ]; then
		cat $path | grep -v "mode: atomic" >> piflab-store-api-go.coverprofile
		rm $path
	fi
done

go tool cover -func=piflab-store-api-go.coverprofile
go tool cover -html=piflab-store-api-go.coverprofile -o piflab-store-api-go.coverprofile.html

if [ -n "$COVERALLS_TOKEN" ]
then
	goveralls -coverprofile=piflab-store-api-go.coverprofile -service circleci -repotoken $COVERALLS_TOKEN
	rm ./piflab-store-api-go.coverprofile
fi
