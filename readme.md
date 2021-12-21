# oNFT

**oNFT** - OmniFlix NFT Module

The code for this module is initially taken from the [irismod/nft](https://github.com/irismod/nft) repository and modified according the structure of oNFT.


## CLI Commands

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
    
### Transactions
  - #### Create Denom / Collection
    Usage
    ```bash
    onftd tx onft create [symbol] [flags] 
    ```
    
    Flags:
      - **name** : name of denom/collection
      - **description**: description for the denom
      - **preview-uri**: display picture url for denom
      - **schema**: json schema for additional properties
      
    Example:
    ```bash
    onftd tx onft create <symbol>  
     --name=<name>
     --description=<description>
     --preview-uri=<preview-uri>
     --schema=<schema>
     --chain-id=<chain-id>
     --fees=<fee>
     --from=<key-name>
    ```
  - #### Mint NFT
    Usage
    ```bash
    onftd tx onft mint [denom-id] [flags]
    ```
    
    Flags:
      - **name** : name of denom/collection (string)
      - **description**: description of the denom (string)
      - **media-uri**: ipfs uri of the nft (url)
      - **preview-uri**: preview uri of the nft (url)
      - **data**: additional nft properties (json string)
      - **recipient**: recipient of the nft (optional, default: minter of the nft)
      
    Example:
    ```bash
    onftd  tx onft mint <denom-id>
     --name=<name>
     --description=<description>
     --media-uri=<preview-uri>
     --preview-uri=<preview-uri>
     --data=<additional nft data json string>
     --recipient=<recipient-account-address>
     --chain-id=<chain-id>
     --fees=<fee>
     --from=<key-name>
      ```
  - #### Transfer NFT
    Usage
    ```bash
    onftd tx onft transfer [recipient] [denom-id] [onft-id] [flags]
    ```
    
    Example:
    ```bash
    onftd  tx onft transfer <recipient> <denom-id> <nft-id>
     --chain-id=<chain-id>
     --fees=<fee>
     --from=<key-name>
    ```

  - #### Burn NFT
    Usage
    ```bash
    onftd tx onft burn [denom-id] [onft-id] [flags]
    ```
    
    Example:
    ```bash
    onftd  tx onft burn <denom-id> <nft-id>
     --chain-id=<chain-id>
     --fees=<fee>
     --from=<key-name>
    ```