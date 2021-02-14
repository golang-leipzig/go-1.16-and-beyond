README.html: README.md
	@pandoc --from=gfm+smart --to=html5 --standalone --css='https://rawcdn.githack.com/edwardtufte/tufte-css/e225f3a0e5f8f42a1d28416c1c85962411711907/tufte.min.css' --metadata 'title=Go-1.16-and-beyond' $< -o $@

open: README.html
	xdg-open README.html

