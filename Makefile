#
# Config
#
MODULE_PATH_LAUNCHER    := github.com/setlog/trivrost/cmd/launcher
MODULE_PATH_HASHER      := github.com/setlog/trivrost/cmd/hasher
MODULE_PATH_VALIDATOR   := github.com/setlog/trivrost/cmd/validator
MODULE_PATH_SIGNER      := github.com/setlog/trivrost/cmd/signer
OUT_DIR                 := out
RELEASE_FILES_DIR       := ${OUT_DIR}/release_files
UPDATE_FILES_DIR        := ${OUT_DIR}/update_files
HASHER_BINARY           := hasher
VALIDATOR_BINARY        := validator
SIGNER_BINARY           := signer

# allow custom program name
LAUNCHER_PROGRAM_NAME   := $(shell GO111MODULE=on go run cmd/echo_field/main.go cmd/launcher/resources/launcher-config.json BinaryName)
LAUNCHER_PROGRAM_EXT    := 
LAUNCHER_BRANDING_NAME  := $(shell GO111MODULE=on go run cmd/echo_field/main.go cmd/launcher/resources/launcher-config.json BrandingName)
MSI_PREFIX              := ${LAUNCHER_PROGRAM_NAME}

GITDESC                 := $(shell git describe --tags 2> /dev/null || echo unavailable)
# Version is latest tag corresponding GLOB pattern "v[0-9]*.[0-9]*.[0-9]*"
LAUNCHER_VERSION        := $(shell git describe --tags --abbrev=0 --match "v[0-9]*.[0-9]*.[0-9]*" 2> /dev/null || echo unavailable)
GITBRANCH               := $(shell git symbolic-ref -q --short HEAD || echo unknown)
GITHASH                 := $(shell git rev-parse --short=8 --verify HEAD || echo unknown)
LDFLAGS                 := -s -X main.gitDescription=${GITDESC} -X main.gitBranch=${GITBRANCH} -X main.gitHash=${GITHASH} -X "github.com/setlog/trivrost/cmd/launcher/launcher.buildTime=$(shell date -u "+%Y-%m-%d %H:%M:%S UTC")"

# Assume version is part of the tag. If not, default to v0.0.0
VERSIONOK := $(shell echo -n "${LAUNCHER_VERSION}" | grep -E ^v[0-9]+\.[0-9]+\.[0-9]+$$ && echo ok)
ifndef VERSIONOK
	LAUNCHER_VERSION := v0.0.0
endif
$(info GIT description: '${GITDESC}' (latest master: '${LAUNCHER_VERSION}'), GIT branch '${GITBRANCH}', GIT hash '${GITHASH}')

TIMESTAMP_SERVER = 'http://timestamp.globalsign.com/scripts/timstamp.dll'

#
# Detect OS
#
OS = unknown
OS_UNAME := $(shell uname -s)
ifneq (,$(filter CYGWIN_NT% MSYS_NT% MINGW64_NT%,${OS_UNAME}))
	OS = windows
	ifdef TRIVROST_FORCECONSOLE
$(info TRIVROST_FORCECONSOLE is set, not hiding console)
	else
		LDFLAGS := ${LDFLAGS} -H=windowsgui
	endif
	LAUNCHER_PROGRAM_EXT := .exe
else ifneq (,$(findstring Linux,${OS_UNAME}))
	OS = linux
else ifneq (,$(findstring Darwin,${OS_UNAME}))
	OS = darwin
	export CGO_CFLAGS=-mmacosx-version-min=10.8
	export CGO_LDFLAGS=-mmacosx-version-min=10.8
endif
$(info Detected uname-id '${OS_UNAME}' as OS '${OS}')

# Globally enable go modules
export GO111MODULE=on
export GOOS=${OS}
export LAUNCHER_PROGRAM_NAME
export LAUNCHER_PROGRAM_EXT

#
# Make targets
#
# See https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: build bundle bundle-msi test copy-test-files generate clean sign dist help

# Default target
build: generate  ## Build (default)
ifeq (${OS},windows)
	# Removing unneeded PNG from Windows binary
	rm cmd/launcher/resources/icon.png.gen.go
endif
	# See https://github.com/golang/go/issues/18400#issuecomment-270414574 for why -installsuffix is needed.
	go build -o "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}${LAUNCHER_PROGRAM_EXT}" -v -installsuffix _separate -ldflags '${LDFLAGS}' ${MODULE_PATH_LAUNCHER}
