# cookiecutter-upp-golang

Project template for UPP microservices in Golang.
To be used with the [Cookiecutter](https://github.com/audreyr/cookiecutter) tool.

# Introduction

This template is meant as an up-to-date skeleton for new Go microservices in the Universal Publishing Platform.

Whenever we need to e.g. bump the Docker image version or improve the README.md template or any other change that new services need to pick up, this is the place to make those changes.
  
Current features
* ready-to-run app
* basic admin endpoints with a sample healthcheck
* Docker config
* CircleCI config
* vendoring support
* README.md skeleton
* .gitignore file

# Usage

1. Install cookiecutter (see [documentation](https://cookiecutter.readthedocs.io/en/latest/installation.html) for detailed instructions)

```bash
pip install --user cookiecutter
```

macOS:

```bash
brew install cookiecutter
```

2. Create new local repository

```bash
cd $GOPATH/src/github.com/Financial-Times
cookiecutter git@github.com:Financial-Times/cookiecutter-upp-golang.git
```

You will be prompted to provide values for template properties (or accept defaults).
See [cookiecutter.json](cookiecutter.json) for supported properties.

3. Make use of yor new repo

```bash
cd $GOPATH/src/github.com/<your_new_repo>
go build
```

The app should build and run with basic admin endpoints

