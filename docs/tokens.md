# API Tokens

## Environment Variables

* `CONCH_TOKEN` 
  : specify the API token
* `CONCH_ENV`
  : `production`, `staging`, or `development` (defaults to `production`)
* `CONCH_URL` 
  : if `CONCH_ENV=development`, specifies the API's URL

## User Commands

* `profile create --name :name --token :token`
  : create a profile using a token. user names and other parameters are ignored
* `profile set token :token`
  : converts a profile to token auth, using the provided token
* `profile upgrade`
  : converts a profile to token auth by auto-generating a token which is never
  shared to the user
* `profile revoke-tokens --tokens-only`
  : when revoking one's access abilities, one can revoke only API tokens instead
  of both API tokens and logins
* `profile change-password --purge-tokens`
  : when changing one's paswords, one can also purge all API tokens at the same
  time
* `user tokens` 
  : list all API tokens
* `user token create :name` 
  : create a token with the given name, displaying the token value
* `user token get :name` 
  : get basic data (including last used time) about a token. Does *not* show the
  token value
* `user token rm :name`
  : remove a token by name


## Admin Commands
* `admin user :id tokens`
  : List a user's API tokens
* `admin user :id token get :name`
  : Get a user's API token by name. Does *not* show the token value
* `admin user :id token rm :name`
  : Remove a user's token by name
* `admin user :id revoke --tokens-only`
  : when revoking a user's access, an admin can revoke just their API tokens
* `admin user :id reset --revoke-tokens`
  : when reseting a user's passwords, an admin can also revoke their API tokens

## Notes

* The API token is obfuscated in the config file. It is not possible to copy
  that value out of the config and use it in another tool.

## Build / Compilation Flags

* `DISABLE_API_TOKEN_CRUD`
  : removes all commands for creating or modifying API tokens
* `TOKEN_OBFUSCATION_KEY` 
  : used in the obfuscation of tokens in the config file. While a default is
  provided, it is *strongly* recommended that this be customized in the build
  environment. *NOTE*: If this value is ever changed, it will render the tokens
  in user configs completely unusable.

