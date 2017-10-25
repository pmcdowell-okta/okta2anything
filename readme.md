![Okta2Anything](./images/okta2anythinglogo.png)


### Description

Okta2Anything is a simple LDAP Proxy which allows someone to have the Okta
IDAAS Service authenticate against almost anything.

This code is based on the great work by: https://github.com/vjeantet/ldapserver

### *Disclaimer*

*Although it is not uncommon for companies to use LDAP Proxies for authentication, this code is developed for Pilots, Proof of Concepts, and testing.*
*Anyone using this code for production is doing os at thier own risk*

### How does it work ?

Okta2Anything acts like an LDAP Service. Using the Okta LDAP Agent, and pointing
the LDAP Agent to Okta2Anything. Okta2Antyhing will delete the authentication
to a local node.js script that will perform the Authentication

![Okta2Anything](./images/flow.png)

![Okta2Anything](./images/linux.png)


### Configuration

#### Prerequisites

An Okta LDAP Agent is required, and node.js and any required artifacts needed for
your authentication scripts to run need to be installed.

You will need to download the binary for the Operating System you are running on

|  OS |   |   
|---|---|
| ![Okta2Anything](./images/linux.png)|   |   
| ![Okta2Anything](./images/macos.png)| https://github.com/pmcdowell-okta/okta2anything/raw/master/okta2anything.linux  |   
| ![Okta2Anything](./images/windows.png)|   |   

####