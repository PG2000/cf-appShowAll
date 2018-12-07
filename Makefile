buildAndExecute:
	go build appShowAll.go && cf install-plugin -f appShowAll && cf app-show-all

acceptance:
	go build appShowAll.go && $$GOPATH/bin/ginkgo