ifeq (${OS},darwin)
	# Mac bundle is special
	mkdir -p "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/MacOS"
	mv "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}" "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/MacOS/launcher"
	cp "cmd/launcher/resources/Info.plist" "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/"
	mkdir -p "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/Resources"
	if [ -f cmd/launcher/resources/icon.icns ]; then cp cmd/launcher/resources/icon.icns "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/Resources/icon.icns"; fi
endif
	$(info # make build finished)

bundle:          ## Bundle OS-specific files. Call after signing
	mkdir -p "${RELEASE_FILES_DIR}/${OS}"
ifeq (${OS},linux)
	tar -cvf "${RELEASE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.tar" -C "${UPDATE_FILES_DIR}/${OS}" .
else ifeq (${OS},darwin)
	# zip is special and needs a cd
	cd "${UPDATE_FILES_DIR}/${OS}"; zip -r "../../../${RELEASE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.zip" "${LAUNCHER_PROGRAM_NAME}.app"
else ifeq (${OS},windows)
	cp "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}" "${RELEASE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}"
endif

# FIXME: The tool always assumes the x64 transform template to be under CWD/build -> parameter? Embed?
DEPLOYMENT_CONFIG ?= trivrost/deployment-config.json
ARCH ?= "386"
bundle-msi: bundle ## Bundle MSI installer packages. Uses cmd/launcher/resources/launcher-config.json. Set DEPLOYMENT_CONFIG to a config with bundles tagged 'msi', ARCH to 386 or amd64 and all public keys of bundles known in public-rsa-keys.pem.
ifneq (${OS},windows)
	$(warning MSI is currently only implemented for windows, skipping)
else
	$(eval OUT_MSI := ${OUT_DIR}/msi/${ARCH})
	mkdir -p "${OUT_MSI}"
	cp "${RELEASE_FILES_DIR}/windows/${LAUNCHER_PROGRAM_NAME}.exe" "${OUT_MSI}/${LAUNCHER_PROGRAM_NAME}.exe"
	echo "Downloading ${ARCH} bundles tagged 'msi'"
	# the launcher detects systemmode by the systembundles directory. It must exist.
	mkdir -p "${OUT_MSI}/systembundles"
	touch "${OUT_MSI}/systembundles/.keep"
	go run cmd/bundown/main.go \
		--deployment-config "${DEPLOYMENT_CONFIG}" \
		--os windows \
		--arch "${ARCH}" \
		--tags "msi" \
		--pub "cmd/launcher/resources/public-rsa-keys.pem" \
		--out "${OUT_MSI}/systembundles"
	echo "Building ${ARCH} MSI packages."
	go run cmd/installdown/main.go \
		--componentgroupdir "${OUT_MSI}" \
		--launcherversion "${LAUNCHER_VERSION}" \
		--launcher-config cmd/launcher/resources/launcher-config.json \
		--arch ${ARCH} \
		--wxstemplate build/launcher.wxs.template \
		--msioutputfile "${MSI_PREFIX}_${ARCH}.msi" \
		--out "${OUT_DIR}"
	cp "${OUT_DIR}/${MSI_PREFIX}_${ARCH}.msi" "${RELEASE_FILES_DIR}/windows/"
endif

tools:           ## Build helper tools like hasher
	go build -o "${OUT_DIR}/${HASHER_BINARY}${LAUNCHER_PROGRAM_EXT}" -v -installsuffix _separate ${MODULE_PATH_HASHER}
	go build -o "${OUT_DIR}/${VALIDATOR_BINARY}${LAUNCHER_PROGRAM_EXT}" -v -installsuffix _separate ${MODULE_PATH_VALIDATOR}
	go build -o "${OUT_DIR}/${SIGNER_BIANRY}${LAUNCHER_PROGRAM_EXT}" -v -installsuffix _separate ${MODULE_PATH_SIGNER}

help:            ## Show this help
	@fgrep -h "##" ${MAKEFILE_LIST} | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

sign:            ## Sign the Windows exe using CERT_FILE (bas64 encoded) and password from CERT_KEY
ifneq (${OS},windows)
	$(warning Signing is currently only implemented for windows, skipping)
else
ifndef CERT_FILE
	$(error CERT_FILE is undefined)
endif
ifndef CERT_KEY
	$(error CERT_KEY is undefined)
endif
	$(info Signing windows release files...)
	echo "$${CERT_FILE}" | base64 -d > ~tmp_launcher_cert
	@signtool sign /debug /a /v /d "${LAUNCHER_BRANDING_NAME}" /f ~tmp_launcher_cert /p "${CERT_KEY}" /t ${TIMESTAMP_SERVER} /fd SHA512 "${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.exe"
	rm ~tmp_launcher_cert
endif

sign-msi:            ## Sign the Windows exe using CERT_FILE (bas64 encoded) and password from CERT_KEY
ifndef CERT_FILE
	$(error CERT_FILE is undefined)
endif
ifndef CERT_KEY
	$(error CERT_KEY is undefined)
endif
ifneq (${OS},windows)
	$(warning Signing is currently only implemented for windows, skipping)
else
	$(info Signing windows release files...)
	echo "$${CERT_FILE}" | base64 -d > ~tmp_launcher_cert
	@signtool sign /debug /a /v /d "${LAUNCHER_BRANDING_NAME}" /f ~tmp_launcher_cert /p "${CERT_KEY}" /t ${TIMESTAMP_SERVER} /fd SHA512 "${RELEASE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}_386.msi"
	@signtool sign /debug /a /v /d "${LAUNCHER_BRANDING_NAME}" /f ~tmp_launcher_cert /p "${CERT_KEY}" /t ${TIMESTAMP_SERVER} /fd SHA512 "${RELEASE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}_amd64.msi"
	rm ~tmp_launcher_cert
endif

TYPE ?= unit
test:           ## Run tests. Set TYPE to 'unit' (default with coverage-report), 'lint' or 'race' (only 64bit)
ifeq (${TYPE},unit)
	mkdir -p out
	# If we use CC=clang for building C code, we could memory-sanatize the code
	#go test -v -installsuffix _separate_test -ldflags '${LDFLAGS}' -msan ./...
	# Race and coverage
	go test -v -installsuffix _separate_test -ldflags '${LDFLAGS}' -covermode=atomic -coverprofile "${OUT_DIR}/coverage_${OS}.cov" ./...
	go tool cover -func="${OUT_DIR}/coverage_${OS}.cov"
	go tool cover -html="${OUT_DIR}/coverage_${OS}.cov" -o "${OUT_DIR}/coverage_${OS}.html"
endif
ifeq (${TYPE},lint)
	# Check formatting issues
	@test -z "$$(gofmt -l .)" || (echo "Following files need gofmt:\n$$(gofmt -l .)" && exit 1)
	# Lint
	golint -set_exit_status ./...
endif
ifeq (${TYPE},race)
	go test -v -installsuffix _separate_test -ldflags '${LDFLAGS}' -race ./...
endif

test-integration: build ## Build and run integration tests
	@echo "Checking whether --build-time still returns valid timestamp."
ifeq (${OS},darwin)
	@"${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}.app/Contents/MacOS/launcher" --build-time | grep --quiet -E "2[0-9]{3}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} UTC"
else
	@"${UPDATE_FILES_DIR}/${OS}/${LAUNCHER_PROGRAM_NAME}${LAUNCHER_PROGRAM_EXT}" --build-time | grep --quiet -E "2[0-9]{3}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2} UTC"
endif
	@echo "Running intergration tests..."
	go test -v -installsuffix _separate_test -ldflags '${LDFLAGS}' -tags=integration ./...


sonar:           ## Runs sonar-scan and sends them to the host configured with SONAR_HOST. Results saved under out/sonar
	SONAR_HOST=$${SONAR_HOST} COVERAGE_FILES=$$(ls -m out/*.cov) PROJECT_VERSION=${LAUNCHER_VERSION} sonar-scanner

copy-test-files: ## Copy example resources into resource directory
	cp examples/launcher-config.json.example cmd/launcher/resources/launcher-config.json
	cp examples/public-rsa-keys.pem.example cmd/launcher/resources/public-rsa-keys.pem
	cp examples/defaulticon.png cmd/launcher/resources/icon.png
	cp examples/defaulticon.ico cmd/launcher/resources/icon.ico

generate:        ## Run go generate
	GO111MODULE=off go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
	go generate -installsuffix _separate -ldflags '${LDFLAGS}' ${MODULE_PATH_LAUNCHER}

clean:           ## Clean generated files
ifndef OUT_DIR
	$(error OUT_DIR is undefined)
endif
	go clean ${MODULE_PATH_LAUNCHER}
	go clean ${MODULE_PATH_HASHER}
	rm -rf "${OUT_DIR}"
	rm -f cmd/launcher/resources/*.gen.go
	rm -f cmd/launcher/*.syso
