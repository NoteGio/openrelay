import yaml
import sys

config = yaml.load(open(sys.argv[1]))
images = set(
    i.get("image").split(":")[0] for i in config["services"].values() if "${TAG:-latest}" in i.get("image", "")
)

for image in images:
    print(image)
