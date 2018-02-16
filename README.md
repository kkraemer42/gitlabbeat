Welcome to Gitlabbeat.



## Getting Started with Gitlabbeat

After cloning this project, you will need to place a 'beat.yml' in the ```_meta``` folder.

It's contents should look like this:

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
./gitlabbeat -c gitlabeat.yml -e -d "*"
```



### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
