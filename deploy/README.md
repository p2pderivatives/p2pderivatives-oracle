# Deploy
## Requirements
 - initialized docker swarm (install docker engine and run `docker swarm init`)
 - access to the docker swarm manager
 You can set up a remote context through ssh using this command
 (need to have your ssh public key registered in the remote machine):  
 `docker context create <context-name> ‐‐docker “host=ssh://<user>@<remote-machine-ip>”`
## Secrets and Configuration
The secrets and configuration are external and need to be created beforehand.  
examples:
- `cat <my-secret-file> | docker secret create <my-secret> -`
- `cat <my-config-file> | docker config create <my-config> -`
## Deployment
This project can be deployed using docker-compose (image tag required and eventual mandatory variables):  
```bash
IMAGE_TAG=<image-tag> \
docker stack deploy \
--compose-file ./deploy/docker-compose.yml \
<stack-name> \
--with-registry-auth
```  
you can check the status of your newly created services using:  
`docker stack services <stack-name>`

You can remove your deployment using `docker stack rm <stack-name>`
(the volumes remains and needs to be manually deleted)

## Update
You can update your services using the same command as deployment.
The secrets and configurations can be updated by creating a new one (immutable) and replacing it manually. 
