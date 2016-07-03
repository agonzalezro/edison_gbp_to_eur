USER=root
IP=192.168.1.133
BINARY=edison_gbp_to_eur
REMOTE_PATH=/home/root

all:
	make build
	make deploy
	make clean
	make run

build:
	GOOS=linux GOARCH=386 go build

deploy:
	scp $(BINARY) $(USER)@$(IP):$(REMOTE_PATH)

clean:
	rm $(BINARY)

run:
	ssh $(USER)@$(IP) "$(REMOTE_PATH)/$(BINARY)"

install:
	scp edison_gbp_to_eur.service $(USER)@$(IP):/etc/systemd/system
	ssh $(USER)@$(IP) "systemctl enable $(BINARY)"
