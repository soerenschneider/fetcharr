# Adding a transmission webhook

## Adding the webhook script

Download the example script [here](contrib/invoke-webhook.sh) that calls fetcharr's webhook endpoint via curl.
You'll need to edit the file to change the host to the correct value in line 3.

```bash
$ sudo curl -o /usr/local/bin/invoke-fetcharr-webhook https://github.com/soerenschneider/fetcharr/blob/main/contrib/invoke-webhook.sh
$ sudo chmod +x /usr/local/bin/invoke-fetcharr-webhook
```

To test if everything works, simply execute the script file. Fetcharr should now try to connect to your seedbox and download files. On case of an error, curl should display logs that help troubleshooting the problem.

```bash
$ /usr/local/bin/invoke-fetcharr-webhook
```

## Configuring Transmission

For Transmission to automatically call the webhook once a file is downloaded, add the following two lines to Transmission's configuration. Make sure to update the path to your webhook.
```json
"script-torrent-done-enabled": true,
"script-torrent-done-filename": "/usr/local/bin/invoke-fetcharr-webhook"
```
