build:
	docker build -t andreaswachs/placeholder_vuln_server:latest -f Dockerfile.vuln_server .

run:
	docker run -d -p 8080:8080 -l placeholder_vuln_server andreaswachs/placeholder_vuln_server:latest

push:
	docker push andreaswachs/placeholder_vuln_server:latest

publish: build push

down:
	docker ps -q -f label=placeholder_vuln_server | xargs docker kill
