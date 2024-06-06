# vault-kv-restore

Using this CLI tool, you can restore Vault KV v2 Secrets Engine Secrets to a Vault instance! :D

Note: The tool is written in Golang and uses the Vault Official Golang API. The Official Vault Golang API documentation is here - https://pkg.go.dev/github.com/hashicorp/vault/api

Note: The tool needs Vault credentials of a user/account that has access to Vault, to create and update the Vault KV v2 Secrets Engine Secrets that you want to restore. So, in short and basic terms - some sort of "write" permission is required. 

Note: We have tested this only with some versions of Vault (like v1.15.x). So beware to test this in a testing environment with whatever version of Vault you are using, before using this in critical environments like production! Also, ensure that the testing environment is as close to your production environment as possible so that your testing makes sense

Note ⚠️‼️🚨: If the Vault instance has some Secrets already defined at the given Vault KV v2 Secrets Engine mount path, and if some of those Secrets have the same path and name as the Secrets present in the Vault KV Backup JSON file, then, when restoring the Secrets to the Vault instance using the Vault KV Backup JSON file, the Secrets in the Vault instance will be overwritten! All the Vault KV v2 Secrets Engine Secrets in Vault KV Backup JSON file will be present in the Vault instance. If the Vault instance has some extra Vault Secrets configured apart from the ones present in the Vault KV Backup JSON file, it would have those untouched and intact

## Building

```bash
CGO_ENABLED=0 go build -v
```

or

```bash
make
```

## Authorization Details for the Vault Credentials

As mentioned before in a note, the tool needs Vault credentials of a user/account that has access to Vault, to create and update the Vault KV v2 Secrets Engine Secrets that you want to restore.

An example Vault Policy that's required to restore all the secrets in a Vault KV v2 Secrets Engine is -

```hcl
# Vault KV v2 Secrets Engine mount path is "secret"
path "secret/*" {
  capabilities = ["create", "update"]
}
```

You can use a similar Vault Policy based on the mount path of the Vault KV v2 Secrets Engine that you are using and want to restore. You can create a Vault Token that has this Vault Policy attached to it and use that token to restore the Vault KV v2 Secrets Engine Secrets to Vault using the `vault-kv-restore` tool and Vault KV Backup JSON file :)

## Usage

```bash
$ ./vault-kv-restore --help

usage: vault-kv-restore [-quiet|--quiet] [-file|-file <vault-kv-backup-json-file-path>] <kv-mount-path>

Note that the flags MUST come before the arguments

arguments of vault-kv-restore:

  <kv-mount-path> string
    vault kv v2 secrets engine mount path where you want to
    restore the secrets to

flags of vault-kv-restore:

  -file / --file string
      vault kv backup json file path (default "vault_kv_backup.json")

  -quiet / --quiet
      quiet progress (default false).
      By default vault-kv-restore CLI will show a lot of details
      about the restore process and detailed progress during the
      restore process

  -h / -help / --help
      show help

examples:

# show help
vault-kv-restore -h
vault-kv-restore --help

# restores all vault kv v2 secrets from the kv backup JSON file.
# also, any existing vault kv v2 secrets with the same secret
# path in the kv backup JSON file will be overwritten.
vault-kv-restore -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --file too instead of -file

vault-kv-restore --file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# quietly restore all vault kv v2 secrets.
# this will just show dots (.) for progress
vault-kv-restore -quiet -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --quiet too instead of -quiet

vault-kv-restore --quiet --file <path-to-vault-kv-backup-json-file> <kv-mount-path>
```

# Demo

I created a new dummy local Vault instance in developer mode for this demo. I ran the Vault server like this -

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

Initially the Vault looks like this -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault secrets list
Path          Type         Accessor              Description
----          ----         --------              -----------
cubbyhole/    cubbyhole    cubbyhole_dc42bced    per-token private secret storage
identity/     identity     identity_90529884     identity store
secret/       kv           kv_61a7d458           key/value secret storage
sys/          system       system_f3b786af       system endpoints used for control, policy and debugging

$ vault kv list secret
No value found at secret/metadata

