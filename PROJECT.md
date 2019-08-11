# Features

# Design
## Command execution
Ansible copies module code to remote target, then generates a script that invoked that module and calls the script.
This is quite a lot of operations to perform.

# Architecture

Playbook
    Play
        Tasks
            Modules
                Transport


# Must have features
Must have features in order of importance. Having these feature should make this a usable project

## VARIABLES (1d)
group_vars, host_vars, set_fact module

## Templates (2d)
Use available library to render jinja templates in go

## WHEN ()
this conditional is use a lot. To support it i will need to support variables. Also, will need to support operators such
as "and,or,in,is,==, !=, >,>=, <,<="

# TODOS:
- re-use session
- remove hardcoded ssh key, add ability to provide your own key
- use normal logging mechanism so i can have different levels of log output
- create first performance test while i don't have too much functionality
