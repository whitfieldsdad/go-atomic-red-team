# Generating Atomic Red Team content archives

## Generating tarball files

To compress the `atomics` folder as a tarball:

```bash
tar -zcf atomics.tar.gz atomic-red-team/atomics
```

End-to-end example:

```bash
git clone https://github.com/redcanaryco/atomic-red-team --depth=1
tar -cf atomics.tar.gz atomic-red-team/atomics
rm -rf atomic-red-team
```

## Generating encrypted tarball files

Compress a directory as a GZIP compressed tar archive (i.e. as a "tarball"):

```bash
tar -zcf atomics.tar.gz atomic-red-team/atomics
```

### Symmetric encryption

#### Using `age`

Encrypt:

```bash
age -p atomics.tar.gz > atomics.tar.gz.age
Enter passphrase (leave empty to autogenerate a secure one): ********
```

Decrypt:

```bash
age -d -i atomics.tar.gz.age > atomics.tar.gz
```