$ vault kv list -mount=secret
No value found at secret/metadata
```

Now, let's do a restore of secrets from a Vault KV Backup JSON file. But before we do that, let's create a token which has very narrow set of permissions / specific set of permissions which can just do a restore of secrets to a particular secrets engine. We will make use of Vault Policy here to do the access control / authorization and do just that using our root user with the root token of the Vault :)

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault policy list
default
root

$ vault policy write kv_restore -
# Vault KV v2 Secrets Engine mount path is "secret"
path "secret/*" {
  capabilities = ["create", "update"]
}
Success! Uploaded policy: kv_restore

$ vault policy list
default
kv_restore
root

$ vault policy read kv_restore
# Vault KV v2 Secrets Engine mount path is "secret"
path "secret/*" {
  capabilities = ["create", "update"]
}

$ vault token create -no-default-policy -policy kv_restore
Key                  Value
---                  -----
token                hvs.CAESIC95G38w0e9VR0kRgjujVbNCoj9pr8Sw9zJzCIWX_12kGh4KHGh2cy5sZENGUU9PN2dIV2lLUlEzWmc4ZFZVeE4
token_accessor       LozCNeDnHBUphr89SUuvocTF
token_duration       768h
token_renewable      true
token_policies       ["kv_restore"]
identity_policies    []
policies             ["kv_restore"]
```

Now, let's use this new token to do the restore / restoration of Vault KV v2 Secrets Engine Secrets

I have a sample Vault KV Backup JSON file here -

```bash
$ cat my_another_secret_backup.json 
{"secrets":{"foo":{"bar":"baz","blah":"bloo","blee":"bley"},"something/over/here":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}},"something/over/here-and-there-haha":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}},"something/over/there/ah":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}}}}

$ cat my_another_secret_backup.json | jq
{
  "secrets": {
    "foo": {
      "bar": "baz",
      "blah": "bloo",
      "blee": "bley"
    },
    "something/over/here": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "something/over/here-and-there-haha": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "something/over/there/ah": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    }
  }
}
```

Let's restore these Secrets into a Vault KV v2 Secrets Engine that's already enabled and available and mounted at the mount path `secret/`

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIC95G38w0e9VR0kRgjujVbNCoj9pr8Sw9zJzCIWX_12kGh4KHGh2cy5sZENGUU9PN2dIV2lLUlEzWmc4ZFZVeE4"

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

restoring secrets to `something/over/here` secrets path in vault

restoring secrets to `something/over/here-and-there-haha` secrets path in vault

restoring secrets to `something/over/there/ah` secrets path in vault
```

That's all, it got restored :D Let's see if all is good within Vault :) and if the restoration was actually successful by looking at the Vault KV v2 Secrets Engine Secrets mounted at `secret/` mount path

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"


$ vault kv list
Not enough arguments (expected 1, got 0)

$ vault kv list secret
Keys
----
foo
something/

$ vault kv list secret/foo
No value found at secret/metadata/foo

$ vault kv get secret/foo
= Secret Path =
secret/data/foo

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T10:41:33.563201Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

==== Data ====
Key     Value
---     -----
bar     baz
blah    bloo
blee    bley

$ vault kv list secret
Keys
----
foo
something/

$ vault kv list secret/something
Keys
----
over/

$ vault kv list secret/something/over
Keys
----
here
here-and-there-haha
there/

$ vault kv list secret/something/over/here
No value found at secret/metadata/something/over/here

$ vault kv get secret/something/over/here
========= Secret Path =========
secret/data/something/over/here

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T10:41:33.564285Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get -format json secret/something/over/here
{
  "request_id": "e67e2ff8-3ce8-1d6e-81ae-61d163e2ab44",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "metadata": {
      "created_time": "2024-06-06T10:41:33.564285Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
  "warnings": null
}

$ vault kv list secret/something/over/here-and-there-haha
No value found at secret/metadata/something/over/here-and-there-haha

$ vault kv get secret/something/over/here-and-there-haha
================= Secret Path =================
secret/data/something/over/here-and-there-haha

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T10:41:33.564991Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get -format json secret/something/over/here-and-there-haha
{
  "request_id": "fd83dc09-6011-8850-ba12-f484db241241",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "metadata": {
      "created_time": "2024-06-06T10:41:33.564991Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
  "warnings": null
}

$ vault kv list secret/something/over/there
Keys
----
ah

$ vault kv list secret/something/over/there/ah
No value found at secret/metadata/something/over/there/ah

$ vault kv get secret/something/over/there/ah
=========== Secret Path ===========
secret/data/something/over/there/ah

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T10:41:33.565621Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get -format json secret/something/over/there/ah
{
  "request_id": "9c809e20-0495-f114-ec89-9c864a12a750",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "metadata": {
      "created_time": "2024-06-06T10:41:33.565621Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
  "warnings": null
}
```

