# GKE Money Saver

If you have a GKE cluster setup to run CI jobs, such as using Gitlab CI Runner, or any other ephemeral workload and you want to auto-delete the node pool
when people go home for the day this should help.

## What

This small Go application can be easily deployed to Google App Engine and be setup to automatically create or delete a node pool for GKE on a schedule.

## Setup

Clone the repo then edit app.yaml with the environment variables for your setup.

```
MACHINE_TYPE: optional, defaults to n1-standard-1, the google cloud machine type to create
PROJECT: required, your Google project id
LOCATION: required, the region you want the node pool to be created in
CLUSTER: required, the name of the cluster that node pool should be long to
NODE_POOL: optional, defaults to pool-1, the name of the node pool
SERVICE_ACCOUNT: required, the name of the service account that the node pool servers should use
OAUTH_SCOPES: required, comma seperated list of oauth_scopes to be attached to the node pool servers
RANDOM_PATH: optional, defaults to 8YkpGd2LQN3nBuWJfXRb, adds a random URL path to the endpoints that control behaviour
PREEMPTIBLE: optional, default to true, should the servers in the node pool be preemptible
DISK_SIZE: optional, default to 50, size in GB of the disks for the servers in the node pool
INITIAL_NODE_COUNT: optional, default to 1, how many servers should be in the pool on create
MAX_NODE_COUNT: optional, default to 2, how many servers should the pool autoscale up to 
```

Then edit cron.yaml to fit your schedule.

## Deploy

Deploy to Google App Engine using the GCloud tool.

```
gcloud app deploy --project "yourprojectid"
```

Test it out by going to `https://your-appengine-url/$RANDOM_PATH/up` to bring up the node pool and `/down` to bring it down again.
If all is well deploy the cron job to bring it up and down on a schedule.

```
gcloud app deploy --project "yourprojectid" cron.yaml
```

Your done!

Pull requests gladly accepted!