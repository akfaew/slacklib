TEST_ARGS = -failfast

update:
	go get -u
	go mod tidy
	go mod verify

fmt:
	go fmt ./...

test: fmt
	go test $(TEST_ARGS) ./...

test-regen:
	rm -rf testdata/output
	mkdir -p testdata/output
	go test $(TEST_ARGS) -regen .

test-cover: fmt
	go test $(TEST_ARGS) -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out |\
		grep -v 100.0% |\
		grep -v total: |\
		perl -nae 'printf("%7s %s %s\n", $$F[2], $$F[0], $$F[1])' | sort -nr
	go tool cover -html=coverage.out

check-workspace: fmt
	# check if workspace is clean
	@if ! git diff-index --quiet HEAD --; then \
		echo "You have unstaged changes"; \
		git status; \
		exit 1; \
	fi

push: test
	git push
	git push --tags

release: check-workspace test push

clean:
	rm -f coverage.out
