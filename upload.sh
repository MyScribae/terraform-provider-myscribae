#!/bin/bash

# Constants
CURRENT_TIME=$(date +%Y%m%d-%H%M%S)
TEMP_DIR="$(pwd)/tmp_$CURRENT_TIME"
CWD="$(pwd)"
DEPLOYMENTS_BUCKET="myscribae-deployments-bucket"
DEPLOYMENT_FOLDER="dev-backend"

# Script specific 
ZIP_ARGS=""
ZIP_FOLDER="."
OUTPUT_EXE="bootstrap"
ZIP_CONTENTS="$OUTPUT_EXE"
CLEAN_UP="$OUTPUT_EXE"

# Args
ZIP_ONLY=false
RELEASE_TYPE="patch"
S3_OR_GITHUB_RELEASE="gh"

# Collect args from command line
while [[ $# -gt 0 ]]; do
	KEY="$1"

	case $KEY in
		--zip-only|-z)
		ZIP_ONLY=true
		shift
		;;
		--release-type|-r)
		RELEASE_TYPE="$2"
		shift
		shift
		;;
		--skip-cleanup|-s)
		CLEAN_UP=""
		shift
		;;
		*)
		echo "Unknown option: ${KEY}"
		shift
		;;
	esac
done

build() {
	GOOS=linux GOARCH=amd64 go build -o "$OUTPUT_EXE" main.go
}

zip_build() {
	echo "Zipping files..."

	# Zip folder is not '' or '.', then cd into it
	if [[ $ZIP_FOLDER != "" && $ZIP_FOLDER != "." ]]; then
		cd $ZIP_FOLDER
	fi

	ls -la 

	# zip files of temp directory into lambda.zip
	echo "zip $ZIP_ARGS deployment.zip $ZIP_CONTENTS"
	zip $ZIP_ARGS deployment.zip $ZIP_CONTENTS

	# If zip fails, exit
	if [ $? -ne 0 ]; then
		echo "Failed to zip files"
		exit 1
	fi
}

upload() {
	if $ZIP_ONLY; then
		echo "Skipping upload to S3"
	else
		if [[ $S3_OR_GITHUB_RELEASE == "gh" ]]; then
			echo "Releasing to GitHub..."

			# Get current release version
			CURRENT_RELEASE=$(gh release list | head -n 1 | awk '{print $1}')
			
			if [[ $CURRENT_RELEASE == "" ]]; then
				CURRENT_RELEASE="0.0.0"
			fi

			# Remove the 'v' prefix from the $CURRENT_RELOEASE
			CURRENT_RELEASE=$(echo $CURRENT_RELEASE | sed 's/v//')

			NEXT_RELEASE=$(pysemver bump $RELEASE_TYPE $CURRENT_RELEASE)

			# Tag current commit with the new release version
			git tag -a v$NEXT_RELEASE -m "Release v$NEXT_RELEASE"

			# Push the tag to GitHub
			git push origin v$NEXT_RELEASE
			
			# Create a new release
			gh release create v$NEXT_RELEASE --generate-notes

			# Upload the deployment.zip to the release
			gh release upload v$NEXT_RELEASE "./deployment.zip"

		elif [[ $S3_OR_GITHUB_RELEASE == "s3" ]]; then
			echo "Uploading to S3..."
			aws s3 cp "./deployment.zip" "s3://$DEPLOYMENTS_BUCKET/$DEPLOYMENT_FOLDER/deployment.zip"
		fi

		rm -f "./deployment.zip"
	fi 
}

cleanup() {
	echo "Cleaning up..."
	if [[ $CLEAN_UP != "" ]]; then
		rm -rf $CLEAN_UP
	fi
}

main() {
	echo "Building..."
	build
	cd "$CWD"

	zip_build
	upload
	cleanup
}

main