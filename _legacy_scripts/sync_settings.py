#!/usr/bin/python3
import re
import json
import sys


def get_fallback_values(source):
    struct_pattern = r"type GameSettings struct \{(.*?)\}"
    struct_match = re.search(struct_pattern, source, re.DOTALL)
    if not struct_match:
        raise ValueError("GameSettings struct not found")

    struct_block = struct_match.group(1)
    fallback_values = {}

    for line in struct_block.splitlines():
        line = line.strip()
        if line and not line.startswith("//"):
            parts = line.split()
            var_name = parts[0]
            var_type = parts[1]
            match = re.search(r"// (.*)", line)
            if match:
                fallback_value = match.group(1)
            else:
                fallback_value = "ERROR: missing"

            # Handling different types
            if var_type in ["int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64"]:
                fallback_value = int(fallback_value)
            elif var_type in ["float32", "float64"]:
                fallback_value = float(fallback_value)
            elif var_type == "bool":
                fallback_value = fallback_value.lower() == "true"
            elif var_type == "string":
                fallback_value = str(fallback_value)
            else:
                fallback_value = str(fallback_value)  # Handling all other types as string

            fallback_values[var_name] = fallback_value

    return fallback_values


def write_settings(fallback_values, output_file):
    with open(output_file, "w") as f:
        json.dump(fallback_values, f, indent=4)


def replace_fallback_settings(fallback_values, source):
    function_pattern = r"func UseFallbackSettings\(\) \{(.*?)\}"
    function_match = re.search(function_pattern, source, re.DOTALL)
    if not function_match:
        raise ValueError("UseFallbackSettings function not found")

    fallback_function = "func UseFallbackSettings() {\n\tsettings = &GameSettings{\n"
    for var, value in fallback_values.items():
        if isinstance(value, str):
            if value.isnumeric():  # Distinguishing between string numbers and string text
                fallback_function += f"\t\t{var}:  {value},\n"
            else:
                fallback_function += f"\t\t{var}:  \"{value}\",\n"
        elif isinstance(value, bool):
            fallback_function += f"\t\t{var}:  {str(value).lower()},\n"
        else:
            fallback_function += f"\t\t{var}:  {value},\n"
    fallback_function += "\t}"

    return re.sub(function_pattern, fallback_function, source, flags=re.DOTALL)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        input_file = "./fw/core/settings/settings.go"
        output_file = "./assets/settings.json"
    else:
        input_file, output_file = sys.argv[1], sys.argv[2]

    with open(input_file, "r") as f:
        source = f.read()

    fallback_values = get_fallback_values(source)
    write_settings(fallback_values, output_file)

    new_source = replace_fallback_settings(fallback_values, source)
    with open(input_file, "w") as f:
        f.write(new_source)
