lib/build: templ/build tailwind/build 
	go build ./...

tailwind/build:
	npx @tailwindcss/cli -i ./static/tailwind.css -o ./static/dist/style.css --minify

templ/build:
	templ generate -v

