all: run

clean:
	rm -rf dist

run:
	sass --watch scss:static/css &
	go run *.go

dist: clean
	sass scss/app.scss static/css/app.css
	go build -o dist/visualizer
	upx dist/visualizer
	mkdir -p dist
	cp -r games gohtml static dist

deploy: dist
	rsync -av --delete dist/ root@phinau.de:/opt/visualizer
