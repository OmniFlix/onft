# oNFT

The `oNFT` module is a part of the OmniFlix Network and enables the classification and tokenization of assets.

- Assets can be classified under `denoms` (aka `collections` across various ecosystems)
- Tokenize media assets by minting NFTs

The module supports the following capabilities for classification and tokenization:

- Creation of collections (denoms)
- Minting of NFTs
- Transferring of NFTs
- Burning of NFTs

Various queries are available to get details about denoms/collections, NFTs, and related metadata including but not limited to ownership. Click here to try them out by interacting with the chain.

The module utilizes the [irismod/nft](https://github.com/irismod/nft) repository and has been modified to meet the requirements of the OmniFlix Network. It can be used through the CLI with various commands and flags to perform the desired actions.

## 1) Create Denom (Collection)
To create an oNFT denom, you will need to use the "onftd tx onft create" command with the following args and flags:

args:
symbol: denom symbol

flags:
name : name of denom/collection 
description: description for the denom
preview-uri: display picture url for denom
schema: json schema for additional properties
creation-fee: denom creation-fee to create denom

Example:
```
    onftd tx onft create <symbol> \
     --name=<name> \
     --description=<description> \
     --preview-uri=<preview-uri> \
     --schema=<schema> \
     --creation-fee=<creation-fee> \
     --chain-id=<chain-id> \
     --fees=<fee> \
     --from=<key-name>
```
  
## 2) Mint an oNFT

To create an oNFT, you will need to use the "onftd tx onft mint" command with the following flags:

denom-id: the ID of the collection in which you want to mint the NFT
name: the name of the NFT
description: a description of the NFT
media-uri: the IPFS URI of the NFT
preview-uri: the preview URI of the NFT
data: any additional properties for the NFT (optional)
recipient: the recipient of the NFT (optional, default is the minter of the NFT)
non-transferable: flag to mint a non-transferable NFT (optional, default is false)
inextensible: flag to mint an inextensible NFT (optional, default is false)
nsfw: flag to mark the NFT as not safe for work (optional, default is false)
royalty-share: the royalty share for the NFT (optional, default is 0.00)

Example:

```
onftd tx onft mint <denom-id>
--name="NFT name"
--description="NFT description"
--media-uri="https://ipfs.io/ipfs/...."
--preview-uri="https://ipfs.io/ipfs/...."
--data=""
--recipient=""
--non-transferable
--inextensible
--nsfw
--chain-id=<chain-id>
--fees=<fee>
--from=<key-name>
```

For a royalty share of 5%:

```
--royalty-share="0.05" # 5%
```

## 3) Transfer an oNFT

To transfer an oNFT, you will need to use the "onftd tx onft transfer" command with the following flags:

recipient: the recipient's account address
denom-id: the ID of the collection in which the NFT is located
onft-id: the ID of the NFT to be transferred
chain-id: the ID of the blockchain where the transaction will be made (required)
fees: the transaction fees (required)
from: the name of the key to sign the transaction with (required)

Example:

```
onftd tx onft transfer <recipient> <denom-id> <onft-id>
--chain-id=<chain-id>
--fees=<fee>
--from=<key-name>
```

## 4) Burn an oNFT

To burn an oNFT, you will need to use the "onftd tx onft burn" command with the following flags:

denom-id: the ID of the collection in which the NFT is located
onft-id: the ID of the NFT to be burned
chain-id: the ID of the blockchain where the transaction will be made (required)
fees: the transaction fees (required)
from: the name of the key to sign the transaction with (required)

Example:

```
onftd tx onft burn <denom-id> <onft-id>
--chain-id=<chain-id>
--fees=<fee>
--from=<key-name>
```

# All CLI Commands

### Queries
  - #### Get List of denoms (collections)
    ```bash
    onftd query onft denoms
    ```
  - #### Get Denom details by it's Id
     ```bash
    onftd query onft denom <denom-id>
    ```    
  - #### Get List of NFTs in a collection
    ```bash
    onftd query onft collection <denom-id>
    ```
  - #### Get Total Count of NFTs in a collection
    ```bash
    onftd query onft supply <denom-id>
    ```
  - #### Get NFT details by it's Id
    ```bash
    onftd query onft asset <denom-id> <nft-id>
    ```
  - #### Get All NFTs owned by an address
    ```bash
    onftd query onft owner <account-address>
    ```
