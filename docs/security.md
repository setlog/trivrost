# Security of trivrost
To secure trivrost, you should only use it via https. However, to increase the security the deployment-config and the bundle info files have to be digitally signed. This way, every download is cryptographically secured. Furthermore, the deployment-config and the bundle info files must contain a timestamp to further strengthen the security.

# Timestamps
To prevent that an attacker can make the client to install an old (and potentially vulnerable) version of a bundle or trivrost itself, the deployment-config and the bundle info files must contain timestamps. The bundle info files must also contain the `UniqueBundleName`, so the timestamps can be assigned to the correct bundle. Thus you have to make sure that the `UniqueBundleName` really is unique.

If trivrost runs for the first time, it will accept any timestamp. After that, trivrost will only accept timestamps that are not older than the last accepted timestamp for the deployment-config or each of the bundle info files. To do so, trivrost will save the last accepted timestamp in a json file called `timestamps.json`. See [Where does trivrost save files?](file_locations.md) to find out where that file is saved.

The launcher will not check the timestamp against the clock of the client or any time server. It will accept a timestamp, if it is the same as the last seen timestamp, even if the deployment-config changed. It is the responsibility of the creator of the deployment-config, to update the timestamp to enhance the security, even if staying at some fixed arbitrary value would not stop trivrost from working.

The timestamp must have the form of `2006-01-02 15:04:05`. You can generate it on Unix systems using the command `date +"%Y-%m-%d %H:%M:%S"`. We strongly advise to always use UTC.

To set the timestamp automatically, the shell script [insert_timestamp](../scripts/insert_timestamp) is provided by this project. It will substitute a given string (e.g. `<TIMESTMAP>`) with a correctly formed timestamp in UTC. You should call this script in your CI/CD-pipeline before signing the deployment-config.

Please note that you have to set the timestamp in the deployment-config before signing the deployment-config. For the bundle info files, the hasher will automatically set the timestamp in UTC.

# Signing
To sign the deployment-config and bundle info files we use `RSA` with the padding algorithm `PSS`. We use `sha256` as the hashing algorithm for signing. The signatures of the deployment-config have to be stored `base64` encoded. The signatures are saved in separate files with the same url as the original files, but with a `.signature` extension. So the signature for the bundle info file `https://example.com/linux/launcher/bundleinfo.json` has the url `https://example.com/linux/launcher/bundleinfo.json.signature.`

The public keys used to validate the signatures are compiled into the trivrost binary. Therefore you have to create the file `cmd/launcher/resources/public-rsa-keys.pem`. It contains all public keys in the `PEM` format (see [example file](../examples/public-rsa-keys.pem.example)). The public keys are separated by additional line breaks. trivrost checks if a signature is valid against at least one of the public keys. Note that signed resources accessed using the `file://`-scheme are not validated. It should only be used for testing.

# Sign with openssl
You can create the keys and sign using `openssl`. To generate a private key with a size of 4096 bit you can use one of the following two commands:
```
openssl genrsa -out private_key.pem 4096
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:4096
```

To extract the public key you can use the following command:
```
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

Copy all the public keys you need in the `resources/public-rsa-keys.pem` file and never ever share the private keys.

To sign a file called `config.json` and `base64`-encode it, you can use the following two commands:
```
openssl dgst -sha256 -sigopt rsa_padding_mode:pss -sign private_key.pem -out /tmp/sign.sha256 config.json
openssl base64 -in /tmp/sign.sha256 -out config.json.signature
```
To sign the deployment-config and bundle info files, you can use the signer utility at `out/signer`. (Build with `make tools`)

# Verify signature with openssl
If you want to check a given signature by hand, you first have to decode the base64 encoded signature file:
```
openssl base64 -d -in config.json.signature -out config.json.signature.decoded
```
Now you can verify the decoded signature file:
```
openssl dgst -verify public_key.pem -sha256 -sigopt rsa_padding_mode:pss -signature config.json.signature.decoded config.json
```
