# Reads .env file and adds them as CMAKE variables
# code partially taken from this stack overflow anser: https://stackoverflow.com/questions/62314497/access-of-outer-environment-variable-in-platformio
from pathlib import Path

Import("env")

# Navigate from project root to repository root (.env is at repo root)
project_root = Path(env.subst("$PROJECT_DIR"))
env_file = project_root.parent.parent / ".env"

try:
    if not env_file.exists():
        raise IOError(f"File {env_file} not found")
    
    with open(env_file, "r") as f:
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
            # safe_value = '"' + value.replace('"', '\\"') + '"'
            # print(f'Adding CPPDEFINE: {key}={safe_value}')
            # env.Append(CPPDEFINES=[f'{key}={safe_value}'])
            env.Append(CPPDEFINES=[(key, value)])
except IOError:
    print("File .env is not accessible")
finally:
    pass
