import random
from datetime import datetime, timedelta
import json
import urllib
import urllib2

COLORS = ['red', 'orange', 'yellow', 'green', 'blue']
PIECES = ['bishop', 'knight', 'rook', 'pawn', 'queen', 'king']
TAGS = ['Luna', 'Phobos', 'Deimos', 'Ceres', 'Pallas', 'Vesta', 'Juno', 'Ganymede', 'Io', 'Europa', 'Callisto', 'Titan', 'Mimas']
COMMAND_BITS = ['zeus', 'hera', 'poseidon', 'demeter', 'dionysus', 'apollo', 'artemis', 'hermes', 'athena', 'ares', 'aphrodite', 'hephaestus']


def get_hostname():
    return random.choice(COLORS) + '.' + random.choice(PIECES)

def get_cmd_name():
    return 'command' + str(random.randint(1, 10) * 10)

for i in xrange(1000):
    data = {
        'name': get_cmd_name(),
        'hostname': get_hostname(),
        'tags': random.sample(TAGS, random.randint(0, 3)),
        'command': random.sample(COMMAND_BITS, random.randint(2, 7)),
        'exit_code': random.randint(0, 1),
        'start_time': (datetime.utcnow() - timedelta(seconds=random.randint(0, 1000))).strftime("%Y-%m-%dT%H:%M:%SZ%Z"),
        'elapsed_seconds': random.random() * 30
    }

    req = urllib2.Request('http://localhost:8080/report', json.dumps(data).encode("utf8"), {'Content-Type': 'application/json'})
    try:
        response = urllib2.urlopen(req)
    except Exception as e:
        print e.read()
        raise
