# SMCCLI

SMCCLI is a CLI tool inspired by the dgkcli from DELA. Here is a complete scenario:

```sh
# Install the CLI
go install .

# Run 3 nodes. Do that in 3 different sessions
LLVL=info smccli --config /tmp/node1 start --routing tree --listen tcp://127.0.0.1:2001
LLVL=info smccli --config /tmp/node2 start --routing tree --listen tcp://127.0.0.1:2002
LLVL=info smccli --config /tmp/node3 start --routing tree --listen tcp://127.0.0.1:2003

# Exchange certificates
smccli --config /tmp/node2 minogrpc join --address //127.0.0.1:2001 $(smccli --config /tmp/node1 minogrpc token)
smccli --config /tmp/node3 minogrpc join --address //127.0.0.1:2001 $(smccli --config /tmp/node1 minogrpc token)

# Initialize DKG on each node. Do that in a 4th session.
smccli --config /tmp/node1 dkg listen
smccli --config /tmp/node2 dkg listen
smccli --config /tmp/node3 dkg listen

# Do the setup in one of the node:
smccli --config /tmp/node1 dkg setup \
    --authority $(cat /tmp/node1/dkgauthority) \
    --authority $(cat /tmp/node2/dkgauthority) \
    --authority $(cat /tmp/node3/dkgauthority)

# Encrypt a message:
smccli --config /tmp/node2 dkg encrypt --message deadbeef

# Decrypt a message
smccli --config /tmp/node3 dkg decrypt --encrypted <...>

# Send transaction to blockchain 
smccli --config /tmp/node1 calypso advertise --blockchain tcp://127.0.0.1:10001