import json

def json_to_env_string(file_path):
    # Read the JSON file
    with open(file_path, 'r') as file:
        data = json.load(file)

    # Convert the JSON object to a string and escape it for environment variable
    return json.dumps(data).replace('\n', '\\n')

# Replace 'path/to/your/service-account.json' with the actual path to your JSON file
escaped_json_string = json_to_env_string('./credentials.json')
print(escaped_json_string)
