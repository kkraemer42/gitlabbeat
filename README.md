Welcome to Gitlabbeat.



## Getting Started with Gitlabbeat

###IMPORTANT (DO THIS FIRST)###
After cloning this project, you will need to change the 'beat.yml' in the ```_meta``` folder:

```
gitlabbeat:
  # Defines how often an event is sent to the output
  period: 20s

  job_timeout: 10s

  access_token: 'YOUR GITLAB ACCESS TOKEN'
  gitlab_address: 'YOUR_GITLAB_ADDRESS/api/v4'
```

### Init Project
To get running with Gitlabbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.


### Build

To build the binary for Gitlabbeat run the command below. This will generate a binary
in the same directory with the name gitlabbeat.

```
make
```


### Run

To run Gitlabbeat with debugging output enabled, run:

```
./gitlabbeat -c gitlabbeat.yml -e -d "*"
```
