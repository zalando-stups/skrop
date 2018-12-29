## Continuous Integration

For CI are using [_Travis CI_](https://travis-ci.org/zalando-stups/skrop).
Every time a PR is merges to master, _Travis_ will create a new release as well.
Both a Github release and a docker image will be published.

### GitHub Token

In order for _Travis CI_ to perform its duty, a valid _GitHub_ token must be encrypted in
the [Travis configuration file](./.travis.yml).

If this token needs to change, here is how to set it up:
- Log on to _GitHub_ on the account you want _Travis CI_ to impersonate when performing operation on the repository
- Go to _GitHub_, in your [Personal access tokens](https://github.com/settings/tokens) configuration
- Click on _Generate new token_
- Check the `public_repo` scope
- Click _Generate token_

We will now encrypt that token with the [_travis CLI_](https://docs.travis-ci.com/user/encryption-keys/). Make sure it
is installed.

```bash
travis encrypt GITHUB_AUTH=the_token_copied_from_github
```

This will encrypt it in a way that **only** _Travis CI_ can decrypt. Your personal token is then quite safe.

The output should look like this.

```
  secure: "someBASE64value"
```

Take the output (the complete YAML key as it appears) and replace, in the [`.travis.yaml` file](./.travis.yml), this
line with the output of the previous command.

```yaml
env:
- secure: "someBASE64value"
```
