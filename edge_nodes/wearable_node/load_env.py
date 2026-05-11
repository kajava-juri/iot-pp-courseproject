# Reads .env file and adds them as CMAKE variables
# code partially taken from this stack overflow anser: https://stackoverflow.com/questions/62314497/access-of-outer-environment-variable-in-platformio
from os.path import isfile

Import("env")

env_file = "../../.env"

assert isfile(env_file)
try:
    f = open(env_file, "r")

    for line in f:
        if line.startswith("#") or not line.strip():
            continue
        key, value = line.split("=", 1)
        key = key.strip()
        value = value.strip()
        # remove surrounding quotes if present
        if (value.startswith('"') and value.endswith('"')) or (value.startswith("'") and value.endswith("'")):
            value = value[1:-1]
        # strip stray trailing tokens like ' -D' if present (actually dont know how it happens, but need to remove it if present)
        if value.endswith(' -D'):
            value = value[:-3].rstrip()
        print(f'Adding CPPDEFINE: {key}={value}')
        env.Append(CPPDEFINES=[(key, value)])
except IOError:
    print("File .env is not accessible")
finally:
    f.close()