As you can see, all is good :) All the secrets have been restored properly :)

# Possible Errors

There are quite some possible errors you can face. Mostly relating to one of the following

- DNS Resolution issues. If you are accessing Vault using it's domain name (DNS record), and not an IP address, then ensure that the DNS resolution works well
- Connectivity issues with Vault. Ensure you have good network connectivity to the Vault system. Ensure the IP you are connecting to is right and belong to the Vault API server, and also check the API server port too.
- Access / Authorization issues. Ensure you have enough access to create and update the Vault KV v2 Secrets Engine that you want to restore to

Example access errors / authorization errors / permission errors -

I'll start off with a brand new Vault instance

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

I'll create a token with half baked access and NOT with full access that's required by the `vault-kv-restore` tool, and use that token with half baked access for the restore. Let's see how that goes ;)

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault policy list
default
root

$ vault policy write half_baked_write_access1 -
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["create"]
}
Success! Uploaded policy: half_baked_write_access1

$ vault policy list
default
half_baked_write_access1
root

$ vault policy read half_baked_write_access1 
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["create"]
}

$ vault token create -no-default-policy -policy half_baked_write_access1
Key                  Value
---                  -----
token                hvs.CAESIL0gzgibdiRNafK-fmqAAgWc2aeenBnSw3OX2Ya71lyjGh4KHGh2cy4wcGtmZmVMWEp3TlU4VjZjNVFWWHp6S1Y
token_accessor       KlQu0ypK2tcTaL9RLpOTgsNs
token_duration       768h
token_renewable      true
token_policies       ["half_baked_write_access1"]
identity_policies    []
policies             ["half_baked_write_access1"]
```

Now, let's use this token to do the restore of Vault KV v2 Secrets Engine Secrets and see if it works :)

I'm going to restore Secrets into a Vault KV v2 Secrets Engine that's already enabled and available and mounted at the mount path `secret/`

And I'm going to restore twice into Vault KV v2 Secrets Engine. I'll tell you why :)

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIL0gzgibdiRNafK-fmqAAgWc2aeenBnSw3OX2Ya71lyjGh4KHGh2cy4wcGtmZmVMWEp3TlU4VjZjNVFWWHp6S1Y"

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

restoring secrets to `something/over/here` secrets path in vault

restoring secrets to `something/over/here-and-there-haha` secrets path in vault

restoring secrets to `something/over/there/ah` secrets path in vault
vault-kv-restore $ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault
error restoring vault policies: error occurred while putting/writing the secrets at path `foo` in vault: error writing secret to secret/data/foo: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/secret/data/foo
Code: 403. Errors:

* 1 error occurred:
	* permission denied
```

As you can see here, when I try to restore the second time, it fails. This is because, the first time it just required `create` access / capability and it had that in it's token's associated Vault Policy in the `half_baked_write_access1` Vault Policy for the given Vault KV v2 Secret Engine's mount path and it's sub-paths `secret/*`

For it to be to be able to run the second time, it needs `update` access / capability, as it's trying to create / write a new version of the secret which already exists. If a secret does NOT already exist, then it just needs `create` access / capability

Now, let's create a different half baked access token with just `update` access and see where that will work and where that won't work

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault policy list
default
half_baked_write_access1
root

