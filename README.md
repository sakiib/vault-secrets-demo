# vault-secrets-demo

## Install KubeDB Enterprise operator chart

```bash
$ helm install kubedb appscode/kubedb \
    --version v2021.12.21 \
    --namespace kubedb --create-namespace \
    --set kubedb-enterprise.enabled=true \
    --set kubedb-autoscaler.enabled=true \
    --set-file global.license=/path/to/the/license.txt
```

## Install KubeVault Enterprise operator chart

```bash
$ helm install kubevault appscode/kubevault \
    --version v2022.01.11 \
    --namespace kubevault --create-namespace \
    --set-file global.license=/path/to/the/license.txt
```

## Install Secret-store CSI Driver

```bash
$ helm install csi-secrets-store secrets-store-csi-driver/secrets-store-csi-driver --namespace kube-system
```

## Install Vault specific CSI Provider

```bash
# using helm
$ helm install vault hashicorp/vault \
    --set "server.enabled=false" \
    --set "injector.enabled=false" \
    --set "csi.enabled=true"
     
# or using provider yaml
$ kubectl apply -f provider.yaml
```

## Deploy TLS Secured VaultServer

A VaultServer is a Kubernetes CustomResourceDefinition (CRD) which is used to deploy a HashiCorp Vault server on Kubernetes clusters in a Kubernetes native way.

When a VaultServer is created, the KubeVault operator will deploy a Vault server and create necessary Kubernetes resources required for the Vault server.

```bash
# create the issuer
$ kubectl apply -f issuer.yaml

# deploy the vault server
$ kubectl apply -f vaultserver.yaml
```

## Export necessary environment variables

```bash
$ export VAULT_ADDR='https://127.0.0.1:8200'

$ export VAULT_SKIP_VERIFY=true

$ export VAULT_TOKEN=(kubectl vault get-root-token vaultserver vault -n demo --value-only) 
```

## Enable MySQL SecretEngine

A SecretEngine is a Kubernetes CustomResourceDefinition (CRD) which is designed to automate the process of enabling and configuring secret engines in Vault in a Kubernetes native way.

```bash
# create mysql DB 
$ kubectl apply -f mysql.yaml

# enable secret engine
$ kubectl apply -f secretengine.yaml
```

## Get decrypted Vault Root Token

```bash
# get the decrypted root token with name
$ kubectl vault get-root-token vaultserver vault -n demo

# get only the value of decrypted root token
$ kubectl vault get-root-token vaultserver vault -n demo --value-only
```

## Create Database Roles

A MySQLRole is a Kubernetes CustomResourceDefinition (CRD) which allows a user to create a database secret engine role in a Kubernetes native way.

When a MySQLRole is created, the KubeVault operator creates a role according to the specification.

```bash
# create the superuser role
$ kubectl apply -f superusr-role.yaml

# create the readonly role
$ kubectl apply -f readonly-role.yaml
```

## Create SecretAccessRequest

A SecretAccessRequest is a Kubernetes CustomResourceDefinition (CRD) which allows a user to request a Vault server for credentials in a Kubernetes native way. A SecretAccessRequest can be created under various roleRef e.g: AWSRole, GCPRole, ElasticsearchRole, MongoDBRole, etc. A SecretAccessRequest has three different phases e.g: WaitingForApproval, Approved, Denied. If SecretAccessRequest is approved, then the KubeVault operator will issue credentials and create Kubernetes secret containing credentials. 

```bash
$ kubectl apply -f secretaccessrequest.yaml
```

## Approve/Deny SecretAccessRequest

```bash
# upon approval of secret access request, secrets with username/password will be created
$ kubectl vault approve secretaccessrequest mysql-cred-req -n dev

# deny secret access request
$ kubectl vault deny secretaccessrequest mysql-cred-req -n dev
```

## Create ServiceAccount & SecretRoleBinding

A SecretRoleBinding is a Kubernetes CustomResourceDefinition (CRD) which allows a user to bind a set of roles to a set of users. Using the SecretRoleBinding it’s possible to bind various roles e.g: AWSRole, GCPRole, ElasticsearchRole, MongoDBRole, etc. to Kubernetes ServiceAccounts. A SecretRoleBinding has three different phases e.g: Processing, Success, Failed. Once a SecretRoleBinding is successful, it will create a VaultPolicy and a VaultPolicyBinding.

```bash
# create the service account
$ kubectl apply -f serviceaccount.yaml

# create the secret role binding
$ kubectl apply -f secretrolebinding.yaml
```

## Create SecretProviderClass using KubeVault CLI

```bash
# Generate secretproviderclass for the MySQL username and password
$ kubectl vault generate secretproviderclass vault-db-provider -n test      \
    --secretrolebinding=dev/secret-r-binding \
    --vaultrole=MySQLRole/readonly-role \
    --keys username=sql-user --keys password=sql-pass -o yaml 
```

## Deploy the Microservice

One of the really cool concepts behind Vault is dynamic secrets. And when we talk about secret sprawl, the ability to have the same username and password distributed out across your fleet allows an attacker to attack one insecure area and then gain secrets across your entire environment.

Dynamic secrets changed this paradigm a little bit by having each of your endpoints get its own username and password for the entity that you're trying to get access to.

Most of these dynamic secrets are timebound and easily revocable, so if you notice that there's an issue or a breach inside your environment, you can revoke one secret, while all the rest of your applications have other usernames and passwords. 

```bash
# create a microservice deployment
$ kubectl apply -f microservice.yaml
```

## MySQL Queries

```bash
# login as the root user
$ mysql -uroot -p$MYSQL_ROOT_PASSWORD

# loging using the username, password
$ mysql -u <username> -p

# show the databases
$ show databases;

# use the <db-name>
$ use <db-name>;

# show tables
$ show tables;

# create a table with name product with column <id, name, price>
$ create table product(id int, name varchar(100), price float);

# insert values into product table
$ insert into product(id, name, price) values(1, "pen", 3.5);
$ insert into product(id, name, price) values(2, "book", 7.5);

# select everything from the product table
$ select * from product;
```

## Revoke the SecretAccessRequest

```bash
$ kubectl vault revoke secretaccessrequest mysql-cred-req -n dev
```

## Delete the VaultServer

```bash
$ kubectl delete -f vaultserver.yaml
```
