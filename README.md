# Gosbile - a naive attempt to rewrite Ansible in go
As title says, this is a very naive attempt to rewrite Ansible in go. 

Why ?
Because Ansible is growing too big, is attempting to do more than it used and becoming too slow. So, having
a simpler and faster incarnation of it might actually be handy for smaller projects.

# Feature list
This is a list of Ansible features that are implemented in gosible.

* Read inventory group_vars
* Read play yaml file
* Module "template" - render jinja2 templates to target host
* Module "command" - execute shell commands
* Task attribute: register. Register task result as variable
* Task attribute: when. Conditionally execute a task
* Module "set_fact" : register given value as variable
* Module "debug"
* Module "copy" - copy files to target hosts

# Building
`make build` should be enough

# Testing
`make test` - runs unit tests which do not have any external dependencies

`make testint` - run integration tests which require some host to exist. Define env var HOST with IP address of a linux 
VM when running these tests and make sure your user can SSH to that VM with default ssh keys.

`make bench` - run few different playbooks with Gosible and Ansible and compare times. Needs Ansible and linux VM (env var HOST)