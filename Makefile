run1:
	go run main.go --addr=":44044" -f=imgs/cat1.png --method=upload

run2:
	go run main.go --addr=":44044" -f=imgs/cat2.png --method=upload

run3:
	go run main.go --addr=":44044" -f=imgs/cat3.png --method=upload

run4:
	go run main.go --addr=":44044" -f=cat1.png --method=download

run5:
	go run main.go --addr=":44044" -f=cat1.png --method=download

run6:
	go run main.go --addr=":44044" -f=cat2.png --method=download

run7:
	go run main.go --addr=":44044" -f=cat3.png --method=download

run8:
	go run main.go --addr=":44044" -f=cat3.png --method=list


