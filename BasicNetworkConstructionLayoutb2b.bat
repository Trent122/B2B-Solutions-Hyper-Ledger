#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
version: '2'

networks:
  basic:

services:
  B2bSolutions.com:
    image: hyperledger/fabric-ca:$TAG
    environment:
      - FABRIC_CA_HOME=/etc/hyperledger/fabric-ca-server
      - FABRIC_CA_SERVER_CA_NAME=ca.example.com
      - FABRIC_CA_SERVER_CA_CERTFILE=/etc/hyperledger/fabric-ca-server-config/ca.org1.example.com-cert.pem
      - FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/689329f7c94d77a1e9b3581a509c108f2160aa964b62dd2f1c9c1fdda77bc258_sk
    ports:
      - "7065:7065"
    command: sh -c 'fabric-ca-server start -b admin:adminpw'
    volumes:
      - ./config/crypto-config/peerOrganizations/org1.example.com/ca/:/etc/hyperledger/fabric-ca-server-config
    container_name: ca.example.com
    networks:
      - basic

  orderer.example.com:
    container_name: "B2B Solutions"
    image: hyperledger/fabric-orderer:$TAG
    environment:
      - ORDERER_GENERAL_LOGLEVEL=debug
      - ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
      - ORDERER_GENERAL_GENESISMETHOD=file
      - ORDERER_GENERAL_GENESISFILE=/etc/hyperledger/configtx/genesis.block
      - ORDERER_GENERAL_LOCALMSPID=OrdererMSP
      - ORDERER_GENERAL_LOCALMSPDIR=/etc/hyperledger/msp/orderer/msp
      - ORDERER_GENERAL_GENESISPROFILE=SampleSingleMSPSolo
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/orderer
    command: orderer
    ports:
      - 7050:7050
    volumes:
        - ./artifacts/:/etc/hyperledger/configtx
        - ./config/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/:/etc/hyperledger/msp/orderer
        - ./config/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/:/etc/hyperledger/msp/peerOrg1
    networks:
      - basic

  peer0.org1.B2BSolutions.com:
    container_name: peer0.org1.B2BSolutions.com
    image: hyperledger/fabric-peer:$TAG
    environment:
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_PEER_ID=peer0.org1.example.com
      - CORE_LOGGING_PEER=debug
      - CORE_CHAINCODE_LOGGING_LEVEL=debug
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/peer/
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      # # the following setting starts chaincode containers on the same
      # # bridge network as the peers
      # # https://docs.docker.com/compose/networking/
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_basic
      - CORE_LEDGER_STATE_STATEDATABASE=LevelDB 
      # - CORE_LEDGER_STATE_COUCHDBCONFIG_COUCHDBADDRESS=couchdb:5984
      # The CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME and CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD
      # provide the credentials for ledger to connect to CouchDB.  The username and password must
      # match the username and password set for the associated CouchDB.
      # - CORE_LEDGER_STATE_COUCHDBCONFIG_USERNAME=
      # - CORE_LEDGER_STATE_COUCHDBCONFIG_PASSWORD=
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric
    command: peer node start --peer-chaincodedev=true
    # command: peer node start --peer-chaincodedev=true
    ports:
      - 7001:7001
      - 7982:7982
    volumes:
        - /var/run/:/host/var/run/
        - ./config/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/msp/peer
        - ./config/crypto-config/peerOrganizations/org1.example.com/users:/etc/hyperledger/msp/users
        - ./artifacts:/etc/hyperledger/configtx
    depends_on:
      - orderer.example.com
      # - couchdb
    networks:
      - basic

  # couchdb:
  #   container_name: couchdb
  #   image: hyperledger/fabric-couchdb
  #   # Populate the COUCHDB_USER and COUCHDB_PASSWORD to set an admin user and password
  #   # for CouchDB.  This will prevent CouchDB from operating in an "Admin Party" mode.
  #   environment:
  #     - COUCHDB_USER=
  #     - COUCHDB_PASSWORD=
  #   ports:
  #     - 5984:5984
  #   networks:
  #     - basic

  cli:
    container_name: cli
    image: hyperledger/fabric-tools:$TAG
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=debug
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
      - CORE_CHAINCODE_KEEPALIVE=10
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: /bin/bash
    volumes:
        - /var/run/:/host/var/run/
        - ./chaincode/:/opt/gopath/src/github.com/
        - ./config/crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
    networks:
        - basic
    #depends_on:
    #  - orderer.example.com
    #  - peer0.org1.example.com
    #  - couchdb

  chaincode:
    container_name: chaincode
    image: hyperledger/fabric-ccenv:$TAG
    tty: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      - CORE_LOGGING_LEVEL=debug
      - CORE_PEER_ID=chaincode
      - CORE_PEER_ADDRESS=peer0.org1.example.com:7051
      - CORE_PEER_LOCALMSPID=Org1MSP
      - CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_basic
      - CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp
    working_dir: /opt/gopath/src/chaincode
    command: /bin/bash -c 'sleep 6000000'
    volumes:
        - /var/run/:/host/var/run/
        - ./config/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/msp:/etc/hyperledger/msp
        - ./chaincode:/opt/gopath/src/chaincode
    networks:
      - basic
    depends_on:
      - cli

  explorer:
    container_name: explorer
    image: zeepeetek/blockchain-explorer:3.1
    environment:
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=postgres
      - PGPASSWORD=postgres
      - LOG_LEVEL=debug
    volumes:
      - ./config/crypto-config/peerOrganizations/org1.example.com:/etc/hyperledger/msp
      - ./config/config.json:/var/blockchain-explorer/config.json
    ports:
      - 3000:3000
    depends_on:
      - postgres
    networks:
      - basic

  postgres:
    container_name: postgres
    image: postgres:10.3-alpine
    ports:
      - 5432:5432
    environment:
      - POSTGRES_PASSWORD=postgres
    networks:
      - basic