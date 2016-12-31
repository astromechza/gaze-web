#!/usr/bin/env python

import os
from textwrap import dedent
import subprocess

"""
This script can be launched to generate the README.md file. This script is a useful way of ensuring that your README
document is up to date and the relevant examples still work as intended.

```
$ ./generate_README.py -o
or
$ ./generate_README.py > README.md
```
"""

# the file that will be written
PROJECT_DIRECTORY = os.path.dirname(__file__)


class Generator(object):
    def __init__(self):
        self.lines = []

    def heading(self, text, level):
        self.lines.append('#' * level + " " + text)
        self.lines.append("")

    def h1(self, text):
        return self.heading(text, 1)

    def h2(self, text):
        return self.heading(text, 2)

    def h3(self, text):
        return self.heading(text, 3)

    def h4(self, text):
        return self.heading(text, 4)

    def paragraph(self, text):
        text = text.rstrip()
        self.lines.append(dedent(text))
        self.lines.append("")

    def command_example(self, command):
        self.lines.append("```")
        self.lines.append("$ {}".format(command))

        try:
            output = subprocess.check_output(command, stderr=subprocess.STDOUT, shell=True, cwd=PROJECT_DIRECTORY)
        except subprocess.CalledProcessError as e:
            output = e.output

        self.lines.append(output.strip())
        self.lines.append("```")
        self.lines.append("")

    def __str__(self):
        text = "\n".join(self.lines)
        if not text.endswith("\n"):
            text += "\n"
        return text


def main():
    g = Generator()

    g.h1("gaze-web` - a web application for serving `gaze` records")
    g.paragraph("""\
    This web app can capture the JSON payloads produced by `gaze` (see [here](https://github.com/AstromechZA/gaze)) and display 
    them while allowing nice paginated and searchable lists of the results.
    """)

    print str(g)


if __name__ == '__main__':
    main()
