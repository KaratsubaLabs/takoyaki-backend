
# takoyaki backend

go backend for からつばLABS's **project takoyaki** - the vps platform.

## RUNNING FOR DEVELOPMENT

If you wish to run **takoyaki** without having to keep building docker
containers, you can run it locally instead. Make sure you have a working go
installation.

Make your own copy of `.env` by copying the provided `dotenv.example`
file.
```
$ cp dotenv.example .env
```

Install packages
```
$ go mod download
```

Run takoyaki
```
$ go run *.go
```

## RUNNING ALL CONTAINERS

To be able to run the stack, **docker** and **docker-compose** are required.
Consult the relevant documentation based on your system on how to get these set
up.

First, make your own copy of `.env` by copying the provided `dotenv.example`
file.
```
$ cp dotenv.example .env
```

Next we can start the containers
```
$ docker-compose up
```

## RESETING THE DATABASE

During testing, if it happens that you wish to reset the database, simply
remove the directory:
```
$ sudo rm -rf db/data/
```

## TODO

- [ ] database initialization migration (+ shell interface to init db)
- [x] validation for requests (as a middleware if possible)
- [x] possibly error middleware
- [x] figure out where to put temp files (cidata.iso etc) for when creating vps
- [x] jwt auth
- [ ] possibly create db struct so methods can all be namespaced
- [x] look at database transactions (+ are they really needed)
- [ ] write the routes
- [ ] allow optional args when configuring vps (ie ssh key)
- [ ] add tests? (might be overkill + annoying)

