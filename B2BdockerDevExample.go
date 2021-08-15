Docker version 18.03.1-ce, build 9ee9f40

OrdererOrgs:
  - Name: OrderingService
    Domain: SmoothTech-network.com
    Specs:
      - Hostname: orderer

PeerOrgs:
  - Name: Nvidia
    Domain: Nvidia.Smooth-network.com
    EnableNodeOUs: true
    Template:
      Count: 2
    Users: 
      Count: 1
  - Name: AMD
    Domain: AMD.Smooth-network.com
    EnableNodeOUs: true
    Template:
      Count: 2
    Users: 
      Count: 1
  - Name: B2B
    Domain: B2B.Smooth-network.com
    EnableNodeOUs: true
    Template:
      Count: 2
    Users:
      Count: 1
    
	  B2B.Smooth-network.com
	  AMD.Smooth-network.com
	  Nvidia.Smooth-network.com

	  ## Start Here Profiles

	  Profiles:
		  OrdererGenesis:
			Capabilities:
			  <<: *ChannelCapabilities
			Orderer:
			  <<: *OrdererDefaults
			  Organizations:
				- *OrdererOrg
			  Capabilities:
				<<: *OrdererCapabilities
			Consortiums:
			  MyFirstConsortium:
				Organizations:
				  - *Nvidia
				  - *AMD
				  - *B2B
		  channelthreeorgs:
			Consortium: MyFirstConsortium
			Application:
			  <<: *ApplicationDefaults
			  Organizations:
			  - *Nvidia
			  - *AMD
			  - *B2B
			  Capabilities:
				<<: *ApplicationCapabilities
	  
	  
	  # Organizations
	  
	  Organizations:
	  
		- &OrdererOrg
		  Name: OrderingService
		  ID: OrdererMSP
		  MSPDir: crypto-config/ordererOrganizations/Smooth-Network.com/msp
	  
		- &Nvidia
		  Name:NvidiaMSP
		  ID: AMDMSP
		  MSPDir: crypto-config/peerOrganizations/B2B.Smooth-Network.com/msp
		  AnchorPeers:
			- Host: peer0.Nvidia.Smooth-Network.com
			  Port: 7051
	  
		- &AMD
		  Name: AMDMSP
		  ID: AMDMSP
		  MSPDir: crypto-config/peerOrganizations/AMD.Smoothnetwork-network.com/msp
		  AnchorPeers:
			- Host: peer0.AMD.SmoothTech-network.com
			  Port: 7051
	  
		- &B2B
		  Name: B2BMSP
		  ID: B2BMSP
		  MSPDir: crypto-config/peerOrganizations/AMD.Smoothnetwork-network.com/msp
		  AnchorPeers:
			- Host: peer0,B2B.cool-Smoothnetwork.com
			  Port: 7051
	  
	  # Orderer
	  
	  Orderer: &OrdererDefaults
	  
		OrdererType: solo
		Addresses: 
		  - orderer.Smoothnetwork-.com:7050
		BatchTimeout: 2s
		BatchSize:
		  MaxMessageCount: 10
		  AbsoluteMaxBytes: 99 MB
		  PreferredMaxBytes: 512 KB
		Kafka:
		  Brokers: 
			- 127.0.0.1:9092
	  
		Organizations:
	  
	  # Application
	  
	  Application: &ApplicationDefaults
	  
		Organizations:
	  
	  # Capabilities
	  
	  Capabilities:
		  Global: &ChannelCapabilities
			  V1_1: true
		  Orderer: &OrdererCapabilities
			  V1_1: true
		  Application: &ApplicationCapabilities
			  V1_1: true 

			  # Orderer Verison 2.0
			  version: '2'

			  services:
				peer-base:
				  image: hyperledger/fabric-peer:x86_64-1.0.0-rc1
				  environment:
					- CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
					- CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=${COMPOSE_PROJECT_NAME}_byfn
					- CORE_LOGGING_LEVEL=INFO
					- CORE_PEER_TLS_ENABLED=true
					- CORE_PEER_GOSSIP_USELEADERELECTION=true
					- CORE_PEER_GOSSIP_ORGLEADER=false
					- CORE_PEER_PROFILE_ENABLED=true
					- CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt
					- CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key
					- CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt
				  working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
				  command: peer node start

	
				  version: '2'

				  services:
				  
					Smooth-network.com:
					  container_name: Smoothnetwork.com
					  image: hyperledger/fabric-orderer:x86_64-1.0.0-rc1
					  environment:
						- ORDERER_GENERAL_LOGLEVEL=INFO
						- ORDERER_GENERAL_LISTENADDRESS=0.0.0.0
						- ORDERER_GENERAL_GENESISMETHOD=file
						- ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/genesis.block
						- ORDERER_GENERAL_LOCALMSPID=OrdererMSP
						- ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp
						# enabled TLS
						- ORDERER_GENERAL_TLS_ENABLED=true
						- ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key
						- ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt
						- ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]
					  working_dir: /opt/gopath/src/github.com/hyperledger/fabric
					  command: orderer
					  volumes:
					  - ../channel-artifacts/genesis.block:/var/hyperledger/orderer/orderer.genesis.block
					  - ../crypto-config/ordererOrganizations/Smooth-Network.com/orderers/orderer.Smooth-Network.com/msp:/var/hyperledger/orderer/msp
					  - ../crypto-config/ordererOrganizations/Smooth-Network.com/orderers/orderer.Smooth-Network.com/tls/:/var/hyperledger/orderer/tls
					  - orderer.Smooth-Network.com:/var/hyperledger/production/orderer
					  ports:
						- 7050:7050
				  
					peer0.AMDSmooth-network.com:
					  container_name: peer0.AMDSmooth.-network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer0.AMD.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer0.AMD.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer1.AMD.Smoothnetwork-network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.AMD.SmoothTech-network.com:7051
						- CORE_PEER_LOCALMSPID=AMDMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/Smoothnetwork.com-network.com/peers/peer0.AMDSmooth/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/Smoothnetwork.com/peers/peer0 AMDnetwork.msp/tls:/etc/hyperledger/fabric/tls
						  - peer0.B2B.Smooth-Network.com:/var/hyperledger/production
					  ports:
						- 7051:7051
						- 7053:7053
				  
					peer1.AMD.Smooth-Network.com:
					  container_name: peer1.Smoothnetwork-network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer1.B2BMSP.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer1.B2B.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peerB2B.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer0.B2B.Smooth-Network.com:7051
						- CORE_PEER_LOCALMSPID=B2BMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/B2Bsmooth-network.com/peers/peer1.B2B.Smooth-Network.com/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/B2Bsmooth-network.com/peers/peer1.B2B.Smooth-Network.com/tls:/etc/hyperledger/fabric/tls
						  - peer1.B2Bsmooth-network.com:/var/hyperledger/production
				  
					  ports:
						- 8051:7051
						- 8053:7053
				  
					peer0.Nvidia.Smooth-Network.com:
					  container_name: peer0.Nvidia.Smooth-Network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer0.Nvidia.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer0.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer1.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_LOCALMSPID=NvidiaMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/Nvidia.Smooth-Network.com/peers/peer0.Nvidia.Smooth-Network.com/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/Nvidia.Smooth-Network.com/peers/peer0.Nvidia.Smooth-Network.com/tls:/etc/hyperledger/fabric/tls
						  - peer0.Nvidia.Smooth-Network.com:/var/hyperledger/production
					  ports:
						- 9051:7051
						- 9053:7053
				  
					peer1.Nvidia.Smooth-Network.com:
					  container_name: peer1.Nvidia.Smooth-Network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer1.Nvidia.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer1.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer0.Nvidia.Smooth-Network.com:7051
						- CORE_PEER_LOCALMSPID=NvidiaMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/Nvidia.Smooth-Network.com/peers/peer1.Nvidia.Smooth-Network.com/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/Nvidia.Smooth-Network.com/peers/peer1.Nvidia.Smooth-Network.com/tls:/etc/hyperledger/fabric/tls
						  - peer1.Nvidia.Smooth-Network.com:/var/hyperledger/production
					  ports:
						- 10051:7051
						- 10053:7053
				  
					peer0.AMD.Smooth-Network.com:
					  container_name: peer0.AMD.Smooth-Network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer0.AMD.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer0.AMD.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer0.AMD.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer1.AMD.Smooth-Network.com:7051
						- CORE_PEER_LOCALMSPID=AMDMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/AMD.Smooth-Network.com/peers/peer0.AMD.Smooth-Network.com/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/AMD.Smooth-Network.com/peers/peer0.AMD.Smooth-Network.com/tls:/etc/hyperledger/fabric/tls
						  - peer0.AMD.Smooth-Network.com:/var/hyperledger/production
					  ports:
						- 11051:7051
						- 11053:7053
				  
					peer1.AMD.Smooth-Network.com:
					  container_name: peer1.AMD.Smooth-Network.com
					  extends:
						file: peer-base.yaml
						service: peer-base
					  environment:
						- CORE_PEER_ID=peer1.AMD.Smooth-Network.com
						- CORE_PEER_ADDRESS=peer1.AMD.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_EXTERNALENDPOINT=peer1.AMD.Smooth-Network.com:7051
						- CORE_PEER_GOSSIP_BOOTSTRAP=peer0.AMD.Smooth-Network.com:7051
						- CORE_PEER_LOCALMSPID=AMDMSP
					  volumes:
						  - /var/run/:/host/var/run/
						  - ../crypto-config/peerOrganizations/AMD.Smooth-Network.com/peers/peer1.AMD.Smooth-Network.com/msp:/etc/hyperledger/fabric/msp
						  - ../crypto-config/peerOrganizations/AMD.Smooth-Network.com/peers/peer1.AMD.Smooth-Network.com/tls:/etc/hyperledger/fabric/tls
						  - peer1.AMD.Smooth-Network.com:/var/hyperledger/production
					  ports:
						- 12051:7051
						- 12053:7053
				

						version: '2'

