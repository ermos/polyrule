# 🛠️ Polyrule
> One file to rule them all - Compile validator rules into multiple languages

Polyrule is a versatile command-line interface tool
that streamlines the process of creating validation
rules for data across multiple programming languages.
With Polyrule, you can easily define and manage a wide
range of validation rules, from simple data type checks
to complex business logic rules, in a single place !

[![Go Reference](https://pkg.go.dev/badge/github.com/ermos/polyrule.svg)](https://pkg.go.dev/github.com/ermos/polyrule)
[![Latest tag](https://img.shields.io/github/v/tag/ermos/polyrule?label=latest)](https://github.com/ermos/polyrule/tags)
[![Go Report Card](https://goreportcard.com/badge/github.com/ermos/polyrule)](https://goreportcard.com/report/github.com/ermos/polyrule)
[![Maintainability](https://api.codeclimate.com/v1/badges/c39c1d80ace4bb344393/maintainability)](https://codeclimate.com/github/ermos/polyrule/maintainability)

## 📦 Installation

### Linux

```bash
curl -sL https://github.com/ermos/polyrule/releases/latest/download/polyrule_Linux_$(uname -m).tar.gz | tar -xvz --wildcards 'polyrule' \
&& sudo mv polyrule /usr/local/bin/
```

### Mac

```bash
curl -sL https://github.com/ermos/polyrule/releases/latest/download/polyrule_macOS_all.tar.gz | tar -xvz --wildcards 'polyrule' \
&& sudo mv polyrule /usr/local/bin/
```

### Windows

Download the right archive from [the latest release page](https://github.com/ermos/polyrule/releases/latest).

### Alternative

From `go install` :
```bash
go install github.com/ermos/polyrule@latest
```

From `source` :
```bash
git clone git@github.com:ermos/polyrule.git
make build/bin
```

## 📚 Documentation

You can find the documentation
[here](https://polyrule.smiti.fr/).

## 🤝 Contributing

Contributions to `polyrule` are always welcome!
If you find a bug or have a feature request, please open an issue on GitHub.
If you want to contribute code, please fork the repository and submit a pull request.