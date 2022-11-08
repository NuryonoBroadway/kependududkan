build:
	env GOOS=windows GOARCH=amd64 go build .

add:
	git add .

commit:
	git commit -m $(message)

push:
	git push -uf origin main


PHONY: build add commit push