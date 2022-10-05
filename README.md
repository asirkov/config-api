# JsonConfig API
API for validation and store configs.
## Development
This project is using [go modules](https://github.com/golang/go/wiki/Modules) for handling dependencies.
Clone the project **outside** your `GOPATH`. 

If you still want to have in your `GOPATH`, you will need to enable the 
modules feature as it is disabled inside of `GOPATH` by default. You can enable it by setting the env variable `GO111MODULE=on`.

### Requirements
* go ([how to install](https://golang.org/doc/install))
* the mysql server running and [configured](config/README.md) properly

### Set up Github with SSH
1. Generate SSH key (if required):
```bash
ssh-keygen -t rsa -b 4096 -C "email@example.com"
ssh-keygen -t ecdsa -b 521 -C "email@example.com"
```
Let's say your private key is located here: ~/.ssh/id_rsa_livescore.

See [documentation](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent):

2. Add SSH key to GitHub ([documentation](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/adding-a-new-ssh-key-to-your-github-account))

3. Clone repository:
```bash
GIT_SSH_COMMAND="ssh -i ~/.ssh/id_rsa_livescore" git clone git@github.com:minelytix/config-api.git
```
And then:
```bash
cd config-api
```

4.Configure repository to use SSH:
```bash
git config url."git@github.com:".insteadOf "https://github.com/"
```

5. Set specific SSH key
```bash
git config core.sshCommand "ssh -i ~/.ssh/id_rsa_livescore -F /dev/null"
```
This step is required if you have multiple SSH keys.

6. Build project first time with:
```bash
GIT_SSH_COMMAND="ssh -i ~/.ssh/id_rsa_livescore" go build
```

### Set up
Setup local environment
```
LS_ENVIRONMENT_NAME=%environment% 
```
`%environment%` should be created from your initials, for example: `OS` for `Oleksandr Sirkov`.

The service will try to read the configuration by the name of your environment - `%environment%.toml`, or the default one - `default.toml`

Just build the service:
```bash
go build
```

Run: 
```bash
./config-api
```

### Running in docker
If you don't want to install go on your machine, you can also use docker to build/run the app.
To do so, run
```bash
docker build -t config-api --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" .
docker run config-api
```
