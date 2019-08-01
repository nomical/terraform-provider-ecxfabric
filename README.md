Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)


Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/nomical/terraform-provider-ecxfabric`

```sh
$ mkdir -p $GOPATH/src/github.com/nomical; cd $GOPATH/src/github.com/nomical
$ git clone git@github.com:nomical/terraform-provider-ecxfabric
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/nomical/terraform-provider-ecxfabric
$ make build
```

Using the provider
------------------

If you want to run Terraform with the ECXFabric provider plugin on your system, complete the following steps:

1. [Download and install Terraform for your system](https://www.terraform.io/intro/getting-started/install.html). 

2. [Download the ECXFabric provider plugin for Terraform](https://github.com/nomical/terraform-provider-ecxfabric/releases).

3. Unzip the release archive to extract the plugin binary (`terraform-provider-ecxfabric_vX.Y.Z`).

4. Move the binary into the Terraform [plugins directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for the platform.
    - Linux/Unix/OS X: `~/.terraform.d/plugins`
    - Windows: `%APPDATA%\terraform.d\plugins`

5. Ensure you have an Equinix developer account. Visit [Equinix Developer site](https://developer.equinix.com).
6. Have your Oauth2 details to hand. See https://developer.equinix.com/user/me/apps to obtain your details, create an App if one doesn't exist.
7. You can either use environment variables to configure the provider or explicitly in the provider configuration.

```sh
export ECXFABRIC_CLIENT_ID="Equinix App OAuth2 Client ID"
export ECXFABRIC_CLIENT_SECRET="Equinix OAuth2 Client Secret"
export ECXFABRIC_USERNAME="Equinix Username"
export ECXFABRIC_PASSWORD="Equinix Password"
```

6. Add the plug-in provider to the Terraform configuration file.

```
provider "ecxfabric" {
  client_id     = "Equinix App OAuth2 Client ID"
  client_secret = "Equinix OAuth2 Client Secret"
  username      = "Equinix Username"
  password      = "Equinix Password"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-ecxfabric
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
In order to run a particular Acceptance test, export the variable `TESTARGS`. For example

```sh
$ export TESTARGS="-run TestAccL2Connection_Basic"
$ make testacc
```

Shorter version

```sh
$ make testacc TESTARGS="-run=TestAccL2Connection_Basic"
```

Issuing `make testacc` will now run the testcase with names matching `TestAccL2Connection_Basic`. This particular testcase is present in
`ecxfabric/resource_l2_connection_test.go`

You will also need to export the following environment variables for running the Acceptance tests.
* `ECXFABRIC_CLIENT_ID`
* `ECXFABRIC_CLIENT_SECRET`
* `ECXFABRIC_USERNAME`
* `ECXFABRIC_PASSWORD`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_REGION`

Additional environment variables may be required depending on the tests being run. Check console log for warning messages about required variables. 