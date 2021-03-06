# Czarcoin V3 Network

[![Go Report Card](https://goreportcard.com/badge/github.com/czarcoin/czarcoin)](https://goreportcard.com/report/github.com/czarcoin/czarcoin)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/czarcoin/czarcoin)
[![Coverage Status](https://coveralls.io/repos/github/czarcoin/czarcoin/badge.svg?branch=master)](https://coveralls.io/github/czarcoin/czarcoin?branch=master)
![Alpha](https://img.shields.io/badge/version-alpha-green.svg)

<img src="https://github.com/czarcoin/czarcoin/raw/master/resources/czr.png" width="100">

Czarcoin is building a decentralized cloud storage network. 
[Check out our white paper for more info!](https://czarcoin.org/white-paper)

----

Czarcoin is an S3-compatible platform and suite of decentralized applications that
allows you to store data in a secure and decentralized manner. Your files are
encrypted, broken into little pieces and stored in a global decentralized
network of computers. Luckily, we also support allowing you (and only you) to
retrieve those files!

## Table of Contents

- [Contributing](#contributing-to-czarcoin)
- [Start using Czarcoin](#start-using-czarcoin)
- [License](#license)
- [Support](#support)

# Contributing to Czarcoin

All of our code for Czarcoin v3 is open source. Have a code change you think would make Czarcoin better? Please send a pull request along! Make sure to sign our [Contributor License Agreement (CLA)](https://docs.google.com/forms/d/e/1FAIpQLSdVzD5W8rx-J_jLaPuG31nbOzS8yhNIIu4yHvzonji6NeZ4ig/viewform) first. See our [license section](#license) for more details.

Have comments, bug reports, or suggestions? Want to propose a PR before hand-crafting it? Jump on to our [Rocketchat](https://community.czarcoin.org) and join the [#dev channel](https://community.czarcoin.org/channel/dev) to say hi to the developer community and to talk to the Czarcoin core team.

### Issue tracking and roadmap

See the breakdown of what we're building by checking out the following resources:

 * [White paper](https://czarcoin.org/white-paper)
 * [Aha! Roadmap](https://czarcoin.aha.io/published/01ee405b4bd8d14208c5256d70d73a38)
 * [Jira Issues](https://czarcoin.atlassian.net/projects/V3)

### Install required packages

To get started running Czarcoin locally, download and install the latest release of Go (at least Go 1.11) at [golang.org](https://golang.org).

You will also need [Git](https://git-scm.com/). (`brew install git`, `apt-get install git`, etc).
If you're building on Windows, you also need to install and have [gcc](https://gcc.gnu.org/install/binaries.html) setup correctly.

We support Linux, Mac, and Windows operating systems. Other operating systems supported by Go should also be able to run Czarcoin.

### Download and compile Czarcoin

> **Aside about GOPATH**: Go 1.11 supports a new feature called Go modules,
> and Czarcoin has adopted Go module support. If you've used previous Go versions,
> Go modules no longer require a GOPATH environment variable. Go by default
> falls back to the old behavior if you check out code inside of the directory
> referenced by your GOPATH variable, so make sure to use another directory,
> `unset GOPATH` entirely, or set `GO111MODULE=on` before continuing with these
> instructions.

First, fork our repo and clone your copy of our repository.

```bash
git clone git@github.com:<your-username>/czarcoin czarcoin
cd czarcoin
```

Then, let's install Czarcoin.

```bash
go install -v ./cmd/...
```

### Make changes and test

Make the changes you want to see! Once you're done, you can run all of the unit tests:

```bash
go test -v ./...
```

You can also execute only a single test package if you like. For example:
`go test ./pkg/kademlia`. Add `-v` for more informations about the executed unit
tests.

### Push up a pull request

Use Git to push your changes to your fork:

```bash
git commit -a -m 'my changes!'
git push origin master
```

Use Github to open a pull request!

# Start using Czarcoin

We have another repository dedicated to [documentation and tutorials](https://github.com/czarcoin/docs).
Check out these three tutorials:

 * [Using the Czarcoin Test Network](https://github.com/czarcoin/docs/blob/master/test-network.md)
 * [Using the Uplink CLI](https://github.com/czarcoin/docs/blob/master/uplink-cli.md)
 * [Using the S3 Gateway](https://github.com/czarcoin/docs/blob/master/s3-gateway.md)

# License

The network under construction (this repo) is currently licensed with the
[AGPLv3](https://www.gnu.org/licenses/agpl-3.0.en.html) license. Once the network
reaches beta phase, we will be licensing all client-side code via the
[Apache v2](https://www.apache.org/licenses/LICENSE-2.0) license.

For code released under the AGPLv3, we request that contributors sign our
[Contributor License Agreement (CLA)](https://docs.google.com/forms/d/e/1FAIpQLSdVzD5W8rx-J_jLaPuG31nbOzS8yhNIIu4yHvzonji6NeZ4ig/viewform) so that we can relicense the
code under Apache v2, or other licenses in the future.

# Support

If you have any questions or suggestions please reach out to us on
[Rocketchat](https://community.czarcoin.org/) or
[Twitter](https://twitter.com/czarcoinproject).
