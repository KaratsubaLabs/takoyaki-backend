
INSTALL_DIR=/usr/local/bin
SERVICE_DIR=/etc/systemd/system

.PHONY: install uninstall

install:
	cp takoyaki-pipe ${INSTALL_DIR}
	cp takoyaki-pipe.service ${SERVICE_DIR}

uninstall:
	rm -f ${INSTALL_DIR}/takoyaki-pipe
	rm -f ${SERVICE_DIR}/takoyaki-pipe.service