$ vault policy write half_baked_write_access2 -
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["update"]
}
Success! Uploaded policy: half_baked_write_access2

$ vault policy list
default
half_baked_write_access1
half_baked_write_access2
root

$ vault policy read half_baked_write_access2
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["update"]
}

$ vault token create -no-default-policy -policy half_baked_write_access2
Key                  Value
---                  -----
token                hvs.CAESIDQDEM89wIneElzbFOUU2LW4ESlziQLWgeSC9zbd0M0CGh4KHGh2cy5OMDVSOWhhNjhXR3RKR0VrMTR6cXpBNkc
token_accessor       6qzhOtMUYEXJacjPQYYQuyag
token_duration       768h
token_renewable      true
token_policies       ["half_baked_write_access2"]
identity_policies    []
policies             ["half_baked_write_access2"]
```

And this Vault already has some secrets

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault kv get secret/foo
= Secret Path =
secret/data/foo

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T11:17:30.401376Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            3

==== Data ====
Key     Value
---     -----
bar     baz
blah    bloo
blee    bley

$ vault kv get secret/something/over/here
========= Secret Path =========
secret/data/something/over/here

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T11:17:30.402023Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            3

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get secret/something/over/here-and-there-haha
================= Secret Path =================
secret/data/something/over/here-and-there-haha

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T11:17:30.399699Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            3

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get secret/something/over/there/ah
=========== Secret Path ===========
secret/data/something/over/there/ah

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-06T11:17:30.400731Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            3

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]
```

Let's use this to do the restore of Vault KV v2 Secrets Engine Secrets mounted at `secret/` mount path. Let's use the 

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIDQDEM89wIneElzbFOUU2LW4ESlziQLWgeSC9zbd0M0CGh4KHGh2cy5OMDVSOWhhNjhXR3RKR0VrMTR6cXpBNkc"

$ cat my_another_secret_backup.json 
{"secrets":{"foo":{"bar":"baz","blah":"bloo","blee":"bley"},"something/over/here":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}},"something/over/here-and-there-haha":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}},"something/over/there/ah":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}}}}

$ cat my_another_secret_backup.json | jq
{
  "secrets": {
    "foo": {
      "bar": "baz",
      "blah": "bloo",
      "blee": "bley"
    },
    "something/over/here": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "something/over/here-and-there-haha": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "something/over/there/ah": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    }
  }
}

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

restoring secrets to `something/over/here` secrets path in vault

restoring secrets to `something/over/here-and-there-haha` secrets path in vault

restoring secrets to `something/over/there/ah` secrets path in vault

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

restoring secrets to `something/over/here` secrets path in vault

restoring secrets to `something/over/here-and-there-haha` secrets path in vault

restoring secrets to `something/over/there/ah` secrets path in vault

$ ./vault-kv-restore --file my_another_secret_backup.json secret 

restoring secrets to `foo` secrets path in vault

restoring secrets to `something/over/here` secrets path in vault

restoring secrets to `something/over/here-and-there-haha` secrets path in vault

restoring secrets to `something/over/there/ah` secrets path in vault
```

As you can see, it's able to restore the Vault KV v2 Secrets Engine Secrets multiple times!

Now, let's create a brand new Vault instance and try to do the same :)

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault policy list
default
root

$ vault policy write half_baked_write_access2 -
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["update"]
}
Success! Uploaded policy: half_baked_write_access2

$ vault policy read half_baked_write_access2
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["update"]
}

$ vault policy list
default
half_baked_write_access2
root

$ vault secrets list
Path          Type         Accessor              Description
----          ----         --------              -----------
cubbyhole/    cubbyhole    cubbyhole_e3af728f    per-token private secret storage
identity/     identity     identity_3ba5fa9e     identity store
secret/       kv           kv_79c4a2ef           key/value secret storage
sys/          system       system_8801de70       system endpoints used for control, policy and debugging

$ vault kv list
Not enough arguments (expected 1, got 0)

$ vault kv list secret
No value found at secret/metadata

$ vault token create -no-default-policy -policy half_baked_write_access2
Key                  Value
---                  -----
token                hvs.CAESIKNuHOy6q7dRF0KRmXttsLYWXdjntmMuVNG3ELgzbr_2Gh4KHGh2cy5CUUtwbU1va1hhT1QyRFdPd09rZXRua1o
token_accessor       UEyQkQ7ZL7naGIutOF890JbA
token_duration       768h
token_renewable      true
token_policies       ["half_baked_write_access2"]
identity_policies    []
policies             ["half_baked_write_access2"]
```

