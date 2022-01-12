# vault-secrets-demo

# Deploy TLS Secured VaultServer

```bash
# create the issuer
$ kubectl apply -f issuer.yaml

# deploy the vault server
$ kubectl apply -f vaultserver.yaml
```

# Enable MySQL SecretEngine

```bash
# create mysql DB 
$ kubectl apply -f mysql.yaml

# enable secret engine
$ kubectl apply -f secretengine.yaml
```

# Get decrypted Vault Root Token

```bash
# get the decrypted root token with name
$ kubectl vault get-root-token vaultserver vault -n demo

# get only the value of decrypted root token
$ kubectl vault get-root-token vaultserver vault -n demo --value-only
```

# Create Database Roles

```bash
# create the superuser role
$ kubectl apply -f superusr-role.yaml

# create the readonly role
$ kubectl apply -f readonly-role.yaml
```

# Create SecretAccessRequest

```bash
$ kubectl apply -f secretaccessrequest.yaml
```

# Approve/Deny SecretAccessRequest

```bash
# upon approval of secret access request, secrets with username/password will be created
$ kubectl vault approve secretaccessrequest mysql-cred-req -n dev

# deny secret access request
$ kubectl vault approve secretaccessrequest mysql-cred-req -n dev
```

# Deploy Microservice to demonstrate Vault Dynamic Secrets

```bash
# create the service account
$ kubectl apply -f serviceaccount.yaml

# create the secret role binding
$ kubectl apply -f secretrolebinding.yaml
```

# Create SecretProviderClass using KubeVault CLI

```bash
# Generate secretproviderclass for the MySQL username and password
$ kubectl vault generate secretproviderclass vault-db-provider -n test      \
    --secretrolebinding=dev/secret-r-binding \
    --vaultrole=MySQLRole/readonly-role \
    --keys username=sql-user --keys password=sql-pass -o yaml 
```

# Create the Microservice

```bash
# create a microservice deployment
$ kubectl apply -f microservice.yaml
```

# MySQL Queries

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

# Revoke the SecretAccessRequest

```bash
$ kubectl vault revoke secretaccessrequest mysql-cred-req -n dev
```

# Delete the VaultServer

```bash
$ kubectl delete -f vaultserver.yaml
```