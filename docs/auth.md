# Authentication

The first step is to get logged in. You will need an account on an instance of
the Conch API. Please work with an admin and get all squared away. Make sure you
can log into the web UI before continuing here.

## Using Password Auth

A new profile is be created via:
```
conch profile new --name <profile name> --user user@example.com
```

You can provide your password via the `--password` argument.  Otherwise, you
will be prompted for your password. 

For Joyent employees, this will attempt to log you into the production instance.
If you'd like to use staging, provide the `--environment=staging` argument. For
development or non-Joyent users, provide `--environment=development
--url=https://joyent.example.com`, adjusting the URL as appropriate for your
environment.

For instance:
```
conch profile new --name <profile name> --user user@example.com --environment=staging
```
or 

```
conch profile new --name <profile name> --user user@example.com --environment=development --url http://localhost:5000
```

### Expiration

Password authentication credentials expire rather rapidly, on the order of 30
days. If you do not use the shell for 30 days, you may receive an "unauthorized"
error message. To log back in, simply execute

```
conch profile relogin
```

You will be prompted for your password.

### Forced Password Change

Sometimes, you might be forced to change your password, usually because you
requested a password reset or an admin reset your password. Typically, you will
receive an "unauthorized" message. You should try to `relogin` and will then be
prompted to change your password.


## Using API Tokens

### What Are API Tokens?

An API token is a special piece of information that allows automation or other
tools to use your Conch account without needing your actual user name or
password.

### Why Should I Use Them?

Probably the best reason to use API tokens is that in the future, tokens will be
the only way to use the shell. At some point, password authentication will go
away. So why not get ahead of the curve?

Tokens also provide a really nice way to control access to your account. For
instance, you could create a token for each device you have. If that device is
lost, you can revoke only that token, cutting off access for the lost device but
not impacting any others.  

Further, API tokens generally don't expire for a long time. They're designed for
use by automation and it'd be very silly to make someone log back in manually
every now and then on every server. Once you've got an API token, you can
mostly put authentication out of your mind. 

Also, as we'll see shortly, it's possible to use conch with no on-disk
configuration. That option is only available when using tokens.

For now, we'll stick to the basics of API tokens. A lot more documentation can
be found over [here](tokens)

*Just as a note, tokens are obfuscated inside the config file. It is not
possible to use the value recorded there in any other tool.*

### Upgrading An Existing Profile

Once you've created a profile using a password, you can easily upgrade that
profile to token auth.

```
conch profile upgrade
```

This command generates a new token specifically for your local device and
switches the profile to use it. No further action is required on your part. 


### Configure A New Profile

If you already have a token (either from the web UI or `conch user tokens new`,
you can create a profile from scratch using:

```
conch profile new --name <profile name> --token <token>
```

Other modifiers like `--environment` are accepted as well.

### Working Without A Profile

If you have a token, it is possible to use the shell without ever creating a
profile. There are two ways to work without a profile.

First, provide the token via the environment.

```
export CONCH_TOKEN='1234....'
conch profile ls
```

Second, provide the token on the command line.

```
conch --token='1234....' profile ls
```

## Multiple Profiles

It is possible to have as many profiles as you want. Developers typically have
one for production, one for staging, and one for their local development
environment. Each profile must have a unique name but they are created just as
above. 

You can switch between profiles by running:

```
conch profile set active <profile name>
```

or you can override the active profile by using:

```
conch -p <profile name>
```