As you can see, this Vault instance does NOT have any existing secrets, let's restore to this brand new Vault KV v2 Secrets Engine using our existing Vault KV Backup JSON File

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIKNuHOy6q7dRF0KRmXttsLYWXdjntmMuVNG3ELgzbr_2Gh4KHGh2cy5CUUtwbU1va1hhT1QyRFdPd09rZXRua1o"

$ ./vault-kv-restore --file my_another_secret_backup.json secret

restoring secrets to `foo` secrets path in vault
error restoring vault policies: error occurred while putting/writing the secrets at path `foo` in vault: error writing secret to secret/data/foo: Error making API request.

URL: PUT http://127.0.0.1:8200/v1/secret/data/foo
Code: 403. Errors:

* 1 error occurred:
	* permission denied
```

As you can see, this token cannot `create` secrets in the Vault KV v2 Secrets Engine mounted at `secret/` mount path. This is because of the lack of access / capability associated with it's token. It only has access to `update` secrets in the Vault KV v2 Secrets Engine mounted at `secret/` mount path, but not `create` secrets in the Vault KV v2 Secrets Engine mounted at `secret/` mount path

Also, if there's not enough permissions in other ways also it won't work - like, access to `create` and/ `update` some secrets but not others, all within the same Vault KV v2 Secrets Engine mounted at a particular path. Also, if there's some explicit `deny` on some secret / secrets, or some secret path / secrets path / paths, then that will also cause problems for the tool. As of now, the tool abruptly stops when there's such an error

# Future Ideas

- Any and all issues in the [GitHub Issues](https://github.com/karuppiah7890/vault-kv-restore/issues) section

- Allow user to say "It's okay if the tool cannot restore some secrets and/ some secret paths, due to permission issues. Just restore the secrets the tool can" and be able to skip intermittent errors here and there and ignore the errors rather than abruptly stop at errors like it does not

- Support restoring of multiple specific KV v2 secrets engine secrets in a single restore at once by providing a file which contains the mount paths of the secrets engines to be restore, or by providing the mount paths of the secrets engines as arguments to the CLI, or provide the ability to use either of the two or even both. Either restore same set of secrets to different mount paths. Or different sets of secrets for different mount paths. Both should be available I think
  - This would mean we can also - Support restoring of all the secrets in all the secrets engines in a single restore

- Support restoring of the KV v2 secrets engine configuration too apart from the secrets in the secrets engine. https://developer.hashicorp.com/vault/api-docs/secret/kv/kv-v2#configure-the-kv-engine

- Support restoring the metadata of the secrets apart from the secrets of the secrets engine. This can be tricky. This has be thought through. But basic feature should be straight forward - restore metadata for all or none - or restore metadata for anything for which the metadata is present in the Vault KV Backup JSON file

- Support restoring all the versions of the secrets in the secrets engine. This is a bit tricky - we need to know which version to restore, or we could restore all versions to start off with and put the latest version as the latest

- Support restoring secrets with a combination of the above features, that is
  - multiple/all secrets
  - one/multiple/all secrets engines
  - only latest/all versions of secrets
  - metadata of secrets
  - configuration of secrets engines

- Allow for flags to come even after the arguments, not just before arguments. This would require using better CLI libraries / frameworks like [`spf13/pflag`](https://github.com/spf13/pflag), [`spf13/cobra`](https://github.com/spf13/cobra), [`urfave/cli`](https://github.com/urfave/cli) etc than using the basic Golang's built-in `flag` standard library / package

# Contributing

Please look at https://github.com/karuppiah7890/vault-tooling-contributions for some basic details on how you can contribute

