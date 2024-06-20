import sys
import yaml


def sum_emotional_profile_fields(file_content):
    # Load the YAML content
    nodes = yaml.safe_load(file_content)

    # Initialize a dictionary to hold the sum of each field
    field_sums = {
        "Tradition": 0,
        "Security": 0,
        "Conformity": 0,
        "Achievement": 0,
        "Power": 0,
        "Hedonism": 0,
        "Stimulation": 0,
        "SelfDirection": 0,
        "Universalism": 0,
        "Benevolence": 0
    }

    # Iterate through each node and response to sum the fields
    for node in nodes:
        for response in node['responses']:
            for field, value in response['influence'].items():
                field_sums[field] += value

    return field_sums


# Example YAML content
with open(sys.argv[1]) as f:
    yaml_content = f.read()
    x = sum_emotional_profile_fields(yaml_content)
    for k, v in x.items():
        print(k, v)
