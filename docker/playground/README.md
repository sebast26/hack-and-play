# Docker: Commands and arguments

These are notes that let me understand the difference between Docker ENTRYPOINT and CMD instuctions. 

## Dockerfile.sleep5

This dockerfile uses ENTRYPOINT to execute process. ENTRYPOINT could be overriden when using `docker run`:

```bash
docker run --entrypoint=date img
```

## Dockerfile.cmd5

This dockerfile uses CMD to execute `sleep 5` command. To override the command you can simply execute:

```bash
docker run sebast26/cmd5 date
```

If you want to sleep for more than 5 seconds and you have `CMD sleep 5` in your Dockerfile, then you can execute the docker run with:

```bash
docker run sebast26/cmd5 sleep 10
```

It's easier to modify command than using ENTRYPOINT, but in this setup you still need to use `sleep 10` as a command to docker run.

## Dockerfile.entrypoint

If you use Dockerfile (with ENTRYPOINT ["sleep"]) when trying to run the container with simply `docker run sebast26:entrypoint` you will get an error, that `sleep` needs an argument.

To use this docker you need to specify the argument in docker run:

```bash
docker run sebast26/entrypoint 10
```

and the ENTRYPOINT command and argument to the docker run will be concatenated. In the result the command that will be run at startup is `sleep 10`.

## ENTRYPOINT vs. CMD based on previous examples

One of the difference between ENTRYPOINT and CMD is that when you execute `docker run img command`:
- when using `CMD` in Dockerfile to `CMD` will be overriden by `command` from docker run
- when using `ENTRYPOINT` in Dockerfile the `command` from docker run will be **appended** to ENTRYPOINT

## Dockerfile.both

When you want to have "default" value for you command when using ENTRYPOINT, you could combine both ENTRYPOINT and CMD.

When you do not specify the `command` in docker run then the command executed at startup will be ENTRYPOINT + CMD (in our case `sleep 5`).

When you specify the `command` part in docker run then the `CMD` from dockerfile will be replaced by `command` and the ENTRYPOINT + command will be executed at startup.

## The use of `ENV`

You can quickly parametrise running docker command with `ENV` instruction. In `Dockerfile.env` file we have an example on how to override default value printed to the user.

You can use `ENV` instruction and then execute docker run with:

```bash
docker run -e"NAME=Ania" sebast26/env
```

