#!/usr/bin/with-contenv bash

### Set Defaults
  tokensFromEnv="LDAP_HOST LDAP_BIND_DN LDAP_BIND_PW LDAP_BASE_DN LDAP_ATTRIBUTE LDAP_SCOPE LDAP_FILTER LDAP_GROUP_ATTRIBUTE"

  LDAP_ATTRIBUTE=${LDAP_ATTRIBUTE:="uid"}
  LDAP_SCOPE=${LDAP_SCOPE:="sub"}
  LDAP_FILTER=${LDAP_FILTER:="(objectClass=person)"}
  LDAP_GROUP_ATTRIBUTE=${LDAP_GROUP_ATTRIBUTE:="uniquemember"}
  UPLOAD_MAX_SIZE=${UPLOAD_MAX_SIZE:="2G"}
  PHP_TIMEOUT=${PHP_TIMEOUT:="180"}

### Adjust NGINX Runtime Variables
  
  sed -i -e "s/<UPLOAD_MAX_SIZE>/$UPLOAD_MAX_SIZE/g" /etc/nginx/nginx.conf
  sed -i -e "s/<PHP_TIMEOUT>/$PHP_TIMEOUT/g" /etc/nginx/conf.d/02-default.conf

### LDAP Setup
  for envVar in $tokensFromEnv; do
    envValue=$(echo "${!envVar}" | sed -e 's/[&\\\$]/\\&/g')
	sed -i -e "s|\$${envVar}|${envValue}|g" /etc/nginx/conf.d/01-ldap.conf;
  done

