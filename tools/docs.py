#!/usr/bin/env python3

import sys
import json

PROVIDER = "transip"
PROVIDER_REGISTRY_NAME = f"registry.terraform.io/aequitas/{PROVIDER}"

INDEX_TEMPLATE = """
# {title} Provider

{description}

## Argument Reference

{arguments}
"""

TEMPLATE = """
# {title} {type}

{description}

## Argument Reference

{arguments}

## Attribute Reference

{attributes}
"""

ARGUMENT_TEMPLATE = "* `{name}` - ({required}) {description}"
ATTRIBUTE_TEMPLATE = "* `{name}` - {description}"

schema = json.loads(sys.stdin.read())

provider = schema["provider_schemas"][PROVIDER_REGISTRY_NAME]

for schematype in ["resource_schemas", "data_source_schemas"]:
    for resource, schema in provider[schematype].items():
        provider[schematype].items()
        resource_name = resource.replace(f"{PROVIDER}_", "")
        title = resource_name.replace("_", " ").title()
        description = schema.get("description", "")
        directory = "resources" if schematype == "resource_schemas" else "data-sources"
        attributes = schema["block"]["attributes"]

        with open(f"docs/{directory}/{resource_name}.md", "w") as f:
            arguments = [
                ARGUMENT_TEMPLATE.format(
                    name=name,
                    description=v.get("description", "n/a"),
                    required="Required" if v.get("required") else "Optional",
                )
                for name, v in attributes.items()
                if not v.get("computed")
            ]
            attributes = [
                ATTRIBUTE_TEMPLATE.format(
                    name=name, description=v.get("description", "n/a")
                )
                for name, v in attributes.items()
                if v.get("computed")
            ]

            f.write(
                TEMPLATE.format(
                    title=title,
                    type="Resource" if schematype == "resource_schemas" else "Data Source",
                    description=description,
                    arguments="\n".join(arguments or ["n/a"]),
                    attributes="\n".join(attributes or ["n/a"]),
                ).strip()
            )

with open("docs/index.md", "w") as f:
    title = PROVIDER.capitalize()

    arguments = [
        ARGUMENT_TEMPLATE.format(
            name=name,
            description=v.get("description", "n/a"),
            required="Required" if v.get("required") else "Optional",
        )
        for name, v in provider["provider"]["block"]["attributes"].items()
    ]

    f.write(
        INDEX_TEMPLATE.format(
            title=title,
            description=description,
            arguments="\n".join(arguments or ["n/a"]),
        ).strip()
    )

