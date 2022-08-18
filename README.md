# This is a telegram bot that monitors a list of addresses and send a warning if they became under a specific limit

# To run the project:
1. 

```
cd tfgrid_monitoring_bot
```

2. Create a .env file containing the following vars [MNEMONICS, TFTS_LIMIT, BOT_TOKEN, CHAT_ID, MINS] where mins is the time between each call.

3. add the list of addresses you want to monitor in app.ts file.

4.

```
yarn install
```
- to run the project via docker:

```
docker build -t <tag> .
docker run -d <image name>
```
- To have your own flist:

```
docker push <DockerHub_username/image>
curl -f -X POST  -H "Authorization: bearer <ZeroHub_Token>" https://hub.grid.tf/api/flist/me/docker -F 'image=<dockerhub_username>/<image>'
```
- To run the deploy a VM based on the flist:

```
cd terraform_script
```
change the flist value with your newly built flist link

```
export MNEMONICS="<mnemonics words>"
export NETWORK="<network>"
terraform init && terraform apply -parallelism=1
```
