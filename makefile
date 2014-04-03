all:
	go build aulasis
	mv aulasis dist/

clean:
	rm aulasis
