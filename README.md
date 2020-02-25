# Detailed features list
This is a list of Ansible features that are implemented in AnsibleGO

* Read inventory group_vars
* Read play yaml file
* Template module - render jinja2 templates to target host
* Command module - execute shell commands
* Task attribute: register. Register task result as variable
* Task attribute: when. Conditionally execute a task
* Set_fact module: register given value as variable

# If you don't know what Ansible is
This is for my future self (if the day comes when I have forgotten what Ansible is...) as well as for those who actually dont' know

# AnsibleGO vs Ansible

## Command execution
Ansible copies module code to remote target, then generates a script that invoked that module and calls the script.
This is quite a lot of operations to perform. AnsibleGO simply executes provided command on target host, this has
certain limitations but is a much faster approach

## Integration tests
Integration tests require a target host.
FOr now multipass can be used (TODO: automate multi pass setup)
