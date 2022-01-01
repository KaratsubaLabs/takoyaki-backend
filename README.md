
# takoyaki backend

go backend for からつばLABS' **project takoyaki** - the vps platform.
**takoyaki** backend is a go api server that comes with a cli to do some
administrative tasks like approve/decline vps requests.

## RUNNING THE STACK

To be able to run the stack, **docker** and **docker-compose** are required.
Consult the relevant documentation based on your system on how to get these set
up.

First, make your own copy of `.env` by copying the provided `dotenv.example`
file.
```
$ cp dotenv.example .env
```

**takoyaki-pipe**, the systemd service that executes libvirt commands for the
**takoyaki** container must also be installed and ran:
```
$ cd takoyaki-pipe
$ make install
$ systemctl enable takoyaki-pipe
$ systemctl start takoyaki-pipe
```

Now start the containers (this basically just runs docker-compose up):
```
$ ./scripts/init
```
and migrate the database
```
$ ./scripts/takocli db migrate
```

To stop the stack, simply run
```
$ ./scripts/stop
```
or to purge the entire stack (deletes db data and vms), run:
```
$ ./scripts/purge
```

## TODO

- [x] database initialization migration (+ shell interface to init db)
- [x] validation for requests (as a middleware if possible)
- [x] possibly error middleware
- [x] figure out where to put temp files (cidata.iso etc) for when creating vps
- [x] jwt auth
- [x] look at database transactions (+ are they really needed)
- [x] write the routes
- [x] allow optional args when configuring vps (ie ssh key)
- [ ] ~~add tests? (might be overkill + annoying)~~
- [ ] get progress of creating vps to show to frontend (possibly)
- [x] snapshot requests
- [ ] look into using RLS
- [ ] return ip address and state of vm as well in vps info endpoint + look into vm networking
- [x] execute vm commands on host
- [x] refactor project into multiple modules
- [x] rewrite api
- [ ] move virt specific commands to host side executable
- [ ] vps status
- [ ] possibly dockerize libvirt??
- [ ] use proper go project structure with cmd/ and pkg/ dirs

