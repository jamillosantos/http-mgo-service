COVERDIR=$(CURDIR)/.cover
COVERAGEFILE=$(COVERDIR)/cover.out
COVERAGEREPORT=$(COVERDIR)/report.html

setup:
	@docker-compose exec mongo mongo test-service-database --eval 'db.createUser({user:"snake.eyes",pwd:"123456",roles:["readWrite"], passwordDigestor: "server"});'

test:
	@ginkgo --failFast ./...

test-watch:
	@ginkgo watch --debug -cover -r ./...

coverage-ci:
	@mkdir -p $(COVERDIR)
	@ginkgo -r -covermode=count --cover --trace ./
	@echo "mode: count" > "${COVERAGEFILE}"
	@find . -type f -name *.coverprofile -exec grep -h -v "^mode:" {} >> "${COVERAGEFILE}" \; -exec rm -f {} \;

coverage: coverage-ci
	@sed -i -e "s|_$(CURDIR)/|./|g" "${COVERAGEFILE}"

coverage-html: coverage
	@go tool cover -html="${COVERAGEFILE}" -o $(COVERAGEREPORT)
	@xdg-open $(COVERAGEREPORT) 2> /dev/null > /dev/null

.PHONY: test test-watch coverage coverage-ci coverage-html