# password

Password generator and password file entry manager

This utility will either generate a password and print it to stdout or modify a
password file. The passwords are stored in bcrypt hashed format.

A build task for use with https://taskfile.dev/#/ is included.

Using the -h flag will give usage information.

This application uses two password oriented libraries:

## go-password:

- a password generating library
- https://github.com/sethvargo/go-password - MIT licence

## htpasswd:

- https://github.com/foomo/htpasswd - MIT licence
- A tool for reading and writing htpasswd format files. The library supports apr1, sha, and bcrypt. Only bcrypt is
  supported.

