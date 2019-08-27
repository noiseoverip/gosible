# Ansible playbook generator for benchamrking Ansible performance
import yaml

OUTPUT_DIR = "test1"
TEMPLATE_STEPS = {
    "count": 100,
    "args": {
        "template": {
            "src": "template_example.j2",
            "dest": "/tmp/ansiblego_test_template"
        }
    }
}
COPY_STEPS = 10
ECHO_STEPS = {
    "count": 10,
    "args": {
        "command": "echo labas"
    }
}
playbook = [
    { "hosts": "all",  "gather_facts": False, "tasks": []}
]

tasks_sets = [TEMPLATE_STEPS]

for t in tasks_sets:
    for n in range(t["count"]):
        playbook[0]["tasks"].append(t["args"])

play_text = yaml.dump(playbook)
open("test_templates_100.yaml", "w").write(play_text)