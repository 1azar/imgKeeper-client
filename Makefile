run1:
	go run main.go --addr=":44044" -f=imgs/cat1.png -m=upload

run2:
	go run main.go --addr=":44044" -f=imgs/cat2.png

run3:
	go run main.go --addr=":44044" -f=imgs/cat3.png

run4:
	go run main.go --addr=":44044" -f=cat1.png method=download