import argparse
import json

def parse_record(record):
    pass

def load_records(path):
    f = open(path)
    for line in f:
        data = json.loads(line)
        print(data['inCitations'])

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        description='Load the citation graph into memory.',
    )
    parser.add_argument('path')
    args = parser.parse_args()
    path = args.path
    load_records(path)
