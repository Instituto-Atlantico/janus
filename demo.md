# Demo

For this we need two agents running, a issuer and a holder.

### connecting

1. On **issuer** create a invitation
2. Copy the invitation and receive it on **holder**.
3. On **holder** accept the invitation
4. On **issuer** accept the request

### sending messages

5. On **holder** send a message to the **issuer**
6. Check the message on the **webhock server**

### providing credentials

7. On **issuer** create a schema and credential definition
8. On **issuer** send a credential offer to **holder**
9. On **holder** send a credential request to **issuer**



## jsons

### schema

``` json
{
   "attributes":[
      "name",
      "age"
   ],
   "schema_name":"test-schema",
   "schema_version":"1.0"
}
```

### credential definition

``` json
{
  "revocation_registry_size": 1000,
  "schema_id": "22222222222222:2:test-schema:1.0",
  "support_revocation": false,
  "tag": "default"
}
```

### send credential

``` json
{
  "auto_remove": true,
  "comment": "second try:D",
  "connection_id": "111111111111",
  "credential_preview": {
    "@type": "issue-credential/2.0/credential-preview",
    "attributes": [
      {
        "mime-type": "plain/text",
        "name": "name", 
        "value": "Bob"
      },
      {
        "mime-type": "plain/text",
        "name": "age", 
        "value": "30"
      }
    ]
  },
  "filter": {
    "indy": {
      "cred_def_id": "22222222222222:3:CL:17926:default",
      "issuer_did": "22222222222222",
      "schema_id": "22222222222222:2:test-schema:1.0",
      "schema_issuer_did": "22222222222222",
      "schema_name": "test-schema",
      "schema_version": "1.0"
    }
  },
  "trace": false
}
```