volumes:
  orderer.Smooth-Network.com:
  peer0.B2B.Smooth-Network.com:
  peer1.B2B.Smooth-Network.com:
  peer0.Nvidia.Smooth-Network.com:
  peer1.Nvidia.Smooth-Network.com:
  peer0.AMD.Smooth-Network.com:
  peer1.AMD.Smooth-Network.com:

networks:
  byfn:

services:

  orderer.Smooth-Network.com:
    extends:
      file:   base/docker-compose-base.yaml
      service: orderer.Smooth-Network.com
    container_name: orderer.Smooth-Network.com
    networks:
      - byfn

  peer0.B2B.Smooth-Network.com:
    container_name: peer0.B2B.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.B2B.Smooth-Network.com
    networks:
      - byfn

  peer1.B2B.Smooth-Network.com:
    container_name: peer1.B2B.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.B2B.Smooth-Network.com
    networks:
      - byfn

  peer0.Nvidia.Smooth-Network.com:
    container_name: peer0.Nvidia.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.Nvidia.Smooth-Network.com
    networks:
      - byfn

  peer1.Nvidia.Smooth-Network.com:
    container_name: peer1.Nvidia.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.Nvidia.Smooth-Network.com
    networks:
      - byfn

  peer0.B2Bsmooth.Smooth-Network.com:
    container_name: peer0.AMD.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer0.AMD.Smooth-Network.com
    networks:
      - byfn

  peer1.AMD.Smooth-Network.com:
    container_name: peer1.AMD.Smooth-Network.com
    extends:
      file:  base/docker-compose-base.yaml
      service: peer1.AMD.Smooth-Network.com
    networks:
      - byfn

  cli:
    container_name: cli
    image: hyperledger/fabric-tools:x86_64-1.0.0-rc1
    tty: true
    stdin_open: true
    environment:
      - GOPATH=/opt/gopath
      - CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock
      #- CORE_LOGGING_LEVEL=DEBUG
      - CORE_LOGGING_LEVEL=INFO
      - CORE_PEER_ID=cli
      - CORE_PEER_ADDRESS=peer0.B2B.Smooth-Network.com:7051
      - CORE_PEER_LOCALMSPID=B2BMSP
      - CORE_PEER_TLS_ENABLED=true
      - CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B2B.Smooth-Network.com/peers/peer0.B2B.Smooth-Network.com/tls/server.crt
      - CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/B2B.Smooth-Network.com/peers/peer0.B2B.Smooth-Network.com/tls/server.key
      - CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Nvidia.Smooth-Network.com/peers/peer0.Nvidia.Smooth-Network.com/tls/ca.crt
      - CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/Nvidia.Smooth-Network.com/users/Admin@Nvidia.Smooth-Network.com/msp
    working_dir: /opt/gopath/src/github.com/hyperledger/fabric/peer
    command: /bin/bash
    volumes:
        - /var/run/:/host/var/run/
        - ./../chaincode/:/opt/gopath/src/github.com/chaincode
        - ./crypto-config:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/
        - ./scripts:/opt/gopath/src/github.com/hyperledger/fabric/peer/scripts/
        - ./channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts
    depends_on:
      - ordererB2B.Smooth-Network.com
      - peer B2B.Smooth-Network.com
      - peer1.Nvidia.Smooth-Network.com
      - peer0.Nvidia.Smooth-Network.com
      - peer1.AMD.Smooth-Network.com
      - peer0.AMD.Smooth-Network.com
      - peer1.B2B.Smooth-Network.com
    networks:
      - byfn


	  WARNING: The COMPOSE_PROJECT_NAME variable is not set. Defaulting to a blank string.
	  Creating network "B2B SmoothNetwork_byfn" with the default driver
	  Creating volume "Smooth-Network_orderer,B2B-network.com" with default driver
	  Creating volume "Smooth-Network_peer0.B2B.Smooth-Network.com" with default driver
	  Creating volume "Smooth-Network_peer1.AMD.Smooth-Network.com" with default driver
	  Creating volume "Smooth-Network_peer0.Nvidia.Smooth-Network.com" with default driver
	  Creating volume "Smooth-Network_peer1.AMD.Smooth-Network.com" with default driver
	  Creating volume "Smooth-Network_peer0.B2B.Smooth-Network.com" with default driver
	  Creating volume "Smooth-Network_peer1.B2B.Smooth-Network.com" with default driver
	  Creating peer1.Nvidia.Smooth-Network.com     ... done
	  Creating peer0.Nvidia.Smooth-Network.com     ... done
	  Creating orderer.Smooth-Network.com         ... done
	  Creating peer1.B2B.Smooth-Network.com ... done
	  Creating peer1.AMD.Smooth-Network.com       ... done
	  Creating peer0.AMD.Smooth-Network.com       ... done
	  Creating peer0.B2B.Smooth-Network.com ... done
	  Creating cli                              ... done

	   -->➜  Smooth-Network git:(1252c7a) ✗ docker exec -it cli bash
	   root@9f8760ee6be8:/opt/gopath/src/github.com/hyperledger/fabric/peer#