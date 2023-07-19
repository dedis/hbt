#!/usr/bin/env bash

# This script tests the encrypt, reencrypt, reveal functionality

# assert that the DKG key is passed to script
if [[ -z "$1" ]]; then
   echo -e "Please provide the DKG public key as parameter"
   exit 1
fi

smccli smc createkeys

CIPHER=$(smccli --config /tmp/smc1 dkg encrypt --message "deadbeef")
echo -e "Ciphertext \t${CIPHER}\n"

DECIPHERED=$(smccli --config /tmp/smc1 dkg decrypt --encrypted ${CIPHER})
echo -e "Message \t${DECIPHERED}\n"


PUBK=$(cat key.pair | cut -d ":" -f2)
echo -e "User pubk: \t${PUBK}\n"
PRIVK=$(cat key.pair | cut -d ":" -f1)
echo -e "User privk: \t${PRIVK}\n"

XHATENC=$(smccli --config /tmp/smc1 dkg reencrypt --encrypted ${CIPHER} --pubk ${PUBK})
echo -e "XhatEnc: \t${XHATENC}\n"

echo -e "DKG key: \t$1\n"

smccli smc reveal --xhatenc ${XHATENC} --encrypted ${CIPHER} --dkgpub $1 --privk ${PRIVK}
